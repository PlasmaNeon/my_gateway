package controller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"my_gateway/middleware"
	"my_gateway/public"
)

//AdminLogout godoc
//@Summary Admin logout.
//@Description Admin logout.
//@Tags Admin API.
//@ID /admin_login/logout
//@Accept  json
//@Produce  json
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /admin_login/logout [get]
func (adminlogin *AdminLoginController) AdminLogout (c *gin.Context){
	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}




