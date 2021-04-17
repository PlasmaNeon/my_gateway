package db

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/io"
	"my_gateway/public"
	"time"
)

type ServiceInfo struct {
	ID          int       `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"http:0, tcp:1, grpc:2"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"Service name."`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"Service description."`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"Update time."`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"Create time"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"Whether user is deleted."`
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB,
	search *ServiceInfo) (*ServiceDetail, error) {
	var err error
	if search.ServiceName == "" {
		info, err := t.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule := &GRPCRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	httpRule := &HTTPRule{ServiceId: search.ID}
	httpRule, err = httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule := &TCPRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		AccessControl: accessControl,
		GRPCRule:      grpcRule,
		HTTPRule:      httpRule,
		LoadBalance:   loadBalance,
		TCPRule:       tcpRule,
	}
	return detail, nil
}

func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB,
	param *io.ServiceListInput) ([]ServiceInfo, int, error) {

	total := 0
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("service_name like ? or service_desc like ?",
			"%"+param.Info+"%",
			"%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	query.Count(&total)
	return list, total, nil
}

func (t *ServiceInfo) CountServices(c *gin.Context, tx *gorm.DB) ([]io.ServiceCount, error) {

	list := []io.ServiceCount{}
	query := tx.SetCtx(public.GetGinTraceContext(c))
	err := query.Table(t.TableName()).Where("is_delete=0").
		Select("load_type, count(*) as count").Group("load_type").Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
