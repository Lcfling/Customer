package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Curl()  {
	
}
func HttpGet(Url string) (string,error){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("HttpGet产生的致命错误：", err) // 这里的err其实就是panic传入的内容
		}
	}()
	resp, err :=   http.Get(Url)
	if err != nil {
		return "",err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}

	return string(body),err
}

func YSGetAccesstoken(appkey ,appsecret string) (string,int64,error){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("获取萤石云失败 panic产生的异常 ：", err) // 这里的err其实就是panic传入的内容
		}
	}()

	//todo 获取当前年份
	data:="appKey="+appkey+"&appSecret="+appsecret
	resp, err :=http.NewRequest("POST","https://open.ys7.com/api/lapp/token/get",

		strings.NewReader(data))
	resp.Header.Set("Content-Type", "application/x-www-form-urlencoded")



	respt, err := (&http.Client{}).Do(resp)
	if err != nil {
		return "",0,err
	}
	//defer respt.Body.Close()
	body, err := ioutil.ReadAll(respt.Body)
	if err != nil {
		return "",0,err
	}
	resultjson:=make(map[string]interface{})
	err=json.Unmarshal(body, &resultjson)


	code,ok:=resultjson["code"]
	if ok && code.(string)=="200"{
		var data map[string]interface{}
		data=resultjson["data"].(map[string]interface{})
		return data["accessToken"].(string),int64(data["expireTime"].(float64)),nil
	}else{
		return "",0,errors.New("错误代码"+code.(string))
	}
}