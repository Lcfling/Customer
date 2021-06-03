package main

import (
	"fmt"
	_ "github.com/Lcfling/Customer/initial"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/users"
	_ "github.com/Lcfling/Customer/routers"
	"github.com/Lcfling/Customer/service"
	"github.com/Lcfling/Customer/socket"
	"github.com/Lcfling/Customer/udpsok"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)
func main() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
	initArgs()
	//beego.InsertFilter("*",beego.BeforeRouter,optionFlater)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//允许访问所有源
		AllowAllOrigins: true,
		//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		//其中Options跨域复杂请求预检
		AllowMethods: []string{"*"},
		//指的是允许的Header的种类
		AllowHeaders: []string{"uid","token"},
		//公开的HTTP标头列表
		ExposeHeaders: []string{"*"},
		//如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))

	beego.ErrorHandler("404", page_not_found)
	beego.ErrorHandler("401", page_note_permission)
	go socket.Run(users.ChangeUserOnline)
	go udpsok.Run()
	go service.Run()
	beego.Run()


}

var optionFlater = func(ctx *context.Context) {
	fmt.Println("优先运行")
	ctx.Output.JSON(map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}, false, false)
	//ctx.Redirect(302, "/login")
}

func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
	}
}

var FilterUser = func(ctx *context.Context) {
	/*_, ok := ctx.Input.Session("userLogin").(string)
	if !ok && !(ctx.Request.RequestURI == "/login" ||strings.Contains(ctx.Request.RequestURI, "/register") ){
		ctx.Redirect(302, "/login")
	}*/
}

func page_not_found(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("404.tpl").ParseFiles("views/404.tpl")
	data := make(map[string]interface{})
	//data["content"] = "page not found"
	t.Execute(rw, data)
}

func page_note_permission(rw http.ResponseWriter, r *http.Request) {
	t, _ := template.New("401.tpl").ParseFiles("views/401.tpl")
	data := make(map[string]interface{})
	//data["content"] = "你没有权限访问此页面，请联系超级管理员。或去<a href='/'>开启我的OPMS之旅</a>。"
	t.Execute(rw, data)
}
