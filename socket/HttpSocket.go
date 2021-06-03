package socket

import (
	"encoding/json"
	log "github.com/golang/glog"
)

type Request struct {
	ip     string
	url    string
	header string
	body   string
}

func (r *Request) GetHeader(data string) interface{} {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(r.header), &mapResult)
	if err != nil {

		log.Info("header get json err: ", err)
	}
	return mapResult[data]
}

func (r *Request) GetPost(data string) interface{} {
	var mapResult map[string]interface{}
	err := json.Unmarshal([]byte(r.body), &mapResult)
	if err != nil {
		log.Info("header get json err: ", err)
	}
	return mapResult[data]
}

/*type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool // whether any patterns contain hostnames
}*/
type Handler interface {
	ServeHTTP(client *Client, r *Request)
}
type HandlerFunc func(*Client, *Request)

func HandleSocket(client *Client) {
	client.Handerfunc("/data", HandlerFunc(SgBetTop3_s))
}
func SgBetTop3_s(c *Client, r *Request) {
	WWriteHttpFailed("失败", c)
}
