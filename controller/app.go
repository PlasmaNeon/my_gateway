package controller

import (
	"errors"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
	"time"
)

type AppController struct{}

func AppRegister(group *gin.RouterGroup) {
	service := &AppController{}
	group.GET("/app_list", service.AppList)
	group.GET("/app_detail", service.AppDetail)
	group.GET("/app_stat", service.AppStat)
	group.GET("/app_delete", service.AppDelete)
	group.POST("/app_add", service.AppAdd)
	group.POST("/app_update", service.AppUpdate)
}

//AppDetail godoc
//@Summary Get service detail
//@Description Get service detail.
//@Tags AppInfo Management
//@ID /app/app_detail
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=db.AppInfo} "success"
//@Router /AppInfo/App_detail [get]
func (service *AppController) AppDetail(c *gin.Context) {
	params := &io.AppDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	search := &db.AppInfo{ID: params.ID}
	detail, err := search.Find(c, tx, search)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	middleware.ResponseSuccess(c, detail)

}

//AppList godoc
//@Summary Show app list
//@Description Show app lists.
//@Tags AppInfo Management
//@ID /app/app_list
//@Accept  json
//@Produce  json
//@Param info query string false "Searching keyword"
//@Param page_size query int true "Entries per page"
//@Param page_no query int true "Page No."
//@Success 200 {object} middleware.Response{data=io.AppListOutput} "success"
//@Router /service/service_list [get]
func (_ *AppController) AppList(c *gin.Context) {
	params := &io.AppListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &db.AppInfo{}
	list, total, err := info.AppList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	outList := []io.AppListItemOutput{}

	for _, item := range list {
		//realQPS := 0
		//realQPD := 0
		outItem := io.AppListItemOutput{
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPs: item.WhiteIPs,
			QPD:      0,
			QPS:      0,
		}
		outList = append(outList, outItem)
	}

	out := &io.AppListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(c, out)

}

//AppDelete godoc
//@Summary Delete a service
//@Description Delete a service.
//@Tags AppInfo Management
//@ID /service/service_delete
//@Accept  json
//@Produce  json
//@Param id query string true "service id"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_delete [get]
func (_ *AppController) AppDelete(c *gin.Context) {
	params := &io.AppDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	info := &db.AppInfo{ID: params.ID}
	info, err = info.Find(c, tx, info)

	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	info.IsDelete = 1
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}

//AppAdd godoc
//@Summary Add app
//@Description  Add app
//@Tags App Management
//@ID /app/app_add
//@Accept  json
//@Produce  json
//@Param body body io.AppAddInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_add_http [post]
func (_ *AppController) AppAdd(c *gin.Context) {
	params := &io.AppAddInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	info := &db.AppInfo{
		AppID: params.AppID,
	}
	if _, err := info.Find(c, tx, info); err == nil {
		middleware.ResponseError(c, 2000, errors.New("AppID already exists"))
		return
	}
	if params.Secret == "" {
		params.Secret = public.GenerateMD5(params.AppID)
	}

	newAppInfo := &db.AppInfo{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPs: params.WhiteIPS,
		QPS:      params.QPS,
		QPD:      params.QPD,
	}
	if err := newAppInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2006, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}

//AppUpdate godoc
//@Summary Update HTTP service
//@Description  Update HTTP service
//@Tags AppInfo Management
//@ID /service/service_update_http
//@Accept  json
//@Produce  json
//@Param body body io.AppUpdateInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_update_http [post]
func (_ *AppController) AppUpdate(c *gin.Context) {
	params := &io.AppUpdateInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	info := &db.AppInfo{ID: params.ID}
	info, err = info.Find(c, tx, info)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	if params.Secret == "" {
		params.Secret = public.GenerateMD5(params.AppID)
	}
	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPs = params.WhiteIPS
	info.QPS = params.QPS
	info.QPD = params.QPD
	if err := info.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, errors.New("app update failed"))
		return
	}
	middleware.ResponseSuccess(c, "")

}

//AppStat godoc
//@Description Get service statistics.
//@Tags AppInfo Management
//@ID /service/service_stat
//@Accept  json
//@Produce  json
//@Param id query string true "app id"
//@Success 200 {object} middleware.Response{data=io.AppStatOutput} "success"
//@Router /service/service_stat [get]
func (_ *AppController) AppStat(c *gin.Context) {
	params := &io.AppStatInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

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
