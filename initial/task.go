package initial

import (
	"errors"
	"fmt"
	"github.com/Lcfling/Customer/models"
	"github.com/Lcfling/Customer/models/device"
	"github.com/Lcfling/Customer/models/logs"
	"github.com/Lcfling/Customer/models/order"
	"github.com/Lcfling/Customer/utils"
	"math/rand"
	"strconv"
	"time"
)

func Inttask(){

	go SettlementBlockQueue()
	go func() {
		for {


			exptime:=time.Now().UnixNano()/1e6+82800000

			for{
				list,err:=device.GetVideobyTime(exptime)
				if err!=nil{
					fmt.Println("err:",err)
					time.Sleep(time.Duration(20)*time.Second)
					continue
				}
				fmt.Println("list:",list)

				if len(list)==0{
					break
				}
				for _,v:=range list{
					accesstoken,expt,err:=utils.YSGetAccesstoken(v.Appid,v.Secret)
					if err!=nil{
						continue
					}
					err=device.UpdateVideoToken(accesstoken,strconv.FormatInt(expt,10),v.Appid)
					if err!=nil {
						fmt.Println("更新Accesstoken失败：",err,v)
					}
				}
			}





			Time:=rand.Intn(10)+14400
			time.Sleep(time.Duration(Time)*time.Second)
		}
	}()
	for{
		Time:=rand.Intn(10)+5
		time.Sleep(time.Duration(Time)*time.Second)
		DateString:=utils.GetDate(time.Now().Unix()-86400)//只结算昨天的

		lockKey:=utils.Md5(strconv.FormatInt(utils.SnowFlakeId(),10))
		err:=utils.SetNx(models.GetRedis(),"settlement_lock",lockKey)
		if err!=nil{
			fmt.Println("分布式锁错误",err)
			continue
		}
		lockvalue,_:=utils.Get(models.GetRedis(),"settlement_lock")
		if lockvalue!=lockKey {
			continue //别人的锁权限
		}
		NowDateString,err:=utils.Get(models.GetRedis(),"NowDateString")
		if err!=nil{
			fmt.Println("定时任务获取redis键值错误",err)
			err=utils.DEL(models.GetRedis(),"settlement_lock")
			if err!=nil{
				fmt.Println("定时任务解锁失败")
			}
			continue
		}
		if NowDateString.(string)==DateString{
			err=utils.DEL(models.GetRedis(),"settlement_lock")
			if err!=nil{
				fmt.Println("定时任务解锁失败")
			}
			goto SlEEP
		}

		//开始结算
		err=Settlement(DateString)
		if err!=nil{
			err=utils.DEL(models.GetRedis(),"settlement_lock")
			if err!=nil{
				fmt.Println("定时任务解锁失败")
			}
			continue
		}
		err=utils.Set(models.GetRedis(),"NowDateString",DateString)
		if err!=nil{
			fmt.Println("更新键值（NowDateString）失败")
		}
		//DateString
		//主动解锁
		err=utils.DEL(models.GetRedis(),"settlement_lock")
		if err!=nil{
			fmt.Println("定时任务解锁失败")
		}
		SlEEP:
		Time=rand.Intn(5)+14400
		time.Sleep(time.Duration(Time)*time.Second)
	}

}

func SettlementBlockQueue(){
	var times int32 = 0
	for{
		times++
		fmt.Println("SettlementBlockQueue start :",times," time",time.Now().Unix())
		arr,err:=utils.BLPop(models.GetRedis(),"settlement_by_day")
		if err!=nil{
			fmt.Println("监听到阻塞队列错误：",arr,err.Error())
			continue
		}
		fmt.Println("监听到阻塞队列错误：",arr)

		//for

		list:=arr.([]string)

		err=Settlement(list[1])
		if err!=nil{
			fmt.Println("对接结算出错：",err.Error())
			continue
		}
		fmt.Println("SettlementBlockQueue end :",times)
	}
}
func Settlement(DateString string)error{

	yesterdaybegin:=utils.GetDateParse(DateString)
	todaybegin:=utils.GetDateParse(DateString)+86400

	//获取销售的商品总数
	sellcounts,err:=order.GetSellCounts(yesterdaybegin,todaybegin)
	if err!=nil {
		fmt.Println("统计销售商品总数失败")
		return errors.New("统计销售商品总数失败")
	}

	//获取当日订单总金额

	handle,err:=order.GetSellMoneySum(yesterdaybegin,todaybegin,-2,"add_time")
	//订单总金额
	SumOrderTotalMoney:=handle.TotalPrice
	handle,err=order.GetSellMoneySum(yesterdaybegin,todaybegin,-2,"pay_time")
	//总成交金额 退款时 请修改此处的金额
	SumTradePayMoney:=handle.PayPrice


	//总订单量
	ordercounts,err:=order.OrderCounts(yesterdaybegin,todaybegin,-2)

	payordercounts,err:=order.OrderCounts(yesterdaybegin,todaybegin,1)
	payordercounts2,err:=order.OrderCounts(yesterdaybegin,todaybegin,-1)
	//成功和退款都算  成交订单量
	payordercounts=payordercounts+payordercounts2

	//总访问量

	counts,err:=logs.GetCounts(yesterdaybegin,todaybegin)

	settlement:=new(order.Settlement)
	settlement.Day=DateString
	settlement.Ordercounts=ordercounts
	settlement.Sellcounts=sellcounts
	settlement.Ordersuccess=payordercounts
	settlement.Ordermoney=SumOrderTotalMoney
	settlement.Trademoney=SumTradePayMoney
	settlement.Visitcounts=counts
	err=order.SaveSettlement(settlement)
	if err!=nil{
		//fmt.Println("统计销售商品总数失败")
		return err
	}
	err=utils.Set(models.GetRedis(),"NowDateString",DateString)
	if err!=nil{
		fmt.Println("更新键值（NowDateString）失败")
	}
	return nil
}