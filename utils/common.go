package utils

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
)

func GetOrderSN() string{
	timeLayout := "20060102150405"
	t := time.Now()
	pre := t.Format(timeLayout)
	rand.Seed(time.Now().UnixNano())
	suff := GetRandomIntString(6)
	order_sn := pre + suff
	return order_sn
}

func StructoMap(obj interface{})map[string]interface{}{

	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

func StoMap(struc interface{}) (map[string]interface{},error ){

	jsonBytes, err := json.Marshal(struc)
	if err!=nil{
		return nil,err
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal(jsonBytes, &mapResult)
	if err != nil {
		return nil,err
	}
	return mapResult,nil
}