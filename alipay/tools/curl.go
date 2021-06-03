package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CurlPost(data ,header map[string]string,url string) (string,error){
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("curl post error：", err) // 这里的err其实就是panic传入的内容
		}
	}()

	/*dataStr:=""
	if data!=nil&& len(data)>0{
		for k,v:=range data{
			dataStr+=k+"="+v+"&"
		}
		dataStr=strings.TrimRight(dataStr,"&")
	}*/



	fmt.Println("最终发送body：",FormatURLParam(data))
	resp, err :=http.NewRequest("POST",url,
		strings.NewReader(FormatURLParam(data)))
	resp.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")


	if header!=nil&& len(header)>0{
		for k,v:=range data{
			resp.Header.Add(k,v)
		}
	}

	respt, err := (&http.Client{}).Do(resp)
	if err != nil {
		return "",err
	}
	//defer respt.Body.Close()
	body, err := ioutil.ReadAll(respt.Body)
	if err != nil {
		return "",err
	}
	return string(body),nil
}
