package models

import (
	"iris-ticket/backend/database"

	"github.com/kataras/golog"
)

type Role struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Desc        string `gorm:"VARCHAR(191)"`
}

type RoleJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Desc        string `json:"desc"`
}

func GetRoleById(id uint) (role *Role, err error) {
	role = new(Role)
	role.ID = id

	if err = database.DB.First(role).Error; err != nil {
		golog.Error("GetRoleByIdErr ", err)
	}

	return
}

func GetRoleByName(name string) (role *Role, err error) {
	role = new(Role)
	role.Name = name

	if err = database.DB.First(role).Error; err != nil {
		golog.Error("GetRoleByNameErr ", err)
	}

	return
}

func DeleteRoleById(id uint) {
	u := new(Role)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		golog.Error("DeleteRoleErr ", err)
	}
}

func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Role, err error) {

	if err = database.GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		golog.Error("GetAllRoleErr ", err)
	}
	return
}

func CreateRole(aul *RoleJson) (role *Role, err error) {

	role = new(Role)
	role.Name = aul.Name
	role.DisplayName = aul.DisplayName
	role.Desc = aul.Desc

	if err = database.DB.Create(role).Error; err != nil {
		golog.Error("CreateRoleErr ", err)
	}

	return
}

func UpdateRole(rj *RoleJson, id uint) (role *Role, err error) {
	role = new(Role)
	role.ID = id

	if err = database.DB.Model(role).Updates(rj).Error; err != nil {
		golog.Error("UpdatRoleErr ", err)
	}

	return
}

func CreateSystemAdminRole(rolename string) (role *Role, err error) {
	aul := new(RoleJson)
	aul.Name = rolename
	aul.DisplayName = "超级管理员"
	aul.Desc = "超级管理员"

	role, err = GetRoleByName(aul.Name)

	if role.ID == 0 {
		golog.Info("创建角色")
		return CreateRole(aul)
	} else {
		golog.Warn("角色已存在")
		return
	}
}
