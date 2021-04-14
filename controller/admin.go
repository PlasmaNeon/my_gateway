package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_gateway/db"
	"my_gateway/io"
	"my_gateway/middleware"
	"my_gateway/public"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup){
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

//AdminInfo godoc
//@Summary Admin info.
//@Description Admin info.
//@Tags Admin API.
//@ID /admin/admin_info
//@Accept  json
//@Produce  json
//@Success 200 {object} middleware.Response{data=io.AdminInfoOutput} "success"
//@Router /admin/admin_info [get]
func (adminlogin *AdminController) AdminInfo (c *gin.Context){
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey).(string)
	adminSessInfo := &io.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo), adminSessInfo); err != nil{
		middleware.ResponseError(c, 2000, err)
		return
	}

	out := &io.AdminInfoOutput{
		ID: adminSessInfo.ID,
		Name: adminSessInfo.UserName,
		LoginTime: adminSessInfo.LoginTime,
		Avatar: "",
		Intro: "Super admin",
		Roles: []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}


//ChangePwd godoc
//@Summary Change admin password.
//@Description Change admin password.
//@Tags Admin API.
//@ID /admin/change_pwd
//@Accept  json
//@Produce  json
//@Param body body io.ChangePwdInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin/change_pwd [post]
func (adminlogin *AdminController) ChangePwd (c *gin.Context){
	params := &io.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey).(string)
	adminSessInfo := &io.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(sessInfo), adminSessInfo); err != nil{
		middleware.ResponseError(c, 2000, err)
		return
	}

	// Get admin info from database.
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	adminInfo := &db.Admin{}
	adminInfo, err = adminInfo.Find(c, tx,
		&db.Admin{UserName: adminSessInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//Generate new salt password
	saltPassword := public.GenerateSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	// Save to database
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}





