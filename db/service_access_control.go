package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceId         int64  `json:"service_id" gorm:"column:service_id" description:"http:0, tcp:1, grpc:2"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"Service name."`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"Service description."`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"Service description."`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"Service description."`
	ClientIpFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"Service description."`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"Service description."`
}

func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (t *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	out := &AccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *AccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
