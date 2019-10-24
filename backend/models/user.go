package models

import (
	"iris-ticket/backend/config"
	"iris-ticket/backend/database"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
)

type User struct {
	gorm.Model
	UserGroupID uint

	Name      string    `gorm:"type:varchar(50);not null"`
	Username  string    `gorm:"type:varchar(50);not null;unique"`
	Password  string    `gorm:"type:varchar(200);not null"`
	Avatar    string    `gorm:"default:'https://apic.douyucdn.cn/upload/avanew/face/201709/04/01/95a344efd1141fd073397fa78cf952ae_big.jpg'"`
	Roles     []*Role   `gorm:"many2many:user_roles;"`
	UserGroup UserGroup `gorm:"foreignkey:UserGroupID"`
	IsActive  bool      `gorm:"default:true"`
}

type UserJson struct {
	Name        string `json:"name" validate:"required,gte=2,lte=50"`
	Username    string `json:"username" validate:"required,gte=2,lte=50"`
	Password    string `json:"password" validate:"required,gte=8,lte=200"`
	Avatar      string `json:"avatar" validate:"required,gte=2,lte=200"`
	UserGroupID uint
	Roles       []string `json:"roles" validate:"required"`
}

type UserPassword struct {
	Password string `json:"password" validate:"required,gte=8,lte=200"`
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) (user *User, err error) {
	user = new(User)
	user.ID = id

	if err = database.DB.Preload("UserGroup").Preload("Roles").First(user).Error; err != nil {
		golog.Error("GetUserByIdErr ", err)
	}

	return
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) (user *User, err error) {
	user = new(User)
	user.Username = username
	if err := database.DB.Preload("UserGroup").Preload("Roles").First(user).Error; err != nil {
		golog.Error("GetUserByUserNameErr ", err)
	}

	return
}

/**
 * 获取所有的账号
 * @method GetAllUser
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllUsers(name, orderBy string, offset, limit int) (users []*User) {
	if err := database.GetAll(name, orderBy, offset, limit).Preload("UserGroup").Preload("Roles").Find(&users).Error; err != nil {
		golog.Error("GetAllUserErr ", err)
	}
	return
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	u := new(User)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		golog.Error("DeleteUserByIdErr ", err)
	}
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 */
func CreateUser(aul *UserJson) (user *User, err error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Name = aul.Name
	user.Username = aul.Username
	user.Password = string(hash)
	user.Avatar = aul.Avatar
	user.UserGroupID = aul.UserGroupID

	if err := database.DB.Create(user).Error; err != nil {
		golog.Error("CreateUserErr ", err)
	}

	roles := []Role{}
	database.DB.Where("name in (?)", aul.Roles).Find(&roles)
	if err := database.DB.Model(&user).Association("Roles").Append(roles).Error; err != nil {
		golog.Error("AppendRolesErr ", err)
	}

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} id int    [description]
 */
func UpdateUser(uj *UserJson, id uint) (user *User, err error) {
	user, _ = GetUserById(id)
	user.Username = uj.Username
	user.Avatar = uj.Avatar

	if err = database.DB.Model(user).Updates(uj).Error; err != nil {
		golog.Error("UpdateUserErr ", err)
	}

	roles := []Role{}
	database.DB.Where("name in (?)", uj.Roles).Find(&roles)
	if err := database.DB.Model(&user).Association("Roles").Replace(roles).Error; err != nil {
		golog.Error("AppendRolesErr ", err)
	}

	return
}

/**
 * 更新密码
 * @method UpdateUserPassword
 * @param  {[type]} password string [description]
 * @param  {[type]} id int    [description]
 */
func UpdateUserPassword(password string, id uint) (user *User, err error) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)

	user = new(User)
	user.ID = id
	user.Password = string(hash)

	err = database.DB.Model(user).Updates(user).Error
	if err != nil {
		golog.Error("UpdateUserPasswordErr ", err)
	}
	return
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}  username string [description]
 */
func UserAdminCheckLogin(username string) User {
	u := User{}
	if err := database.DB.Where("username = ?", username).First(&u).Error; err != nil {
		golog.Error("UserAdminCheckLoginErr ", err)
	}
	return u
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "user is not exist"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			expireTime := time.Now().Add(time.Hour * time.Duration(config.Conf.Get("jwt.timeout").(int64))).Unix()
			jwtSecret := config.Conf.Get("jwt.secert").(string)
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = expireTime
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			Tokenstring, err := token.SignedString([]byte(jwtSecret))

			if err != nil {
				msg = err.Error()
				return
			}

			oauthToken := new(OauthToken)
			oauthToken.Token = Tokenstring
			oauthToken.UserId = user.ID
			oauthToken.Secret = jwtSecret
			oauthToken.Revoked = false
			oauthToken.ExpressIn = expireTime
			oauthToken.CreatedAt = time.Now()
			response = oauthToken.OauthTokenCreate()
			status = true
			msg = "success"

			return
		} else {
			msg = "error"
			return
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(userId uint) string {
	ot := UpdateOauthTokenByUserId(userId)
	return ot.Secret
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(aul *UserJson) (user *User, err error) {
	user, err = GetUserByUserName(aul.Username)

	if user.ID == 0 {
		golog.Info("创建账号")
		return CreateUser(aul)
	} else {
		golog.Warn("账号已存在")
		return
	}
}
