package socket

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
	"html"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetRandomSalt() string {
	return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成随机字符串
func GetRandomIntString(lens int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Info(err.Error())
		return err
	} else {
		return nil
	}
}

func Checkdir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func GetDistance(lat1, lng1, lat2, lng2 string) string {
	s := GetDistanceNone(lat1, lng1, lat2, lng2)
	s = s / 10000
	var dist string
	if s < 1 {
		s = Decimal(s * 1000)
		dist = fmt.Sprintf("%.2fm", s)
	} else {
		dist = fmt.Sprintf("%.2fkm", s)
	}

	return dist
}
func GetDistanceNone(lat11, lng11, lat22, lng22 string) float64 {
	lat1, _ := strconv.ParseFloat(lat11, 64)
	lng1, _ := strconv.ParseFloat(lng11, 64)
	lat2, _ := strconv.ParseFloat(lat22, 64)
	lng2, _ := strconv.ParseFloat(lng22, 64)
	radius := 6378137.00 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func SendCode(mobile string, sign string) (bool, string) {

	conn := redis_pool.Get()
	key := fmt.Sprintf("Code_%s_%s", sign, mobile)
	keyTime := fmt.Sprintf("Time_%s_%s", sign, mobile)

	lastSendTime, err := redis.Int64(conn.Do("GET", keyTime))
	if err != nil && lastSendTime != 0 {
		log.Info(err, "last:", lastSendTime)
		return false, "系统错误！"
	}

	if (time.Now().Unix() - lastSendTime) < 60 {
		return false, "一分钟之后再次尝试！"
	}

	mobileCode := GetRandomIntString(6)
	_, err = conn.Do("SET", key, mobileCode)
	if err != nil {
		log.Info(err)
		return false, "系统错误！"
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAItiitVSV4yRgr", "xBhAurjJKQHRpMjDqgD1ztSitp0rRX")
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = mobile
	request.SignName = "晶晶的IM"
	request.TemplateCode = "SMS_172735374"
	request.TemplateParam = fmt.Sprintf("{'code':%s}", mobileCode)

	response, err := client.SendSms(request)
	if err != nil {
		log.Info(err)
	}
	log.Info("response is %#v\n", response)
	//发送成功记录发送时间

	_, err = conn.Do("SET", keyTime, time.Now().Unix())
	if err != nil {
		log.Info(err)
		return false, "系统错误！"
	}
	return true, ""
}

func GetCode(mobile string, sign string) string {
	conn := redis_pool.Get()
	key := fmt.Sprintf("Code_%s_%s", sign, mobile)

	Code, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Info(err)
		return ""
	}
	return Code
}

/*
链接数据库
*/
func Mysql() *sql.DB {

	db, err := sql.Open("mysql", config.mysqldb_datasource)
	if err != nil {
		log.Info("mysql connect error")
	}
	return db
}

//set
func SetRedis(k string, v string) bool {
	//newRedis := NewOrbRedisPool()
	conn := redis_pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", k, v)
	if err != nil {
		log.Info("redis set error", err)
		return false
	}
	return true
}

//get
func GetRedis(k string) (string, bool) {
	var val string
	//newRedis := NewOrbRedisPool()
	conn := redis_pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", k))
	if err != nil {
		log.Info("redis get error", err)
		return "", false
	}
	return val, true
}

//del
func DelRedis(redis_pool *redis.Pool, k string) bool {
	conn := redis_pool.Get()

	defer conn.Close()

	_, err := redis.String(conn.Do("DEL", k))
	if err != nil {
		log.Info("redis get error", err)
		return false
	}
	return true
}

//INCR
func IncrRedis(k string) bool {
	conn := redis_pool.Get()
	defer conn.Close()
	_, err := conn.Do("INCR", k)
	if err != nil {
		log.Info("redis set error", err)
		return false
	}
	return true
}

//lh struct->json string




func DoHttpsss(host string, r *http.Request) {

	cli := &http.Client{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("io.ReadFull(r.Body, body) ", err.Error())
		log.Info("io.ReadFull(r.Body, body) ", err.Error())
		//return,没有数据也是可以的，不需要直接结束
	}
	fmt.Println("req count :", len(body), "\n")
	fmt.Println("body=", string(body))
	//fmt.Print(len(body))
	rreqUrl := "http://" + r.Host + r.URL.String()

	//reqUrl := r.URL.String()
	fmt.Println("url=", rreqUrl)

	reqUrl := ChangeHost(rreqUrl, host)
	fmt.Println("res url=", reqUrl)

	req, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(string(body)))
	if err != nil {
		fmt.Print("http.NewRequest ", err.Error())
		return
	}

	//用遍历header实现完整复制
	//contentType := r.Header.Get("Content-Type")
	//req.Header.Set("Content-Type", contentType)

	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}
	res, err := cli.Do(req)
	if err != nil {
		fmt.Print("cli.Do(req) ", err.Error())
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print("err body ::", err.Error())
		return
	}
	fmt.Print("数据处理完毕 body ", string(body))
	// n, err = io.ReadFull(res.Body, body)
	// if err != nil {
	//     fmt.Print("io.ReadFull(res.Body, body) ", err.Error())
	//     return
	// }
	//fmt.Print("count body bytes: ", n, "\n")

	/*for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}
	io.Copy(w, res.Body)*/

	//这样复制对大小控制较差，不建议。用copy即可
	// io.WriteString(w, string(body[:n]))
	// fmt.Print(string(body))
}

func DoHttp(host string, r *http.Request, boot_num string, pave_num string, game_num string) {

	cli := &http.Client{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("io.ReadFull(r.Body, body) ", err.Error())
		log.Info("io.ReadFull(r.Body, body) ", err.Error())
		//return,没有数据也是可以的，不需要直接结束
	}
	fmt.Println("req count :", len(body), "\n")
	//fmt.Print(len(body))
	rreqUrl := "http://" + r.Host + r.URL.String()

	//reqUrl := r.URL.String()

	reqUrl := ChangeHost(rreqUrl, host)

	req, err := http.NewRequest(r.Method, reqUrl, strings.NewReader("boot_num="+boot_num+"&pave_num="+pave_num+"&game_num="+game_num))
	if err != nil {
		fmt.Print("http.NewRequest ", err.Error())
		return
	}

	//用遍历header实现完整复制
	//contentType := r.Header.Get("Content-Type")
	//req.Header.Set("Content-Type", contentType)

	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}
	res, err := cli.Do(req)
	if err != nil {
		fmt.Print("cli.Do(req) ", err.Error())
		return
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Print("err body ::", err.Error())
		return
	}
	fmt.Print("http result ", string(body))
	// n, err = io.ReadFull(res.Body, body)
	// if err != nil {
	//     fmt.Print("io.ReadFull(res.Body, body) ", err.Error())
	//     return
	// }
	//fmt.Print("count body bytes: ", n, "\n")

	/*for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}
	io.Copy(w, res.Body)*/

	//这样复制对大小控制较差，不建议。用copy即可
	// io.WriteString(w, string(body[:n]))
	// fmt.Print(string(body))
}

func ChangeHost(Url string, host string) string {
	rawUrl := Url
	changeHost := host
	newUrl, _ := url.Parse(rawUrl)
	newUrl.Host = changeHost + ":" + newUrl.Port()
	log.Info(newUrl)
	return newUrl.String()
}

//使用粗函数 后期可以过滤更多字段 方便维护
func FilteredHtmlsql(out string) string {
	out = html.EscapeString(out)
	return out
}
