package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

type GetJsConfig struct {
	controllers.IndexController
}

func (this *GetJsConfig)Get(){


	url:=this.GetString("url")
	resp,_:=utils.Get(models.GetRedis(),"wxjsticket")
	ticket:=resp.(string)
	if ticket==""{
		access_token,err:=Access_token()
		if err!=nil{
			this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误代码50001"+err.Error()}
			this.ServeJSON()
			return
		}
		ticket,err=GetjsTicket(access_token)
		if err!=nil{
			this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误代码50001"+err.Error()}
			this.ServeJSON()
			return
		}
	}

	noncestr:=utils.RandChar(12)
	timestamp:=strconv.FormatInt(time.Now().Unix(),10)

	sign:=Signer(ticket,noncestr,timestamp,url)

	var config map[string]string

	config=make(map[string]string)
	config["appid"]=beego.AppConfig.String("wechat_appid")
	config["noncestr"]=noncestr
	config["timestamp"]=timestamp
	config["sign"]=sign

	jsonText,err:=json.Marshal(config)
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误代码10001"}
		this.ServeJSON()
		return
	}
	err=utils.SetEx(models.GetRedis(),"wxjsconfig",string(jsonText),"7200")
	if err!=nil{
		this.Data["json"]=map[string]interface{}{"code": 0, "message": "错误代码20001"}
		this.ServeJSON()
		return
	}
	this.Data["json"]=map[string]interface{}{"code": 1, "message": "success","data":config}
	this.ServeJSON()
	return
}

func Signer(jsapi_ticket,noncestr,timestamp,url string) string {
	str:="jsapi_ticket="+jsapi_ticket+"&noncestr="+noncestr+"&timestamp="+timestamp+"&url="+url
	str=utils.SHA1(str)
	return str
}
func GetjsTicket(accesstoken string)(string,error){

	res,err:=utils.Get(models.GetRedis(),"wxjsticket")

	if res.(string)!=""{
		return res.(string),nil
	}
	url:="https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token="+accesstoken+"&type=jsapi"
	Body,err:=utils.HttpGet(url)
	if err!=nil{
		return "",err
	}
	var mapBody map[string]interface{}
	err= json.Unmarshal([]byte(Body), &mapBody)
	if err!=nil{
		return "",errors.New("错误代码：10001")
	}
	errCode,ok:=mapBody["errcode"].(float64)
	if !ok{
		return "",errors.New("错误代码：50002")
	}
	if int64(errCode)!=0{
		return "",errors.New("错误代码：50002")
	}
	err=utils.SetEx(models.GetRedis(),"wxjsticket",mapBody["ticket"].(string),strconv.FormatInt(int64(mapBody["expires_in"].(float64)),10))
	if err!=nil{
		return "",errors.New("错误代码：20001")
	}
	return mapBody["ticket"].(string),nil
}

func Access_token() (string,error) {
	//https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
	apiUrl:="https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+beego.AppConfig.String("wechat_appid")+"&secret="+beego.AppConfig.String("wechat_appsecret")
	Body,err:=utils.HttpGet(apiUrl)
	if err!=nil{
		return "",err
	}
	var mapBody map[string]interface{}
	err= json.Unmarshal([]byte(Body), &mapBody)

	_,ok:=mapBody["access_token"].(string)

	fmt.Println(mapBody)
	if !ok{
		return "",errors.New("获取accesstoken失败")
	}
	return mapBody["access_token"].(string),nil
}