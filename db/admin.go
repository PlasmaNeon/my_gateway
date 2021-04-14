package db

import (
	"errors"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"my_gateway/io"
	"my_gateway/public"
	"time"
)
type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"Auto increasing primary key."`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"Admin username."`
	Salt      string    `json:"salt" gorm:"column:salt" description:"Salt."`
	Password  string    `json:"password" gorm:"column:password" description:"Admin password."`
	UpdatedAt  time.Time `json:"update_at" gorm:"column:update_at" description:"Update time."`
	CreatedAt  time.Time `json:"create_at" gorm:"column:create_at" description:"Create time"`
	IsDelete  int		`json:"is_delete" gorm:"column:is_delete" description:"Whether user is deleted."`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error
}

func (t *Admin) LoginCheck (c *gin.Context, tx *gorm.DB, param *io.AdminLoginInput)(*Admin, error){
	admin_info, err := t.Find(c, tx, (&Admin{UserName: param.UserName, IsDelete: 0}))
	if err != nil{
		return nil, errors.New("User does not exist.")
	}
	salt_password := public.GenerateSaltPassword(admin_info.Salt, param.Password)
	if admin_info.Password != salt_password{
		return nil, errors.New("Wrong password.")
	}
	return admin_info, nil
}