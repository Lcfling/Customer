package socket

import ()
import "encoding/json"

func WWriteHttpError(status int, err string, w *Client) {

	obj := make(map[string]interface{})
	meta := make(map[string]interface{})
	meta["code"] = status
	meta["message"] = err
	obj["meta"] = meta
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}

func WWriteHttpFailed(info string, w *Client) {
	meta := make(map[string]interface{})
	meta["status"] = 0
	meta["info"] = info
	meta["data"] = nil
	b, _ := json.Marshal(meta)
	w.HttpWrite(string(b))
}

func WWriteHttpObj(data map[string]interface{}, info string, w *Client) {

	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}

func WWriteHttpFile(data map[string]interface{}, w *Client) {

	b, _ := json.Marshal(data)
	w.HttpWrite(string(b))
}
func WWriteHttpOutLogin(w *Client) {

	obj := make(map[string]interface{})
	obj["data"] = nil
	obj["info"] = "登录状态验证失败"
	obj["status"] = 2
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}

func WWriteHttpEmpty(w *Client) {

	obj := make(map[string]interface{})
	obj["data"] = nil
	obj["info"] = "空列表数据"
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}

func WWriteHttpArray(data []map[string]interface{}, info string, w *Client) {

	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}

func WWriteHttpArrays(data []map[string]map[string]interface{}, info string, w *Client) {
	obj := make(map[string]interface{})
	obj["data"] = data
	obj["info"] = info
	obj["status"] = 1
	b, _ := json.Marshal(obj)
	w.HttpWrite(string(b))
}
