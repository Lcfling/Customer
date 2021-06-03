package utils

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func LPush(pool *redis.Pool,key string,value interface{}) error {

	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("LPUSH",key,value)
	if err != nil {
		fmt.Println("redis LPUSH failed:", err)
	}
	return err
}
func RPush(pool *redis.Pool,key string,value interface{}) error {

	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("RPUSH",key,value)
	if err != nil {
		fmt.Println("redis LPUSH failed:", err)
	}
	return err
}
func LPop(pool *redis.Pool,key string) (interface{},error ){

	conn:= pool.Get()
	defer conn.Close()
	reply,err:=redis.String(conn.Do("LPOP",key))
	if err != nil {
		return nil,err
	}
	return reply,nil
}
func BLPop(pool *redis.Pool,key string) (interface{},error ){

	conn:= pool.Get()
	defer conn.Close()
	reply,err:=redis.Strings(conn.Do("BLPOP",key,0))
	if err != nil {
		return nil,err
	}
	return reply,nil
}
func LRange(pool *redis.Pool,key string,start int64,end int64)([]interface{},error){
	conn:= pool.Get()
	defer conn.Close()
	reply,err:=redis.Values(conn.Do("LRANGE",key,start,end))
	if err != nil {
		return nil,err
	}
	return reply,nil
}

func Get(pool *redis.Pool,key string) (interface{},error) {
	conn:= pool.Get()
	defer conn.Close()

	exsit,err:=redis.Int64(conn.Do("EXISTS",key))
	if err==nil && exsit==0{
		return "",nil
	}
	reply,err:=redis.String(conn.Do("GET",key))
	if err != nil {

		fmt.Println("redis Get failed:", err)
		return "",err
	}
	return reply,nil
}
func Set(pool *redis.Pool,key string,value string) error{
	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("SET",key,value)
	if err != nil {
		fmt.Println("redis SET failed:", err)
	}
	return err
}

func LRem(pool *redis.Pool,key string,count int64,value interface{}) error{
	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("LREM",key,count,value)
	if err != nil {
		fmt.Println("redis LREM failed:", err)
	}
	return err
}
func DEL(pool *redis.Pool,key string)error{
	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("DEL",key)
	if err != nil {
		fmt.Println("redis DEL failed:", err)
	}
	return err
}
//分布式锁
func SetNx(pool *redis.Pool,key string,value string) error{
	conn:= pool.Get()
	defer conn.Close()
	replay,err:=conn.Do("SET",key,value,"NX","PX","30000")



	//fmt.Println("redis SetNx failed:", replay)
	if err != nil {
		fmt.Println("redis SetNx failed:", err)
	}
	switch replay.(type) {
	case string:
		if replay.(string)=="ok"{
			return nil
		}
	case nil:
		return errors.New("key has exsit")
	}
	return err
}

func SetEx(pool *redis.Pool,key string,value string,s string) error{
	conn:= pool.Get()
	defer conn.Close()
	_,err:=conn.Do("SET",key,value,"EX",s)

	//fmt.Println("redis SetNx failed:", replay)
	if err != nil {
		fmt.Println("redis SetEX failed:", err)
	}
	if err != nil {
		fmt.Println("redis SET failed:", err)
	}
	return err
}