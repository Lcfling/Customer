package users

import (
	"encoding/json"
	"fmt"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/utils"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type Users struct {
	Id          int64 `orm:"pk;column(uid);"`
	Account     string
	Balance     int64
	Pwd         string
	Token       string
	Openid      string
	Accesstoken string
	Nickname    string
	Avatar      string
	Phone       string
	AddTime     int64
	AddIp       string
	LastTime    int64
	LastIp      string
	Status      int
	Isleave     int
	Online      int
	UserType    int
	PayCount    int64
	LoginType   string
}

func (this *Users) TableName() string {
	return models.TableName("user")
}
func init() {
	orm.RegisterModel(new(Users))
}
func ChangeUserOnline(uid int64, status int) {
	var mark string
	user, _ := GetUser(uid)
	var types int64
	if status == 0 {
		mark = "TCP断开"
		types = 3
	} else {
		mark = "TCP接通"
		types = 4
	}
	if user.UserType == 4 {
		go logs.AddUserlog(uid, mark, types, user.UserType)
	}
	UpdateOnline(uid, status)
}

func LoginByOpenid(openid string) (Users, bool) {

	o := orm.NewOrm()
	user := Users{Openid: openid}
	err := o.Read(&user, "openid")
	if err != nil && err.Error() == "<QuerySeter> no row found" {
		return user, false
	}
	return user, true
}

func GreateUser(user *Users) (uid int64, token string, err error) {
	o := orm.NewOrm()
	user.AddTime = time.Now().Unix()
	user.LastTime = time.Now().Unix()
	tokeni := utils.SnowFlakeId()
	token = strconv.FormatInt(tokeni, 10)
	user.Token = utils.Md5(token)
	id, err := o.Insert(user)
	return id, user.Token, err
}

func UpdateUser(uid int64, user Users) error {
	o := orm.NewOrm()

	userinfo := Users{Id: uid}

	err := o.Read(&userinfo, "id")
	if nil != err {
		return err
	}
	if user.Phone != "" {
		userinfo.Phone = user.Phone
	}
	if user.Nickname != "" {
		userinfo.Nickname = user.Nickname
	}
	if user.Avatar != "" {
		userinfo.Avatar = user.Avatar
	}

	_, err = o.Update(&userinfo, "avatar", "nickname", "phone")
	if err != nil {
		return err
	}
	err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	return err
}
func UpdateLeave(uid int64, status int) error {
	o := orm.NewOrm()

	userinfo := Users{Id: uid}

	err := o.Read(&userinfo, "id")
	if nil != err {
		return err
	}
	userinfo.Isleave = status
	_, err = o.Update(&userinfo, "isleave")
	if err != nil {
		return err
	}
	err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	return err
}
func UpdateAccesstoken(uid int64, token string) error {
	o := orm.NewOrm()

	userinfo := Users{Id: uid}

	err := o.Read(&userinfo, "id")
	if nil != err {
		return err
	}
	userinfo.Accesstoken = token
	_, err = o.Update(&userinfo, "accesstoken")
	if err != nil {
		return err
	}
	err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	return err
}
func UpdateOnline(uid int64, status int) error {

	utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	o := orm.NewOrm()

	userinfo := Users{Id: uid}

	err := o.Read(&userinfo, "id")
	if nil != err {
		return err
	}
	userinfo.Online = status
	_, err = o.Update(&userinfo, "online")
	if err != nil {
		return err
	}
	go func() {
		//延时双删 双写一致性
		time.Sleep(time.Duration(1) * time.Second)
		err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	}()
	fmt.Println("UpdateOnline err", err)
	return err
}
func GetUser(uid int64) (Users, error) {
	info, err := utils.Get(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	fmt.Println("redis:", info)
	if err == nil && info == "" {
		o := orm.NewOrm()
		user := Users{Id: uid}
		err := o.Read(&user, "id")

		userJson, err := json.Marshal(user)
		if err != nil {
			return user, err
		}
		err = utils.Set(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10), string(userJson))
		return user, err
	} else if info != "" {
		var user Users
		err := json.Unmarshal([]byte(info.(string)), &user)
		return user, err
	} else {
		return Users{}, err
	}
}
func GetUserByAccount(account string) (Users, error) {
	o := orm.NewOrm()
	user := Users{Account: account}
	err := o.Read(&user, "account")
	return user, err
}
func UpdateToken(uid int64, token string) error {
	o := orm.NewOrm()
	utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	userinfo := Users{Id: uid}
	userinfo.Token = token
	_, err := o.Update(&userinfo, "token")
	if err != nil {
		return err
	}
	go func() {
		//延时双删 双写一致性
		time.Sleep(time.Duration(1) * time.Second)
		err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	}()

	return nil
}
func UpdatePwd(uid int64, newpwd string) error {
	o := orm.NewOrm()
	utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	userinfo := Users{Id: uid}
	userinfo.Pwd = newpwd
	_, err := o.Update(&userinfo, "pwd")
	if err != nil {
		return err
	}
	go func() {
		//延时双删 双写一致性
		time.Sleep(time.Duration(1) * time.Second)
		err = utils.DEL(models.GetRedis(), "userinfo_"+strconv.FormatInt(uid, 10))
	}()

	return nil
}
