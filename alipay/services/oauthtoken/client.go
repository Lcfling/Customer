package oauthtoken

import (
	"github.com/Lcfling/Customer/alipay/sdk"
	"github.com/Lcfling/Customer/alipay/tools"
	"time"
)
type Client struct {
	sdk.Client
	GrantType string
	Code 	string
	RefreshToken string
}

func NewClient() *Client {
	client := new(Client)
	client.Method="alipay.system.oauth.token"
	client.Charset="utf-8"
	client.Timestamp=tools.GetDateMHS(time.Now().Unix())
	client.Datamap=make(map[string]string)
	//client.GrantType
	return client
}
func (client *Client)Execute()(string,error){
	data:=client.FormatInit()
	data["grant_type"]=client.GrantType
	if client.GrantType=="authorization_code" {
		data["code"]=client.Code
	}else{
		data["refresh_token"]=client.RefreshToken
	}

	//client.Datamap=data
	for k,v:=range data{
		client.Datamap[k]=v
	}
	respose,err:=client.DoHttps()
	return respose,err
}