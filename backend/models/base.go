package models

import (
	"iris-ticket/backend/config"
)

// 初始化系统 账号 权限 角色
func CreateSystemData(env string) {

	rolename := config.Conf.Get(env + ".role").(string)
	role, _ := CreateSystemAdminRole(rolename)

	if role.ID != 0 {
		aul := new(UserJson)
		aul.Name = config.Conf.Get(env + ".user").(string)
		aul.Username = config.Conf.Get(env + ".user").(string)
		aul.Password = config.Conf.Get(env + ".pass").(string)
		aul.Roles = []string{rolename}
		CreateSystemAdmin(aul)
	}
}
