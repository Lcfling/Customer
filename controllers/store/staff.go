package store

import (
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"strconv"
)

type StaffLogin struct {
	controllers.IndexController
}

func (this *StaffLogin) Post() {
	account := this.GetString("account")
	passwrod := this.GetString("password")
	if passwrod == "" || account == "" {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "账号密码不能为空"}
		this.ServeJSON()
		return
	}
	user, err := users.GetUserByAccount(account)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "用户不存在"}
		this.ServeJSON()
		return
	}
	if user.Pwd != utils.Md5(passwrod) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "密码错误"}
		this.ServeJSON()
		return
	}
	if user.UserType != 4 {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有登录权限"}
		this.ServeJSON()
		return
	}

	token := utils.Md5(strconv.FormatInt(utils.SnowFlakeId(), 10))
	err = users.UpdateToken(user.Id, token)
	user.Token = token
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "登录失败：" + err.Error()}
		this.ServeJSON()
		return
	} else {

		this.Data["json"] = map[string]interface{}{"code": 1, "message": "登录成功", "data": user}
		this.ServeJSON()
		return
	}
}
