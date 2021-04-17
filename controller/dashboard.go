package controller

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
	"time"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	service := &DashboardController{}
	group.GET("/sum_data", service.SumData)
	group.GET("/global_flow_stat", service.GlobalFlowStat)
	group.GET("/global_services_count", service.GlobalServicesCount)
}

//SumData godoc
//@Summary Dashboard summarize data
//@Description Dashboard summarize data
//@Tags Dashboard Management
//@ID /dashboard/sum_data
//@Accept  json
//@Produce  json
//@Success 200 {object} middleware.Response{data=io.DashboardSumDataOutput} "success"
//@Router /dashboard/sum_data [get]
func (_ *DashboardController) SumData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &db.ServiceInfo{}
	_, serviceNum, err := info.PageList(c, tx, &io.ServiceListInput{
		PageSize: 1,
		PageNo:   1,
	})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	appInfo := &db.AppInfo{}
	_, appNum, err := appInfo.AppList(c, tx, &io.AppListInput{
		PageSize: 1,
		PageNo:   1,
	})

	output := &io.DashboardSumDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          appNum,
		TodayRequestNum: 0,
		CurrentQPS:      0,
	}
	middleware.ResponseSuccess(c, output)

}

//GlobalFlowStat godoc
//@Description Get service statistics.
//@Tags Dashboard Management
//@ID /dashboard/global_flow_stat
//@Accept  json
//@Produce  json
//@Success 200 {object} middleware.Response{data=io.DashboardGlobalFlowStatOutput} "success"
//@Router /dashboard/global_flow_stat [get]
func (_ *DashboardController) GlobalFlowStat(c *gin.Context) {

	//tx, err := lib.GetGormPool("default")
	//if err != nil {
	//	middleware.ResponseError(c, 2001, err)
	//	return
	//}
	//serviceInfo := &db.AppInfo{ID: params.ID}
	//serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	//if err != nil {
	//	middleware.ResponseError(c, 2002, err)
	//	return
	//}
	//
	//serviceDetail, err := serviceInfo.AppDetail(c, tx, serviceInfo)
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

	middleware.ResponseSuccess(c, &io.AppStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})

}

//GlobalServicesCount godoc
//@Description Get service statistics.
//@Tags Dashboard Management
//@ID /dashboard/global_services_count
//@Accept  json
//@Produce  json
//@Success 200 {object} middleware.Response{data=io.DashboardGlobalServiceCountOutput} "success"
//@Router /dashboard/global_services_count [get]
func (_ *DashboardController) GlobalServicesCount(c *gin.Context) {

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &db.ServiceInfo{}
	list, err := info.CountServices(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	services := []string{}
	for index, item := range list {
		list[index].Name = public.LoadTypeMap[item.LoadType]
		services = append(services, item.Name)
	}

	middleware.ResponseSuccess(c, &io.DashboardGlobalServiceCountOutput{
		Services: services,
		Count:    list,
	})

}
