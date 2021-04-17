package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
	"strings"
)

type LoadBalance struct {
	ID                     int    `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceID              int    `json:"service_id" gorm:"column:service_id" description:"http:0, tcp:1, grpc:2"`
	CheckMethod            int    `json:"check_method" gorm:"column:check_method" description:"Service name."`
	CheckTimeout           int    `json:"check_timeout" gorm:"column:check_timeout" description:"Service description."`
	CheckInterval          int    `json:"check_interval" gorm:"column:check_interval" description:"Service description."`
	RoundType              int    `json:"round_type" gorm:"column:round_type" description:"Service description."`
	IpList                 string `json:"ip_list" gorm:"column:ip_list" description:"Service description."`
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

func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}
