package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/public"
)

type HTTPRule struct {
	ID              int    `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	ServiceId       int    `json:"service_id" gorm:"column:service_id" description:"http:0, tcp:1, grpc:2"`
	RuleType        int    `json:"rule_type" gorm:"column:rule_type" description:"Rule type."`
	Rule            string `json:"rule" gorm:"column:rule" description:"Service description."`
	NeedHttps       int    `json:"need_https" gorm:"column:need_https" description:"Service description."`
	NeedStripUri    int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"Service description."`
	NeedWebsocket   int    `json:"need_websocket" gorm:"column:need_websocket" description:"Service description."`
	UrlRewrite      string `json:"url_rewrite" gorm:"column:url_rewrite" description:"Service description."`
	HeaderTransform string `json:"header_transform" gorm:"column:header_transfor" description:"Service description."`
}

func (t *HTTPRule) TableName() string {
	return "gateway_service_http_rule"
}
func (t *HTTPRule) Find(c *gin.Context, tx *gorm.DB, search *HTTPRule) (*HTTPRule, error) {
	out := &HTTPRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *HTTPRule) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
