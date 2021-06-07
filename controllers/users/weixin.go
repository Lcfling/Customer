package users

import (
	"encoding/json"
	"fmt"
	"github.com/Lcfling/Customer/alipay/services/oauthtoken"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/cmd"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/socket"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
	wx "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"strconv"
)

//更加微信小程序code 换取票据 openid unionid、session_key
//https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
type CodeSession struct {
	controllers.IndexController
}

func (this *CodeSession) Post() {
	code := this.GetString("code")

	if this.UserAgent == "weixin" {
		config := make(map[string]interface{})
		config["appid"] = beego.AppConfig.String("wx_appid")
		config["secret"] = beego.AppConfig.String("wx_secret")
		//todo 获取配置参数
		apiUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + config["appid"].(string) + "&secret=" + config["secret"].(string) + "&js_code=" + code + "&grant_type=authorization_code"

		if code == "" {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "code参数为空："}
			this.ServeJSON()
			return
		}

		Body, err := utils.HttpGet(apiUrl)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "微信通讯失败："}
			this.ServeJSON()
			return
		}
		var mapBody map[string]interface{}
		err = json.Unmarshal([]byte(Body), &mapBody)
		fmt.Println("map", mapBody)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "微信通讯失败：map error"}
			this.ServeJSON()
			return
		}

		openid, ok := mapBody["openid"].(string)
		if !ok {
			errCode, ok := mapBody["errcode"].(int)
			if !ok {
				this.Data["json"] = map[string]interface{}{"code": 0, "message": "微信通讯失败：get errCode err"}
				this.ServeJSON()
				return
			}
			if errCode != 0 {
				this.Data["json"] = map[string]interface{}{"code": 0, "message": "微信通讯失败：" + mapBody["errmsg"].(string)}
				this.ServeJSON()
				return
			}
		}
		UserInfo, state := users.LoginByOpenid(openid)
		if !state {

			user := new(users.Users)

			user.Account = "wx" + strconv.FormatInt(utils.SnowFlakeId(), 10)
			user.LoginType = "weixin"
			user.Status = 1
			user.Openid = openid
			user.UserType = 1
			uid, token, err := users.GreateUser(user)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"code": 0, "message": "获取用户信息失败!"}
				this.ServeJSON()
			} else {
				user.Id = uid
				user.Token = token
				users.UpdateAccesstoken(user.Id, mapBody["session_key"].(string))
				this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": map[string]interface{}{"first": true, "userinfo": user}, "session_key": mapBody["session_key"].(string)}
				this.ServeJSON()
			}

		} else {
			users.UpdateAccesstoken(UserInfo.Id, mapBody["session_key"].(string))
			this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": map[string]interface{}{"first": false, "userinfo": UserInfo}, "session_key": mapBody["session_key"].(string)}
			this.ServeJSON()
		}
	} else {
		client := oauthtoken.NewClient()
		client.AppId = beego.AppConfig.String("alipay_x_appid")
		client.AppId = "2021002145616351"
		client.Keypath = beego.AppConfig.String("ali_privatekey")
		client.Keypath = "MIIEowIBAAKCAQEAoprutP5tIXd2hs71ilWL9Kjg4twseUdMJ0RxbdO/irr1g/ty6dLSHN28lk7aXWRZmJE4mUNbWrMZA/O05bWDgi15pL1xTRBvvRQ21zVbxXw4Bdzq7WYNE89C7mGO7TDxsZDkRF1WlWRrYRCdHSOkyTPHd2mhMTfL/zXNvtLMmzrDTy3pssU02vwoHbM5iRIPVzYx/r32QXJvSKo9c5OqGsKg8JlWIzuj8l53UVEKII2pnsHk4VO5qs/aVhquWXEiUPxXspTpSyrxavnwaMSi0a4jjyrO1zPUr5hTAh5e1GsbLSdCS8pRBLh5rEwcXFOEOn23GNGGInkR6P79J/s/hwIDAQABAoIBAQCXqdix+olBaNKlpI2C7I2wsn+nOWNF70lJat49aP5D4GO1KagiDaAqimsm6v9jkoC6++CFmzyvGVNgy0PT6XxyxAWssYHnNkhyXFNWYY9qYJVEaqy4prHV40BzZY1REJCuZQ1z8ncaumIpU7ynfCJsBB6s81oEtR1RuhZgQO/Ua/lqidKF0tUzsy/62jPzJoqwJG1Vn4Nb7lHpGXf5J2AIP6sM0nOzWYLYFfyP1LmNg/OpJhDWUKHJBQCOB3jlcinKLbN1ncBHxqMCJQ7xAn1ir07eInG0AVdU9vHigBWY7P9SSxzy3zwxIl5XRFXrFWltGPw6CLJGg3Igu7SRJNDZAoGBAOD05FxWzTDwzrNnJAh/3fpVboGVXjlBj4ri9y0fZU8AST33iZ7I9W9l9wAA4lS2DpS7RTVzvJ9wjQFrBWVEuH+8vIzzFQeiPSJf0Z/bEVx31IzGF6H/9qdRkOgYt/gZ4yic+hHqN+g3bPifW9GRKOtW+sRIBoZHQHGHvihraO4FAoGBALkLVYeUz6JsmUPJ6U0qT82ZQWHPnLP+VXoA5ZiNtc/UK+LqD+K/4fXfXF6NIdUDctGB50/V2T0seaDQ0N3MF7fuqoVCZy1I3FdtQcgTGmTyojzxTKKeOj6EEUOlt8Pv0ABE8QGifUbF2s2QC59ORToB5MzAMjyQkjwGt7xJpaEbAoGANlDrEqCiyr5aKlctDCBTqK4YEJHQPmLmFdLXe72o6HpZNO0f/YboPA2Sph2QiIOs4ZyWCWH4mUbDxSPiGaGOKsmXfTD0UvOJb1NTehWbC4ijeZoa+rKjC6NWKbRON0mI37WHa+vxs9AuL5nKwb8a8jf+NIZvjNyHYuIzt+63V0ECgYBCHjukY1bBiZ5F64pyKRE0vHLxORab9d+i5VkkZlY1eXFo9gtRERDzIqlFm5YgH8hR9eGp1BZ4VkDrZlGLPtamwR+q1+w38RXSI1bi33iJ42x27B1e6byUA+qLSlZcK38d6YRX+jBbLm0dEEAm3ve7X1vakT4iB+JIknnqTEJjSwKBgG7PFJv/ZjqQnjKr5QWqqjQb2MuH15yLlN8SLQ8RV91groGVPBGSLQr1pELH0oBT1wowGVRWNIyjgdvN+zCauc7SIHdliHwvPVLuRO/i1phKBTlnavrITtPjl+n9tDnDNYYrZEoBB4SGoBvbb/OWeq2uCJClqvGMut3xpA7sgBiB"
		client.SignType = "RSA2"
		client.Version = "1.0"
		client.GrantType = "authorization_code"
		client.Code = code
		//err:=client.SetAppSn("./cert/appCertPublicKey.cer")

		//err=client.SetRootSn("./cert/alipayRootCert.cer")
		resJson, err := client.Execute()
		//fmt.Println("res:",resJson)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "err:" + err.Error()}
			this.ServeJSON()
			return
		}
		var mapResult map[string]interface{}
		err = json.Unmarshal([]byte(resJson), &mapResult)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "json格式错误"}
			this.ServeJSON()
			return
		}

		res, ok := mapResult["alipay_system_oauth_token_response"].(map[string]interface{})

		fmt.Println("ers:::", mapResult)
		if !ok {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "远端错误代码"}
			this.ServeJSON()
			return
		}
		accesstoken := res["access_token"].(string)
		openid := res["user_id"].(string)

		UserInfo, state := users.LoginByOpenid(openid)
		if !state {

			user := new(users.Users)

			user.Account = "ali" + strconv.FormatInt(utils.SnowFlakeId(), 10)
			user.LoginType = "alipay"
			user.Status = 1
			user.Openid = openid
			user.Accesstoken = accesstoken
			user.UserType = 2
			uid, token, err := users.GreateUser(user)
			if err != nil {
				this.Data["json"] = map[string]interface{}{"code": 0, "message": "获取用户信息失败!"}
				this.ServeJSON()
			} else {
				user.Id = uid
				user.Token = token
				users.UpdateAccesstoken(user.Id, accesstoken)
				this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": map[string]interface{}{"first": true, "userinfo": user}, "access_token": accesstoken}
				this.ServeJSON()
			}
		} else {
			users.UpdateAccesstoken(UserInfo.Id, accesstoken)
			this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": map[string]interface{}{"first": false, "userinfo": UserInfo}, "access_token": accesstoken}
			this.ServeJSON()
		}
	}

}

type Modeltest struct {
	controllers.IndexController
}

func (this *Modeltest) Get() {
	//utils.Set(models.GetRedis(),"user_4","hshhshshshshshshsh")
	utils.RPush(models.GetRedis(), "settlement_by_day", "2021-5-12")
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "微信通讯失败：", "data": ""}
	this.ServeJSON()
}

//***微信支付回调
type WxPayNotify struct {
	controllers.IndexController
}

func (this *WxPayNotify) Post() {
	body := this.Ctx.Input.RequestBody

	apiv3Key := beego.AppConfig.String("apiv3key")

	var bodyMap map[string]interface{}
	err := json.Unmarshal(body, &bodyMap)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": "faild", "message": "数据解析失败"}
		this.ServeJSON()
		return
	}

	var resource map[string]interface{}
	resource = make(map[string]interface{})
	resource, ok := bodyMap["resource"].(map[string]interface{})

	if !ok {
		this.Data["json"] = map[string]interface{}{"code": "false", "message": "bodyMap 解析错误"}
		this.ServeJSON()
		return
		//return
	}

	jsonText, err := wx.DecryptToString(apiv3Key, resource["associated_data"].(string), resource["nonce"].(string), resource["ciphertext"].(string))

	var resultMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonText), &resultMap)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": "false", "message": "resultMap 解析错误"}
		this.ServeJSON()
		return
	}
	fmt.Println("resultMap:", resultMap)

	order_id := resultMap["out_trade_no"].(string)
	out_trade_sn := resultMap["transaction_id"].(string)
	money := int64(resultMap["amount"].(map[string]interface{})["total"].(float64))
	orderinfo, err := order.GetOrderByOrderId(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": "false", "message": "订单不存在"}
		this.ServeJSON()
		return
	}
	err = order.OrderPaid(order_id, money, out_trade_sn)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": "false", "message": err.Error()}
		this.ServeJSON()
		return
	}
	service_id, _ := service.GetService(orderinfo.StoreId)
	change := cmd.SendOrderStatus{Cmd: 5, Uid: orderinfo.Uid, Storeid: orderinfo.StoreId, Status: 1, Ordersn: order_id}
	go socket.SendMessageToPeer(service_id, change)

	this.Data["json"] = map[string]interface{}{"code": "success", "message": "成功"}
	this.ServeJSON()
	return
}
