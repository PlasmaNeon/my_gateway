package io

import (
	"github.com/gin-gonic/gin"
	"my_gateway/public"
	"time"
)

type AppAddInput struct {
	//ID            int  `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	AppID    string `json:"app_id" form:"app_id" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	Name     string `json:"name" form:"name" comment:"Service Description" example:"" validate:"required,min=1,max=255"`
	Secret   string `json:"secret" form:"secret" comment:"black_list" example:"" validate:""`
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"white_list" example:"" validate:""`
	QPS      int    `json:"qps" form:"qps" comment:"upstream_connect_time_out" example:"" validate:"min=0"`
	QPD      int    `json:"qpd" form:"qpd" comment:"upstream_header_time_out" example:"" validate:"min=0"`
}

//Validation process
func (param *AppAddInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppDeleteInput struct {
	ID int `json:"id" form:"id" comment:"Service id" example:"56" validate:"required"`
}

//Validation process
func (param *AppDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppDetailInput struct {
	ID int `json:"id" form:"id" comment:"Service id" example:"56" validate:"required"`
}

//Validation process
func (param *AppDetailInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppListInput struct {
	Info     string `json:"info" form:"info" comment:"Keyword" example:"" validate:""`
	PageNo   int    `json:"page_no" form:"page_no" comment:"page number" example:"1" validate:"required"`
	PageSize int    `json:"page_size" form:"page_size" comment:"number of entries per page" example:"20" validate:"required"`
}

//Validation process
func (param *AppListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppStatInput struct {
	ID int `json:"id" form:"id" comment:"Service id" example:"56" validate:"required"`
}

//Validation process
func (param *AppStatInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppUpdateInput struct {
	ID       int    `json:"id" form:"id" comment:"Service id" example:"56" validate:"required,min=1"`
	AppID    string `json:"app_id" form:"app_id" comment:"Service Name" example:"" validate:"required,is_valid_service_name"`
	Name     string `json:"name" form:"name" comment:"Service Description" example:"" validate:"required,min=1,max=255"`
	Secret   string `json:"secret" form:"secret" comment:"black_list" example:"" validate:""`
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"white_list" example:"" validate:""`
	QPS      int    `json:"qps" form:"qps" comment:"upstream_connect_time_out" example:"" validate:"min=0"`
	QPD      int    `json:"qpd" form:"qpd" comment:"upstream_header_time_out" example:"" validate:"min=0"`
}

//Validation process
func (param *AppUpdateInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AppListItemOutput struct {
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"Admin username."`
	Name      string    `json:"name" gorm:"column:name" description:"Salt."`
	Secret    string    `json:"secret" gorm:"column:secret" description:"Admin password."`
	WhiteIPs  string    `json:"white_ips" gorm:"column:white_ips" description:"Admin password."`
	QPD       int       `json:"qpd" gorm:"column:qpd" description:"Auto increasing primary key."`
	QPS       int       `json:"qps" gorm:"column:qps" description:"Auto increasing primary key."`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"Update time."`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"Create time"`
}
type AppListOutput struct {
	Total int                 `json:"total" form:"total" comment:"Total entries"`
	List  []AppListItemOutput `json:"list" form:"list" comment:"list number"`
}
type AppStatOutput struct {
	Today     []int `json:"today" form:"today" comment:"Today"`
	Yesterday []int `json:"yesterday" form:"yesterday" comment:"Yesterday"`
}
