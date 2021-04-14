package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
	"strings"
)

type LoadBalance struct {
	ID                     int64  `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceId              int64    `json:"service_id" gorm:"column:service_id" description:"http:0, tcp:1, grpc:2"`
	CheckMethod            int    `json:"service_name" gorm:"column:service_name" description:"Service name."`
	CheckTimeout           int    `json:"service_desc" gorm:"column:service_desc" description:"Service description."`
	CheckInterval          int    `json:"need_https" gorm:"column:need_https" description:"Service description."`
	RoundType              int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"Service description."`
	IPList                 string `json:"need_websocket" gorm:"column:need_websocket" description:"Service description."`
	WeightList             string `json:"weight_list" gorm:"column:weight_list" description:"Service description."`
	ForbidList             string `json:"forbid_list" gorm:"column:forbid_list" description:"Service description."`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"Service description."`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"Service description."`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"Service description."`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"Service description."`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) Find(c *gin.Context, tx *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	out := &LoadBalance{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *LoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

func (t *LoadBalance) GetIPListByModel () []string{
	return strings.Split(t.IPList, ",")
}