package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"time"
)
var redis_pool *redis.Pool

func init() {
	redis_pool=NewRedisPool(beego.AppConfig.String("redis_host"),beego.AppConfig.String("redis_password"),0)
}
func GetRedis()*redis.Pool{
	return redis_pool
}
func TableName(name string) string {
	return beego.AppConfig.String("mysqlpre") + name
}
func migirateTable() {
	//o := orm.NewOrm()
	//	tSql := "alter table `case` modify `story` varchar(400) NOT NULL"
	//tSql := "alter table `material` modify `name` varchar(50) NOT NULL"
	//tSql := "alter table `user` modify `password` varchar(100)  NULL"
	//o.Raw(tSql).Exec()
}
func Syncdb() {
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	//initRole()
	//migirateTable()
}
//redis句柄函数
func NewRedisPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   500,
		IdleTimeout: 480 * time.Second,
		Dial: func() (redis.Conn, error) {
			timeout := time.Duration(2) * time.Second
			c, err := redis.DialTimeout("tcp", server, timeout, 0, 0)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if db > 0 && db < 16 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}
