package controller

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup){
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)

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
func (service *ServiceController) ServiceList (c *gin.Context){
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
	if err != nil{
		middleware.ResponseError(c, 2002, err)
		return
	}

	outList := []io.ServiceListItemOutput{}

	for _, listItem := range list{
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil{
			middleware.ResponseError(c, 2003, err)
			return
		}
		serviceAddr := "unknown"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRulePrefixURL{
			if serviceDetail.HTTPRule.NeedHttps == 0 {
				serviceAddr = clusterIP + clusterPort + serviceDetail.HTTPRule.Rule
			} else {
				serviceAddr = clusterIP + clusterSSLPort + serviceDetail.HTTPRule.Rule
			}
		}

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain{
			serviceAddr = serviceDetail.HTTPRule.Rule
		}

		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = clusterIP + string(rune(serviceDetail.TCPRule.Port))
		}

		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = clusterIP + string(rune(serviceDetail.GRPCRule.Port))
		}

		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		outItem := io.ServiceListItemOutput{
			ID: listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			QPD: 0,
			QPS: 0,
			TotalNode: len(ipList),

		}
		outList = append(outList, outItem)
	}

	out := &io.ServiceListOutput{
		Total: total,
		List: outList,
	}
	middleware.ResponseSuccess(c, out)

}