package routers

//import "C"
import (
	"github.com/Lcfling/Customer/controllers/customer"
	"github.com/Lcfling/Customer/controllers/service"
	"github.com/Lcfling/Customer/controllers/store"
	"github.com/Lcfling/Customer/controllers/users"
	"github.com/astaxie/beego"
)

func init() {

	//beego.Router("/", &users.MainController{})

	//websocket
	beego.Router("/test", &users.Modeltest{})
	beego.Router("/testmq", &service.Mqstatus{})

	beego.Router("/opendoor", &service.OpenDoor{})

	//用户

	//文件管理
	/*客服接口 ----------  start----------*/
	beego.Router("/service/version", &service.Version{})             //登录
	beego.Router("/service/login", &service.ServiceLogin{})          //登录
	beego.Router("/service/device", &service.Device{})               //门禁列表
	beego.Router("/service/videos", &service.GetVideo{})             //摄像头列表
	beego.Router("/service/closestore", &service.CloserController{}) //关闭店铺
	beego.Router("/service/opendoor", &service.OpenDoor{})           //开门

	beego.Router("/service/leave", &service.LeaveMoment{})             //客服暂时离开
	beego.Router("/service/logout", &service.LoginOut{})               //登出
	beego.Router("/service/userinfo", &service.UserInfo{})             //用户信息
	beego.Router("/service/selldetail", &service.GetOrderDetail{})     //售卖详情
	beego.Router("/service/working", &service.WorkingStart{})          //开始接单
	beego.Router("/service/int", &service.IntData{})                   //初始化
	beego.Router("/service/video/gettoken", &service.GetAccessToken{}) //摄像头列表
	beego.Router("/service/storeinfo", &service.StoreInfo{})           //店铺信息
	beego.Router("/service/closed", &service.StoreStatus{})            //店铺信息
	beego.Router("/service/report", &service.ReportInfo{})             //用户信息
	/* --------  end   ---------客服接口*/

	/*客户端接口 ----------  start----------*/
	//beego.Router("/alipay/token", &users.AlipayGetToken{})//获取openid
	beego.Router("/weixin/code", &users.CodeSession{})         //获取openid
	beego.Router("/weixin/notify", &users.WxPayNotify{})       //回调
	beego.Router("/alipay/notify", &users.AlipayNotify{})      //回调
	beego.Router("/user/enter", &customer.EnterController{})   //用户进入
	beego.Router("/user/leave", &customer.LeaveController{})   //
	beego.Router("/goods", &customer.GoodsController{})        //
	beego.Router("/order/sub", &customer.SubOrderController{}) //
	beego.Router("/user/update", &users.UpdateUser{})          //
	beego.Router("/storeinfo", &customer.StoreInfo{})          //
	beego.Router("/orderlist", &customer.OrderList{})          //
	beego.Router("/orderdetail", &customer.OrderDetail{})      //
	beego.Router("/customerphone", &customer.GetPhone{})

	beego.Router("/orderpaid", &customer.OrderPaid{})
	/* --------  end   ---------客户端接口*/

	/*店家接口 ----------  start----------*/
	beego.Router("/store/login", &store.Login{})                    //获取openid
	beego.Router("/store/goods", &store.GoodsController{})          //获取openid
	beego.Router("/store/goods/add", &store.ProductAddController{}) //添加商品
	beego.Router("/store/enterlog", &store.EnterLogs{})             //进入记录
	beego.Router("/store/enterdetail", &store.EnterDetail{})        //进入记录详情
	beego.Router("/store/goods/list", &store.ProductList{})         //获取商品列表
	beego.Router("/store/goods/edit", &store.ProductEdit{})         //编辑商品
	beego.Router("/store/order/list", &store.OrderList{})           //获取订单列表
	beego.Router("/store/withdraw", &store.WithDraw{})              //提现
	beego.Router("/store/storelist", &store.StoreList{})            //获取店铺列表
	beego.Router("/store/storeinfo", &store.StoreInfo{})            //获取 修改
	beego.Router("/store/storestatus", &store.StoreStatus{})        //更改店铺营业状态
	beego.Router("/store/videos", &store.GetVideos{})               //摄像头列表
	beego.Router("/store/wxconfig", &store.GetJsConfig{})           //微信js
	beego.Router("/store/bill", &store.BillLists{})                 //商家账单
	beego.Router("/store/finance", &store.Finance{})                //财务信息
	beego.Router("/store/userinfo", &store.UserInfo{})              //获取用户信息
	beego.Router("/store/catetree", &store.GateTree{})              //获取商品分类
	beego.Router("/store/suborder", &store.SubOrderController{})    //提交订单
	beego.Router("/store/orderpaid", &store.OrderPaid{})            //更改订单为支付
	/* --------  end   ---------店家接口*/
}
