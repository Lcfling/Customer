package sdk

import (
	"crypto/rsa"
	"fmt"
	"github.com/Lcfling/Customer/alipay/tools"
)

type Client struct {
	AppId 	string
	Method  string
	Format 	string
	Charset string
	SignType string
	Sign 	string
	Timestamp string
	Version		string
	privateKey *rsa.PrivateKey
	GeteWay string
	Keypath string
	Datamap map[string]string
}

func (client *Client)SignWithRSA2(source string)(string,error) {
	return tools.Sha256WithRsa(source,client.privateKey)
}

func (client *Client)DoHttp()  {

}
func (client *Client)DoHttps() (string,error) {


	//client.SignWithRSA2()
	//keystr:=
	//key,err:=tools.LoadPrivateKeyWithPath(client.Keypath)
	keystr:=tools.FormatPrivateKey(client.Keypath)
	key,err:=tools.LoadPrivateKey(keystr)
	if err!=nil{
		return "",err
	}
	client.privateKey=key
	sign,err:=tools.Signer(client.Datamap,client.privateKey)

	fmt.Println("签名后字符串sign：",sign)
	if err!=nil{
		fmt.Println("签名后字符串err：",err.Error())
		return "",err
	}
	client.Datamap["sign"]=sign

	fmt.Println("data:",client.Datamap)

	return tools.CurlPost(client.Datamap,nil,client.GeteWay)

	//return ""
}

func (client *Client) SetRootSn(strPath string)error{
	sn,err:=tools.GetCertRootSn(strPath)
	if err!=nil{
		return err
	}
	client.Datamap["alipay_root_cert_sn"]=sn
	return nil
}
func (client *Client) SetAppSn(strPath string)error{

	sn,err:=tools.GetCertRootSn(strPath)
	if err!=nil{
		return err
	}
	client.Datamap["app_cert_sn"]=sn
	return nil
}

func (client *Client) LoadPrivateKeyWithPath (path string)error{
	key,err:=tools.LoadPrivateKeyWithPath(path)
	if err!=nil{
		return err
	}
	client.privateKey=key
	return nil
}
func (client *Client)FormatInit()map[string]string{
	client.GeteWay="https://openapi.alipay.com/gateway.do"
	mapdata:=make(map[string]string)
	mapdata["app_id"]=client.AppId
	mapdata["method"]=client.Method
	//mapdata["format"]=client.Format
	mapdata["charset"]=client.Charset
	mapdata["sign_type"]=client.SignType
	mapdata["timestamp"]=client.Timestamp
	mapdata["version"]=client.Version
	return mapdata
}