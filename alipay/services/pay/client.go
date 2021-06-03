package pay

import (
	"github.com/Lcfling/Customer/alipay/sdk"
	"github.com/Lcfling/Customer/alipay/tools"
	"time"
)

type Client struct {
	sdk.Client
	NotifyUrl string
	AppAuthToken string
	BizContent string
}

func NewClient() *Client {
	client := new(Client)
	client.Method="alipay.trade.create"
	client.Charset="utf-8"
	client.Timestamp=tools.GetDateMHS(time.Now().Unix())
	client.Datamap=make(map[string]string)
	return client
}
func (client *Client)Execute()(string,error){
	data:=client.FormatInit()
	data["notify_url"]=client.NotifyUrl
	data["biz_content"]=client.BizContent
	//data["app_auth_token"]=client.AppAuthToken
	for k,v:=range data{
		client.Datamap[k]=v
	}
	respose,err:=client.DoHttps()
	return respose,err
}