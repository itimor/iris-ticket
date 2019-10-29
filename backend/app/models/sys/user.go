package sys

import (
	"time"

	"iris-ticket/backend/app/models/basemodel"
	"iris-ticket/backend/app/models/db"

	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
)

// 后台用户
type User struct {
	basemodel.Model
	Memo     string `gorm:"column:memo;size:64;" json:"memo" form:"memo"`                                                          //备注
	UserName string `gorm:"column:user_name;size:32;unique_index:uk_admins_user_name;not null;" json:"user_name" form:"user_name"` // 用户名
	RealName string `gorm:"column:real_name;size:32;" json:"real_name" form:"real_name"`                                           // 真实姓名
	Password string `gorm:"column:password;type:char(32);not null;" json:"password" form:"password"`                               // 密码(sha1(md5(明文))加密)
	Email    string `gorm:"column:email;size:64;" json:"email" form:"email"`                                                       // 邮箱
	Avatar   string `gorm:"default:'https://apic.douyucdn.cn/upload/avanew/face/201709/04/01/95a344efd1141fd073397fa78cf952ae_big.jpg'" json:"avatar" form:"avatar"`
	Status   uint8  `gorm:"column:status;type:tinyint(1);not null;" json:"status" form:"status"` // 状态(1:正常 2:未激活 3:暂停使用)
}

// 表名
func (User) TableName() string {
	return TableName("user")
}

// 添加前
func (m *User) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *User) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 删除用户及关联数据
func (User) Delete(userids []uint64) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("id in (?)", userids).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("user_id in (?)", userids).Delete(&UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetUserById(id uint64) (user *User, err error) {
	user = new(User)
	user.ID = id

	if err = db.DB.First(user).Error; err != nil {
		golog.Error("GetUserByIdErr ", err)
	}

	return
}
