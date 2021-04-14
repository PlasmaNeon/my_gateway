package io

import (
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type ServiceListInput struct {
	Info string `json:"info" form:"info" comment:"Keyword" example:"" validate:""`
	PageNo int `json:"page_no" form:"page_no" comment:"page number" example:"1" validate:"required"`
	PageSize int `json:"page_size" form:"page_size" comment:"number of entries per page" example:"20" validate:"required"`
}
//Validation process
func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}


type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	LoadType    int  `json:"load_type" form:"load_type"`
	ServiceAddr string  `json:"service_addr" form:"service_addr"`
	QPS         int64  `json:"qps" form:"qps"`
	QPD         int64  `json:"qpd" form:"qpd"`
	TotalNode   int  `json:"total_node" form:"total_node"`
}

type ServiceListOutput struct {
	Total int64 `json:"total" form:"total" comment:"Total entries"`
	List []ServiceListItemOutput `json:"list" form:"list" comment:"list number"`
}
