package store

import (
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"html"
	"strconv"
	"strings"
)

type Login struct {
	controllers.IndexController
}

func (this *Login) Post() {
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
	if user.UserType != 3 {
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

type BillLists struct {
	controllers.MobileController
}

func (this *BillLists) Get() {
	lastid, _ := this.GetInt64("lastid")
	list, err := users.BillList(this.Uid, lastid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有更多的信息"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}

type Finance struct {
	controllers.MobileController
}

func (this *Finance) Post() {
	mod, _ := this.GetInt64("type")
	if !(mod > 0 && mod < 5) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的信息类型"}
		this.ServeJSON()
		return
	}
	name := html.EscapeString(this.GetString("name"))
	bankname := html.EscapeString(this.GetString("bankname"))
	code := html.EscapeString(this.GetString("code"))

	var pro users.Finance
	pro.Uid = this.Uid
	pro.Bankname = bankname
	pro.Code = code
	pro.Type = mod
	pro.Name = name
	id, err := users.AddFinance(pro)
	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "入库信息成功", "data": fmt.Sprintf("%d", id)}
	} else {
		if strings.Contains(err.Error(), "Error 1062") {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "银行卡重复"}
		} else {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "入库信息失败：" + err.Error()}
		}
	}
	this.ServeJSON()
}
func (this *Finance) Get() {
	list, err := users.GetFinanceList(this.Uid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "获取失败：" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "获取成功", "data": list}
		this.ServeJSON()
		return
	}
}

type UserInfo struct {
	controllers.StaffController
}

func (this *UserInfo) Get() {
	u, err := users.GetUser(this.Uid)
	u.Token = ""
	u.Accesstoken = ""
	u.Pwd = ""
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "获取失败：" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "获取成功", "data": u}
		this.ServeJSON()
		return
	}
}
