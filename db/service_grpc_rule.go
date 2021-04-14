package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type GRPCRule struct {
	ID              int64  `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceId       int64    `json:"service_id" gorm:"column:service_id" description:"Service description.""`
	Port            int    `json:"port" gorm:"column:port" description:"Port."`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"Service description."`
}

func (t *GRPCRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *GRPCRule) Find(c *gin.Context, tx *gorm.DB, search *GRPCRule) (*GRPCRule, error) {
	out := &GRPCRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *GRPCRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
