package io

import (
	"github.com/gin-gonic/gin"
	"my_gateway/public"
	"time"
)


type AdminLoginInput struct{
	UserName string `json:"username" form:"username" comment:"Username" example:"admin" validate:"required,is_valid_username"`
	Password string `json:"password" form:"password" comment:"Password" example:"123456" validate:"required"`
}

//Validation process
func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct{
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""`
}

type AdminSessionInfo struct{
	ID int `json:"id"`
	UserName string `json:"username"`
	LoginTime time.Time `json:"login_time"`
}
