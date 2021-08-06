package customer

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/alipay/services/pay"
	"github.com/Lcfling/Customer/controllers"
	"github.com/Lcfling/Customer/models/cmd"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/models/users"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/socket"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/signers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	wx "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type SubOrderController struct {
	controllers.UserBaseController
}

func (this *SubOrderController) Post() {
	sublistJson := this.GetString("proList")
	var paytype int64
	if this.UserAgent == "weixin" {
		paytype = 0
	} else {
		paytype = 1
	}

	var sublist []order.SubList

	err := json.Unmarshal([]byte(sublistJson), &sublist)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "json格式错误"}
		this.ServeJSON()
		return
	}

	token := this.GetString("dsn")
	doorinfo, _ := device.GetDiviveByToken(token)
	store_id := doorinfo.StoreId
	if !(store_id > 0) {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "商家信息不存在"}
		this.ServeJSON()
		return
	}
	//生成订单号
	order_id := utils.GetOrderSN()
	enterlog, _ := logs.GetEid(this.Uid)

	err = order.CreatOrder(sublist, store_id, this.Uid, order_id, paytype, enterlog.Id)
	//logs.UpdateOrderid(storeid,this.Uid,order_id)

	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "创建订单错误，请联系客服处理！"}
		this.ServeJSON()
		return
	}

	var mapDate map[string]string

	fmt.Println("userAgent:", this.UserAgent)
	if this.UserAgent == "weixin" {
		mapDate = CreateWxOrder(order_id, this.Openid)

		mapDate["order_id"] = order_id
	} else {
		mapDate, err = CreatAliOrder(order_id, this.Uid)
		if err != nil {
			this.Data["json"] = map[string]interface{}{"code": 0, "message": "err:" + err.Error()}
			this.ServeJSON()
			return
		}
	}
	orderSend := cmd.SendOrder{Cmd: 3, Uid: this.Uid, Storeid: store_id, Ordersn: order_id}
	service_id, _ := service.GetService(store_id)

	if service_id > 0 {
		go socket.SendMessageToPeer(service_id, orderSend)
	}

	/*if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "数据处理失败，联系管理员"}
		this.ServeJSON()
		return
	}*/

	fmt.Println("mapData:", mapDate)
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": mapDate}
	this.ServeJSON()
	return
}

type PayOrder struct {
	controllers.UserBaseController
}

func (this *PayOrder) Post() {
	order_id := this.GetString("orderid")
	orderInfo, err := order.GetOrderByOrderId(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "订单不存在"}
		this.ServeJSON()
		return
	}

	mapDate := CreateWxOrder(orderInfo.OrderId, this.Openid)

	if err == nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "数据处理失败，联系管理员"}
		this.ServeJSON()
		return
	}
	mapDate["order_id"] = order_id
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": mapDate}
	this.ServeJSON()
	return
}

type OrderPaid struct {
	controllers.IndexController
}

func (this *OrderPaid) Get() {
	order_id := this.GetString("orderid")
	m, _ := this.GetInt64("m")
	err := order.OrderPaid(order_id, m, "sdbshbahw")

	orderinfo, _ := order.GetOrderByOrderId(order_id)

	service_id, _ := service.GetService(orderinfo.StoreId)
	change := cmd.SendOrderStatus{Cmd: 5, Uid: orderinfo.Uid, Storeid: orderinfo.StoreId, Status: orderinfo.Status}
	go socket.SendMessageToPeer(service_id, change)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": err.Error()}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success"}
		this.ServeJSON()
	}
}

type OrderList struct {
	controllers.UserBaseController
}

func (this *OrderList) Get() {
	lastid, _ := this.GetInt64("lastid")
	orderlist, err := order.GetOrdersByUid(this.Uid, lastid)
	var sellarray [][]order.SellDetail
	var orderMaplist []map[string]interface{}
	for _, v := range orderlist {
		var selllist []order.SellDetail
		_, selllist, _ = order.ListByOrder(v.OrderId)
		sellarray = append(sellarray, selllist)
		//var value map[string]interface{}
		value := make(map[string]interface{})
		value, _ = utils.StoMap(v)
		Store, _ := order.GetStoreById(v.StoreId)
		value["Storename"] = Store.Name
		orderMaplist = append(orderMaplist, value)
	}
	data := map[string]interface{}{"orderlist": orderMaplist, "sellarray": sellarray}
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "没有更多的信息"}
		this.ServeJSON()
		return
	} else {
		this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": data}
		this.ServeJSON()
		return
	}
}

type OrderDetail struct {
	controllers.UserBaseController
}

func (this *OrderDetail) Get() {

	order_id := this.GetString("orderid")
	orderInfo, err := order.GetOrderByOrderId(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "订单不存在"}
		this.ServeJSON()
		return
	}
	_, selllist, err := order.ListByOrder(order_id)
	if err != nil {
		this.Data["json"] = map[string]interface{}{"code": 0, "message": "订单详情为空"}
		this.ServeJSON()
		return
	}
	data := map[string]interface{}{"orderinfo": orderInfo, "selllist": selllist}
	this.Data["json"] = map[string]interface{}{"code": 1, "message": "success", "data": data}
	this.ServeJSON()

}

/*
   Package core 微信支付api v3 go http-client 基础库，你可以使用它来创建一个client，并向微信支付发送http请求
   只需要你在初始化客户端的时候，传递credential以及validator
   credential用来生成http header中的authorization信息
   validator则用来校验回包是否被篡改
   如果http请求返回的err为nil，一般response.Body 都不为空，你可以尝试对其进行序列化
   请注意及时关闭response.Body
   注意：使用微信支付apiv3 go库需要引入相关的包，该示例代码必须引入的包名有以下信息

   "context"
   "crypto/x509"
   "fmt"
   "io/ioutil"
   "log"
   "github.com/wechatpay-apiv3/wechatpay-go/core"
   "github.com/wechatpay-apiv3/wechatpay-go/core/option"
   "github.com/wechatpay-apiv3/wechatpay-go/utils"

*/
func SetUp() (opt []option.ClientOption, err error) {
	//商户号
	mchID := beego.AppConfig.String("wechat_mchID")
	//商户证书序列号
	mchCertSerialNumber := beego.AppConfig.String("mchCertSerialNumber")
	//商户私钥文件路径
	privateKeyPath := beego.AppConfig.String("privateKeyPath")
	//平台证书文件路径
	wechatCertificatePath := beego.AppConfig.String("wechatCertificatePath")

	// 加载商户私钥
	privateKey, err := wx.LoadPrivateKeyWithPath(privateKeyPath)
	if err != nil {
		log.Printf("load private err:%s", err.Error())
		return nil, err
	}
	// 加载微信支付平台证书
	wechatPayCertificate, err := wx.LoadCertificateWithPath(wechatCertificatePath)
	if err != nil {
		log.Printf("load certificate err:%s", err)
		return nil, err
	}
	//设置header头中authorization信息
	opts := []option.ClientOption{
		option.WithMerchant(mchID, mchCertSerialNumber, privateKey),     // 设置商户相关配置
		option.WithWechatPay([]*x509.Certificate{wechatPayCertificate}), // 设置微信支付平台证书，用于校验回包信息用
		option.WithoutValidator(),
	}
	return opts, nil
}

func CreateWxOrder(order_id, openid string) map[string]string {
	// 初始化客户端
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("微信创建远端订单错误：：", err) // 这里的err其实就是panic传入的内容
		}
	}()

	mchID := beego.AppConfig.String("wechat_mchID")
	appid := beego.AppConfig.String("wx_appid")
	orderInfo, _ := order.GetOrderByOrderId(order_id)
	ctx := context.TODO()
	opts, err := SetUp()
	if err != nil {
		return nil
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("init client err:%s", err)
		return nil
	}
	//设置请求地址
	URL := "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
	//设置请求信息,此处也可以使用结构体来进行请求
	mapInfo := map[string]interface{}{
		"mchid":        mchID,
		"out_trade_no": order_id,
		"appid":        appid,
		"description":  "U云智能无人超市订单支付",
		"notify_url":   beego.AppConfig.String("wechat_notifyUrl"),
		"amount": map[string]interface{}{
			"total":    orderInfo.TotalPrice,
			"currency": "CNY",
		},
		"payer": map[string]interface{}{
			"openid": openid,
		},
	}

	// 发起请求
	response, err := client.Post(ctx, URL, mapInfo)
	if err != nil {
		log.Printf("client post err:%s", err)
		//return ""
	}

	// 校验回包内容是否有逻辑错误
	err = core.CheckResponse(response)
	if err != nil {
		log.Printf("check response err:%s", err)
		return nil
	}
	// 读取回包信息
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read response body err:%s", err)
		return nil
	}
	fmt.Println(string(body))
	var prepaymap map[string]string
	err = json.Unmarshal(body, &prepaymap)

	//开始签名
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := utils.Md5(utils.RandChar(6))
	packageStr := "prepay_id=" + prepaymap["prepay_id"]
	signBody := appid + "\n"
	signBody += timeStamp + "\n"
	signBody += nonceStr + "\n"
	signBody += packageStr + "\n"
	privateKeyPath := "./cert/apiclient_key.pem"
	privateKey, err := wx.LoadPrivateKeyWithPath(privateKeyPath)

	//Signer:=&signers.Sha256WithRSASigner{PrivateKey: privateKey, MchCertificateSerialNo: mchCertSerialNumber}
	sign, err := signers.Sha256WithRsa(signBody, privateKey)

	resultMap := make(map[string]string)
	resultMap["timeStamp"] = timeStamp
	resultMap["nonceStr"] = nonceStr
	resultMap["package"] = packageStr
	resultMap["signType"] = "RSA"
	resultMap["paySign"] = sign

	return resultMap
}

func CreatAliOrder(order_id string, uid int64) (map[string]string, error) {

	userinfo, _ := users.GetUser(uid)

	orderInfo, _ := order.GetOrderByOrderId(order_id)

	client := pay.NewClient()
	client.AppId = beego.AppConfig.String("alipay_x_appid")
	client.Keypath = beego.AppConfig.String("ali_privatekey")
	client.SignType = "RSA2"
	client.Version = "1.0"
	client.NotifyUrl = beego.AppConfig.String("alipay_notifyUrl")
	//client.AppAuthToken=userinfo.Accesstoken
	//client.Code=code
	//err:=client.SetAppSn("./cert/appCertPublicKey.cer")

	//err=client.SetRootSn("./cert/alipayRootCert.cer")

	BizContent := make(map[string]interface{})
	BizContent["out_trade_no"] = order_id
	BizContent["total_amount"] = float64(orderInfo.TotalPrice) / 100
	BizContent["buyer_id"] = userinfo.Openid
	BizContent["subject"] = "微淘云智服超市订单支付"
	BizJson, err := json.Marshal(BizContent)
	if err != nil {
		return nil, err
	}

	fmt.Println("biz::", string(BizJson))
	client.BizContent = string(BizJson)

	resJson, err := client.Execute()
	if err != nil {
		fmt.Println("err::::", err)
		return nil, err
	}
	fmt.Println("json:::", resJson)
	var mapResult map[string]interface{}
	err = json.Unmarshal([]byte(resJson), &mapResult)
	if err != nil {
		return nil, err
	}
	res, ok := mapResult["alipay_trade_create_response"].(map[string]interface{})

	if !ok {
		return nil, errors.New("date error")
	}

	back := make(map[string]string)
	back["trade_no"] = res["trade_no"].(string)
	back["order_id"] = res["out_trade_no"].(string)
	return back, nil
}
