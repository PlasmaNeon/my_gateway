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
	"time"
)

type AdminLoginController struct{}


//AdminLogin godoc
//@Summary Admin login.
//@Description Admin login.
//@Tags Admin API.
//@ID /admin_login/login
//@Accept  json
//@Produce  json
//@Param body body io.AdminLoginInput true "body"
//@Success 200 {object} middleware.Response{data=io.AdminLoginOutput} "success"
//@Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin (c *gin.Context){
	params := &io.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &db.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}


	//Set session
	sessInfo := &io.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &io.AdminLoginOutput{Token:admin.UserName}
	middleware.ResponseSuccess(c, out)
}

func AdminLoginRegister(group *gin.RouterGroup){
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLogout)
}



