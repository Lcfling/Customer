package store

import (
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"html"
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

//员工列表
type StaffList struct {
	controllers.MobileController
}

func (this *StaffList) Get() {
	list, err := order.GetStaffListByMerId(this.Uid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "failed：" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}

}

type AddStaff struct {
	controllers.MobileController
}

func (this *AddStaff) Get() {

}
func (this *AddStaff) Post() {
	account := html.EscapeString(this.GetString("account"))
	phone := html.EscapeString(this.GetString("phone"))
	name := html.EscapeString(this.GetString("name"))
	pwd := "123456"
	userType := 5

	user := new(users.Users)
	user.Account = account
	user.Pwd = utils.Md5(pwd)
	user.Nickname = name
	user.Phone = phone
	user.UserType = userType
	uid, _, err := users.GreateUser(user)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "添加员工失败"}
		this.ServeJSON()
		return
	}
	err = order.AddStaffMer(uid, this.Uid)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "添加员工失败"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "添加员工成功"}
		this.ServeJSON()
		return
	}

}

//店铺绑定员工列表
type StaffBindList struct {
	controllers.MobileController
}

func (this *StaffBindList) Get() {
	store_id, _ := this.GetInt64("storeid")
	storeInfo, err := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有店铺权限"}
		this.ServeJSON()
		return
	}
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正确的店铺"}
		this.ServeJSON()
		return
	}

	list, err := order.GetStaffBindList(store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "failed：" + err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": list}
		this.ServeJSON()
		return
	}
}

type StoreStaffBind struct {
	controllers.MobileController
}

func (this *StoreStaffBind) Post() {
	store_id, _ := this.GetInt64("storeid")
	staff_id, _ := this.GetInt64("staffid")
	if !(staff_id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "请选择正取的员工id"}
		this.ServeJSON()
		return
	}
	storeInfo, _ := order.GetStoreById(store_id)
	if storeInfo.Uid != this.Uid {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "错误的门店"}
		this.ServeJSON()
		return
	}
	err := order.BindStaff(staff_id, this.Uid, store_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "添加员工失败"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "添加员工成功"}
		this.ServeJSON()
		return
	}
}

//密码修改
type ChangePwd struct {
	controllers.StaffController
}

func (this *ChangePwd) Post() {
	oldpwd := html.EscapeString(this.GetString("oldpwd"))
	newpwd := html.EscapeString(this.GetString("newpwd"))
	userInfo, _ := users.GetUser(this.Uid)
	if utils.Md5(oldpwd) != userInfo.Pwd {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "原始密码错误"}
		this.ServeJSON()
		return
	}
	err := users.UpdatePwd(this.Uid, utils.Md5(newpwd))
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "修改失败"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "修改成功"}
		this.ServeJSON()
		return
	}
}
