package controllers

import (
	"github.com/Lcfling/Customer/utils"

	//"github.com/virteman/OPMS/initial"
	"github.com/Lcfling/Customer/models/users"
	"github.com/astaxie/beego"
	"strconv"
)

//客服

type BaseController struct {
	beego.Controller
	IsLogin bool
	//UserInfo string
	Uid        int64
	UserType   int
	UserAvatar string
	UserAgent  string
}

//MicroMessenger 微信

func (this *BaseController) Prepare() {
	token := this.Ctx.Request.Header.Get("token")
	uidstr := this.Ctx.Request.Header.Get("uid")
	useragent := this.Ctx.Request.Header.Get("user-agent")

	this.UserAgent = utils.GetUserAgent(useragent)
	if token == "" || uidstr == "" {
		this.Data["json"] = map[string]interface{}{"code": 2, "message": "登录效验失败-效验数据为空", "data": nil}
		this.ServeJSON()
		return
	}
	//效验合法性
	uid, _ := strconv.ParseInt(uidstr, 10, 64)

	user, err := users.GetUser(uid)

	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 2, "message": "登录效验失败" + err.Error()}
		this.ServeJSON()
		return
	}
	if user.Token != token {
		this.Data["json"] = map[string]interface{}{"code": 2, "message": "登录效验失败"}
		this.ServeJSON()
		return
	}
	if user.UserType != 4 {
		this.Data["json"] = map[string]interface{}{"code": 2, "message": "登录用户类型错误", "data": nil}
		this.ServeJSON()
		return
	}
	this.IsLogin = true
	this.Uid = user.Id
	this.UserAvatar = user.Avatar
	this.UserType = user.UserType
}
