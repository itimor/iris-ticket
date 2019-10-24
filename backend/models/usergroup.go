package models

import (
	"iris-ticket/backend/database"

	"github.com/kataras/golog"
)

type UserGroup struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Desc        string `gorm:"VARCHAR(191)"`
}

type UserGroupJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Desc        string `json:"desc"`
}

func GetUserGroupById(id uint) (group *UserGroup, err error) {
	group = new(UserGroup)
	group.ID = id

	if err = database.DB.First(group).Error; err != nil {
		golog.Error("GetUserGroupByIdErr ", err)
	}

	return
}

func GetUserGroupByName(name string) (group *UserGroup, err error) {
	group = new(UserGroup)
	group.Name = name

	if err = database.DB.First(group).Error; err != nil {
		golog.Error("GetUserGroupByNameErr ", err)
	}

	return
}

func DeleteUserGroupById(id uint) {
	u := new(UserGroup)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		golog.Error("DeleteUserGroupErr ", err)
	}
}

func GetAllUserGroups(name, orderBy string, offset, limit int) (UserGroups []*UserGroup, err error) {

	if err = database.GetAll(name, orderBy, offset, limit).Find(&UserGroups).Error; err != nil {
		golog.Error("GetAllUserGroupErr ", err)
	}
	return
}

func CreateUserGroup(aul *UserGroupJson) (group *UserGroup, err error) {

	group = new(UserGroup)
	group.Name = aul.Name
	group.DisplayName = aul.DisplayName
	group.Desc = aul.Desc

	if err = database.DB.Create(group).Error; err != nil {
		golog.Error("CreateUserGroupErr ", err)
	}

	return
}

func UpdateUserGroup(rj *UserGroupJson, id uint) (group *UserGroup, err error) {
	group = new(UserGroup)
	group.ID = id

	if err = database.DB.Model(group).Updates(rj).Error; err != nil {
		golog.Error("UpdatUserGroupErr ", err)
	}

	return
}
