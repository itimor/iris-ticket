package sys

import (
    "iris-ticket/backend/app/models/basemodel"
	"iris-ticket/backend/app/models/db"

    "time"
	
    "github.com/kataras/golog"
	"github.com/jinzhu/gorm"
)

type OauthToken struct {
	basemodel.Model
	Token     string `gorm:"not null default '' comment('Token') VARCHAR(191)"`
	UserId    uint64 `gorm:"not null default '' comment('UserId') VARCHAR(191)"`
	Secret    string `gorm:"not null default '' comment('Secret') VARCHAR(191)"`
	ExpressIn int64  `gorm:"not null default 0 comment('是否是标准库') BIGINT(20)"`
	Revoked   bool
}

type Token struct {
	Token string `json:"access_token"`
}

// 表名
func (OauthToken) TableName() string {
	return TableName("oauth_token")
}

// 添加前
func (m *OauthToken) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

// 更新前
func (m *OauthToken) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = time.Now()
	return nil
}

// 删除角色及关联数据
func (OauthToken) Delete(roleids []uint64) error {
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
	if err := tx.Where("id in (?)", roleids).Delete(&OauthToken{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

/**
 * oauth_token
 * @method OauthTokenCreate
 */
func (ot *OauthToken) OauthTokenCreate() (response Token) {
	db.DB.Create(ot)
	response = Token{ot.Token}

	return
}

/**
 * 通过 token 获取 access_token 记录
 * @method GetOauthTokenByToken
 * @param  {[type]}       token string [description]
 */
func GetOauthTokenByToken(token string) (ot *OauthToken) {
	ot = new(OauthToken)
	db.DB.Where("token =  ?", token).First(&ot)
	return
}

/**
 * 通过 user_id 更新 oauth_token 记录
 * @method UpdateOauthTokenByUserId
 *@param  {[type]}       user  *OauthToken [description]
 */
func UpdateOauthTokenByUserId(userId uint) (ot *OauthToken) {
	ot = new(OauthToken)
	db.DB.Model(ot).Where("revoked = ?", false).Where("user_id = ?", userId).Updates(map[string]interface{}{"revoked": true})

	return
}

/**
 * 删除token
 */
func DeleteRequestTokenByToken(token *OauthToken) (err error) {
	if err := db.DB.Delete(token).Error; err != nil {
		golog.Error("DeleteRequestTokenByTokenErr: ", err)
	}
	return
}
