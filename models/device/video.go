package device

import (
	"fmt"
	"github.com/Lcfling/Customer/models"
	"github.com/astaxie/beego/orm"
)

type Video struct {
	Id          int64
	StoreId     int64
	Describe    string
	Url         string
	Isprimary   int
	Appid       string
	Secret      string
	Accesstoken string
	Expiretime  int64
	Safekey     string
	Devicesn    string
}

func (this *Video) TableName() string {
	return models.TableName("store_video")
}
func init() {
	orm.RegisterModel(new(Video))
}

func GetVideo(store_id int64) ([]Video, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(models.TableName("store_video"))
	cond := orm.NewCondition()
	cond = cond.And("store_id", store_id)
	qs = qs.SetCond(cond)
	var videos []Video
	_, err := qs.All(&videos)
	return videos, err
}
func GetVideoById(id int64) (Video, error) {
	o := orm.NewOrm()
	v := Video{Id: id}
	err := o.Read(&v)
	return v, err
}

func UpdateVideoToken(token, exp, appid string) error {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Update("eb_store_video").Set("accesstoken=?,expiretime=?").Where("appid=?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, token, exp, appid).Exec()
	return err
}

func GetVideobyTime(t int64) ([]Video, error) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("appid,secret").From("eb_store_video").
		Where("expiretime<?").Or("accesstoken=?").GroupBy("appid,secret")
	sql := qb.String()
	o := orm.NewOrm()

	var Videos []Video
	var empty = ""
	_, err := o.Raw(sql, t, empty).QueryRows(&Videos)
	fmt.Println(Videos)
	return Videos, err
}
