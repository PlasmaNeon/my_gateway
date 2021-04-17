package controller

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
	"strings"
	"time"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
	group.POST("/service_add_tcp", service.ServiceAddTCP)
	group.POST("/service_update_tcp", service.ServiceUpdateTCP)
	group.POST("/service_add_grpc", service.ServiceAddGRPC)
	group.POST("/service_update_grpc", service.ServiceUpdateGRPC)
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
			serviceAddr = clusterIP + ":" + string(rune(serviceDetail.TCPRule.Port))
		}

		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = clusterIP + ":" + string(rune(serviceDetail.GRPCRule.Port))
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
	serviceInfo := &db.ServiceInfo{ID: params.ID}
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

	if len(strings.Split(params.IPList, "\n")) !=
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
		NeedStripUri:    params.NeedStripURI,
		NeedWebsocket:   params.NeedWebsocket,
		UrlRewrite:      params.URLRewrite,
		HeaderTransform: params.HeadTransform,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2006, errors.New("IP list length and Weight list length are not equal!!"))
		return
	}

	accessControl := &db.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2007, errors.New("accessControl save failed."))
		return
	}

	loadBalance := &db.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IPList,
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
	if len(strings.Split(params.IPList, "\n")) !=
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
	serviceDetail.Info.ServiceDesc = params.ServiceDesc
	//serviceDetail.Info.Save(c, tx)

	if err := serviceDetail.Info.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2004, errors.New("httpRule update failed!!"))
		return
	}

	serviceDetail.HTTPRule.NeedHttps = params.NeedHttps
	serviceDetail.HTTPRule.NeedStripUri = params.NeedStripURI
	serviceDetail.HTTPRule.NeedWebsocket = params.NeedWebsocket
	serviceDetail.HTTPRule.UrlRewrite = params.URLRewrite
	serviceDetail.HTTPRule.HeaderTransform = params.HeadTransform

	if err := serviceDetail.HTTPRule.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2005, errors.New("httpRule update failed!!"))
		return
	}

	serviceDetail.AccessControl.OpenAuth = params.OpenAuth
	serviceDetail.AccessControl.BlackList = params.BlackList
	serviceDetail.AccessControl.WhiteList = params.WhiteList

	serviceDetail.AccessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	serviceDetail.AccessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := serviceDetail.AccessControl.Save(c, tx); err != nil {
		tx = tx.Rollback()
		middleware.ResponseError(c, 2005, errors.New("accessControl update failed."))
		return
	}

	serviceDetail.LoadBalance.RoundType = params.RoundType
	serviceDetail.LoadBalance.IpList = params.IPList
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

// ServiceAddTCP godoc
// @Summary Add TCP service.
// @Description Add TCP service.
// @Tags Service Management
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body io.ServiceAddTCPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (admin *ServiceController) ServiceAddTCP(c *gin.Context) {
	params := &io.ServiceAddTCPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &db.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &db.TCPRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &db.GRPCRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &db.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	loadBalance := &db.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IPList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	httpRule := &db.TCPRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &db.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateTCP godoc
// @Summary Update TCP service.
// @Description Update TCP service
// @Tags Service Management
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body io.ServiceUpdateTCPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (admin *ServiceController) ServiceUpdateTCP(c *gin.Context) {
	params := &io.ServiceUpdateTCPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &db.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2003, err)
		return
	}

	loadBalance := &db.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IPList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	tcpRule := &db.TCPRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	accessControl := &db.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceAddGRPC godoc
// @Summary Add GRPC service.
// @Description Add GRPC service.
// @Tags Service Management
// @ID /service/service_add_grpc
// @Accept  json
// @Produce  json
// @Param body body io.ServiceAddGRPCInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_grpc [post]
func (admin *ServiceController) ServiceAddGRPC(c *gin.Context) {
	params := &io.ServiceAddGRPCInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//验证 service_name 是否被占用
	infoSearch := &db.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(c, lib.GORMDefaultPool, infoSearch); err == nil {
		middleware.ResponseError(c, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}

	//验证端口是否被占用?
	tcpRuleSearch := &db.TCPRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(c, lib.GORMDefaultPool, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &db.GRPCRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(c, lib.GORMDefaultPool, grpcRuleSearch); err == nil {
		middleware.ResponseError(c, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()
	info := &db.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	loadBalance := &db.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IPList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	grpcRule := &db.GRPCRule{
		ServiceID:       info.ID,
		Port:            params.Port,
		HeaderTransform: params.HeadTransform,
	}
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &db.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

// ServiceUpdateGRPC godoc
// @Summary Update GRPC Service
// @Description Update GRPC Service
// @Tags Service Management
// @ID /service/service_update_grpc
// @Accept  json
// @Produce  json
// @Param body body io.ServiceUpdateGRPCInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_grpc [post]
func (admin *ServiceController) ServiceUpdateGRPC(c *gin.Context) {
	params := &io.ServiceUpdateGRPCInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	//ip与权重数量一致
	if len(strings.Split(params.IPList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}

	tx := lib.GORMDefaultPool.Begin()

	service := &db.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(c, lib.GORMDefaultPool, service)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, err)
		return
	}

	loadBalance := &db.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IPList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	grpcRule := &db.GRPCRule{}
	if detail.GRPCRule != nil {
		grpcRule = detail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	//grpcRule.Port = params.Porta
	grpcRule.HeaderTransform = params.HeadTransform
	if err := grpcRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	accessControl := &db.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "")
	return
}

//ServiceDetail godoc
//@Summary Get service detail
//@Description Get service detail.
//@Tags Service Management
//@ID /service/service_detail
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=db.ServiceDetail} "success"
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
	serviceInfo := &db.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, serviceDetail)

}

//ServiceStat godoc
//@Description Get service statistics.
//@Tags Service Management
//@ID /service/service_stat
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=io.ServiceStatOutput} "success"
//@Router /service/service_stat [get]
func (service *ServiceController) ServiceStat(c *gin.Context) {
	params := &io.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//tx, err := lib.GetGormPool("default")
	//if err != nil {
	//	middleware.ResponseError(c, 2001, err)
	//	return
	//}
	//serviceInfo := &db.ServiceInfo{ID: params.ID}
	//serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	//if err != nil {
	//	middleware.ResponseError(c, 2002, err)
	//	return
	//}
	//
	//serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	//if err != nil {
	//	middleware.ResponseError(c, 2003, err)
	//	return
	//}

	todayList := []int{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yesterdayList := []int{}
	for i := 0; i < 24; i++ {
		yesterdayList = append(yesterdayList, 0)
	}

	middleware.ResponseSuccess(c, &io.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})

}
