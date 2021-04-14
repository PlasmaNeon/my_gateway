package io

import (
	"github.com/gin-gonic/gin"
	"my_gateway/public"
	"time"
)

type AdminInfoOutput struct {
	ID        int       `json:"id"`
	Name      string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
	Avatar    string    `json:"avatar"`
	Intro     string    `json:"intro"`
	Roles     []string  `json:"roles"`
}

type ChangePwdInput struct{
	Password string `json:"password" form:"password" comment:"Password" example:"123456" validate:"required"`
}

//Validation process
func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}