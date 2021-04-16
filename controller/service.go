package controller

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
	"strconv"
	"strings"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDelete)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
}

//ServiceList godoc
//@Summary Service list
//@Description Show service lists.
//@Tags Service Management
//@ID /service/service_list
//@Accept  json
//@Produce  json
//@Param info query string false "Searching keyword"
//@Param page_size query int true "Entries per page"
//@Param page_no query int true "Page No."
//@Success 200 {object} middleware.Response{data=io.ServiceListOutput} "success"
//@Router /service/service_list [get]
func (service *ServiceController) ServiceList(c *gin.Context) {
	params := &io.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &db.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outList := []io.ServiceListItemOutput{}

	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := "unknown"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRulePrefixURL {
			if serviceDetail.HTTPRule.NeedHttps == 0 {
				serviceAddr = clusterIP + ":" + clusterPort + serviceDetail.HTTPRule.Rule
			} else {
				serviceAddr = clusterIP + ":" + clusterSSLPort + serviceDetail.HTTPRule.Rule
			}
		}

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}

		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = clusterIP + ":" + strconv.FormatInt(serviceDetail.TCPRule.Port, 10)
		}

		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = clusterIP + ":" + strconv.FormatInt(serviceDetail.GRPCRule.Port, 10)
		}

		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		outItem := io.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			QPD:         0,
			QPS:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}

	out := &io.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)

}

//ServiceDelete godoc
//@Summary Delete a service
//@Description Delete a service.
//@Tags Service Management
//@ID /service/service_delete
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &io.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &db.ServiceInfo{ID: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)

	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}

//ServiceAddHTTP godoc
//@Summary Add HTTP service
//@Description  Add HTTP service
//@Tags Service Management
//@ID /service/service_add_http
//@Accept  json
//@Produce  json
//@Param body body io.ServiceAddHTTPInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_add_http [post]
func (adminlogin *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &io.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	serviceInfo := &db.ServiceInfo{ServiceName: params.ServiceName}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx = tx.Begin()

	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("Service name has already exists!"))
		return
	}

	httpUrl := &db.HTTPRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Find(c, tx, httpUrl); err == nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("Service rule type has already exists!"))
		return
	}

	if len(strings.Split(params.IpList, "\n")) !=
		len(strings.Split(params.WeightList, "\n")) {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2004, errors.New("IP list length and Weight list length are not equal!!"))
		return
	}

	serviceModel := &db.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2005, errors.New("IP list length and Weight list length are not equal!!"))
		return
	}
	httpRule := &db.HTTPRule{
		ServiceId:       serviceModel.ID,
		RuleType:        params.RuleType,
		Rule:            params.Rule,
		NeedHttps:       params.NeedHttps,
		NeedStripUri:    params.NeedStripUri,
		NeedWebsocket:   params.NeedWebsocket,
		UrlRewrite:      params.UrlRewrite,
		HeaderTransform: params.HeadTransform,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2006, errors.New("IP list length and Weight list length are not equal!!"))
		return
	}

	accessControl := &db.AccessControl{
		ServiceId:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIpFlowLimit: params.ClientIpFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2007, errors.New("accessControl save failed."))
		return
	}

	loadBalance := &db.LoadBalance{
		ServiceId:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeOut,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeOut,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeOut,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2008, errors.New("loadBalance add failed."))
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")

}

//ServiceUpdateHTTP godoc
//@Summary Update HTTP service
//@Description  Update HTTP service
//@Tags Service Management
//@ID /service/service_update_http
//@Accept  json
//@Produce  json
//@Param body body io.ServiceUpdateHTTPInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_update_http [post]
func (adminlogin *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &io.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	if len(strings.Split(params.IpList, "\n")) !=
		len(strings.Split(params.WeightList, "\n")) {
		middleware.ResponseError(c, 2001, errors.New("IP list length and Weight list length are not equal!!"))
		return
	}

	serviceInfo := &db.ServiceInfo{ServiceName: params.ServiceName}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	tx = tx.Begin()
	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err == nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("Service does not exist!"))
		return
	}

	serviceDetail.HTTPRule.NeedHttps = params.NeedHttps
	serviceDetail.HTTPRule.NeedStripUri = params.NeedStripUri
	serviceDetail.HTTPRule.NeedWebsocket = params.NeedWebsocket
	serviceDetail.HTTPRule.UrlRewrite = params.UrlRewrite
	serviceDetail.HTTPRule.HeaderTransform = params.HeadTransform

	if err := serviceDetail.HTTPRule.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2004, errors.New("httpRule update failed!!"))
		return
	}

	serviceDetail.AccessControl.OpenAuth = params.OpenAuth
	serviceDetail.AccessControl.BlackList = params.BlackList
	serviceDetail.AccessControl.WhiteList = params.WhiteList

	serviceDetail.AccessControl.ClientIpFlowLimit = params.ClientIpFlowLimit
	serviceDetail.AccessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := serviceDetail.AccessControl.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2005, errors.New("accessControl update failed."))
		return
	}

	serviceDetail.LoadBalance.RoundType = params.RoundType
	serviceDetail.LoadBalance.IpList = params.IpList
	serviceDetail.LoadBalance.WeightList = params.WeightList
	serviceDetail.LoadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeOut
	serviceDetail.LoadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeOut
	serviceDetail.LoadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeOut
	serviceDetail.LoadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := serviceDetail.LoadBalance.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2006, errors.New("loadBalance update failed."))
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")

}

//ServiceDetail godoc
//@Summary Delete a service
//@Description Delete a service.
//@Tags Service Management
//@ID /service/service_detail
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_detail [get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &io.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &db.ServiceInfo{ID: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)

	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceInfo.IsDelete = 1
	if err := serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}
