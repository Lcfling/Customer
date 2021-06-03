package users

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Lcfling/Customer/alipay/tools"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
)

//登录服务器  获取服务许可证
type LoginController struct {
	controllers.IndexController
}

func (this *LoginController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")
	this.SetSession("userLogin", 1)
	this.Data["json"] = map[string]interface{}{"code": 1, "info": "登录成功", "data": username + password}
	this.ServeJSON()
}

type UpdateUser struct {
	controllers.UserBaseController
}

func (this *UpdateUser) Post() {
	avatar:=this.GetString("avatar")
	nickname:=this.GetString("nickname")
	phone:=this.GetString("phone")

	userinfo,err:=users.GetUser(this.Uid)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "获取用户信息失败!"}
		this.ServeJSON()
	}
	uid:=this.Uid
	user:=users.Users{}
	if avatar!=""{
		user.Avatar=avatar
	}
	if nickname!=""{
		user.Nickname=nickname
	}
	if phone!=""{
		if this.UserAgent=="weixin"{
			iv:=this.GetString("iv")
			wx:=utils.NewWXBizDataCrypt(beego.AppConfig.String("wx_appid"),userinfo.Accesstoken)
			cipherText,err:=wx.Decrypt(phone,iv)
			if err!=nil{

				this.Data["json"]=map[string]interface{}{"code": 0, "message": "解密失败!"}
				this.ServeJSON()
				return
			}
			var phInfo utils.PhoneStruct
			err= json.Unmarshal([]byte(cipherText), &phInfo)
			if err!=nil{
				this.Data["json"]=map[string]interface{}{"code": 0, "message": "json格式化错误!"}
				this.ServeJSON()
				return
			}
			phone=phInfo.PurePhoneNumber
		}else{

			iv:=make([]byte,16)
			cipherText,_:=base64.StdEncoding.DecodeString(phone)
			key,_:=base64.StdEncoding.DecodeString(beego.AppConfig.String("ali_aes_secret"))
			jsonbyte,err:=tools.CBCDecryptIvData(cipherText,key,iv)
			if err!=nil{
				this.Data["json"]=map[string]interface{}{"code": 0, "message": "解密失败!"}
				this.ServeJSON()
				return
			}
			var phInfo utils.APhoneStruct
			err= json.Unmarshal(jsonbyte, &phInfo)
			if err!=nil{
				this.Data["json"]=map[string]interface{}{"code": 0, "message": "json格式化错误!"}
				this.ServeJSON()
				return
			}
			if !(phInfo.Code=="10000"){
				this.Data["json"]=map[string]interface{}{"code": 0, "message": "支付宝响应失败!"+phInfo.Msg}
				this.ServeJSON()
				return
			}
			phone=phInfo.Mobile
		}

		user.Phone=phone
	}

	err=users.UpdateUser(uid,user)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "更新失败!"}
		this.ServeJSON()
	}
	user,err=users.GetUser(this.Uid)
	if user.Phone==""{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "空的手机号码!"}
		this.ServeJSON()
	}
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "更新失败!"}
		this.ServeJSON()
	}else{
		this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":user}
		this.ServeJSON()
	}
}