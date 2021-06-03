package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//百分比例函数
func Ratio(a int64,s int64,j int) float64{
	r:=float64(a)/float64(s)
	return FloatRound(r,j)
}
// 截取小数位数
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}
const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte(char[rand.Int63() % int64(len(char))])
	}
	return s.String()
}
