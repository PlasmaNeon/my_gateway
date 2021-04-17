package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/io"
	"my_gateway/public"
	"time"
)

type AppInfo struct {
	ID        int       `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"Admin username."`
	Name      string    `json:"name" gorm:"column:name" description:"Salt."`
	Secret    string    `json:"secret" gorm:"column:secret" description:"Admin password."`
	WhiteIPs  string    `json:"white_ips" gorm:"column:white_ips" description:"Admin password."`
	QPD       int       `json:"qpd" gorm:"column:qpd" description:"Auto increasing primary key."`
	QPS       int       `json:"qps" gorm:"column:qps" description:"Auto increasing primary key."`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"Update time."`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"Create time"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"Whether user is deleted."`
}

func (t *AppInfo) TableName() string {
	return "gateway_app"
}
func (t *AppInfo) Find(c *gin.Context, tx *gorm.DB, search *AppInfo) (*AppInfo, error) {
	out := &AppInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *AppInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}
func (t *AppInfo) AppList(c *gin.Context, tx *gorm.DB,
	params *io.AppListInput) ([]AppInfo, int, error) {
	total := 0
	list := []AppInfo{}
	offset := (params.PageNo - 1) * params.PageSize
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if params.Info != "" {
		query = query.Where("service_name like ? or service_desc like ?",
			"%"+params.Info+"%",
			"%"+params.Info+"%")
	}
	if err := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	query.Count(&total)
	return list, total, nil
}
