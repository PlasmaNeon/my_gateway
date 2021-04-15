package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type TCPRule struct {
	ID        int64 `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceId int64 `json:"service_id" gorm:"column:service_id" description:"Service description.""`
	Port      int64 `json:"port" gorm:"column:port" description:"Port."`
}

func (t *TCPRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TCPRule) Find(c *gin.Context, tx *gorm.DB, search *TCPRule) (*TCPRule, error) {
	out := &TCPRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *TCPRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
