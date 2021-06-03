/**
 * Copyright (c) 2014-2015, GoBelieve
 * All rights reserved.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

package socket

//换靴
const SVR_NOTIFYROOM_BOOT_NUM = 2

//洗牌中
const SVR_NOTIFYROOM_SHUFFLE = 3

//新一轮开始更改铺次 开始下注
const SVR_NOTIFYROOM_PAVE_NUM = 4

//房间状态更改
const SVR_NOTIFYROOM_SETTLEMENTOVER = 5

//游戏记录结果推送
const SVR_NOTIFYUSER_Winner = 6

//通知用户结算
const SVR_NOTIFYUSER_SETTLEMENTInfo = 7

//牛牛定位
const SVR_NOTIFYROOM_ORIENTATION = 8

//牛牛发牌
const SVR_NOTIFYROOM_GAMECARD = 9

//推送桌上筹码
const SVR_NOTIFYROOM_BET = 10

//推送牛牛彩池
const SVR_NOTIFYROOM_JACKPOT = 11

//推送龙虎彩池
const SVR_NOTIFYROOM_LHPOOL = 12

//推送百家乐彩池
const SVR_NOTIFYROOM_BJLPOOL = 13

//推送会员打赏主播
const SVR_NOTIFYUSER_REWARD = 14

//提前停止倒计时 (开牌中)
const SVR_NOTIFYROOM_SHOWCARD = 15

//下注前三推送
const SVR_NOTIFYROOM_TOP3 = 16

//用户余额实时推送
const SVR_USERBALANCE = 20

//房间清楚聊天信息
const SVR_NOTIFYROOM_CLEAR_TALK = 50

//推送用户实时余额
type UserBalance struct {
	Cmd         int
	UserBalance float64
}

//龙虎下注总筹码
type LHSumBets struct {
	Cmd         int //操作指令
	DragonMoney int
	TigerMoney  int
	TieMoney    int
}

//百家乐下注总筹码
type BJLSumBets struct {
	Cmd             int //操作指令
	PlayerMoney     int
	PlayerPairMoney int
	TieMoney        int
	BankerMoney     int
	BankerPairMoney int
}

// 牛牛下注总筹码
type NnSumBets struct {
	Cmd            int //操作指令
	DeskId         int //台桌id
	IdleOneMoney   int
	IdleTwoMoney   int
	IdleThreeMoney int
}

// 三公总筹码g
type SgSumBets struct {
	Cmd            int //操作指令
	DeskId         int //台桌id
	IdleOneMoney   int
	IdleTwoMoney   int
	IdleThreeMoney int
	//IdleFourMoney  int //闲4总筹码
	//IdleFiveMoney  int //闲5总筹码
	//IdleSixMoney   int //闲6总筹码
}

// 桌上筹码推送 推送
type Bet struct {
	Cmd      int //操作指令
	DeskId   int //台桌id
	Boot_num int //靴次
	Pave_num int //铺次
	BetMoney int //桌上筹码
}

// 牛牛定位 推送
type Location struct {
	Cmd         int //操作指令
	DeskId      int //台桌id
	Boot_num    int //靴次
	Pave_num    int //铺次
	LocationNum int //定位
}

// 龙虎&&百家乐&&牛牛 房间状态推送
type Phase struct {
	Cmd             int   //操作指令
	DeskId          int   //台桌id
	Boot_num        int   //靴次
	Pave_num        int   //铺次
	Phase           int   //房间状态
	CountDown       int   //倒计时
	GameStarTime    int64 //游戏开始时间
	Systime         int64 //服务器当前时间
	BetMoney        int   //桌上筹码
	IdleOneMoney    int64 //闲1总筹码
	IdleTwoMoney    int64 //闲2总筹码
	IdleThreeMoney  int64 //闲3总筹码
	IdleFourMoney   int64 //闲4总筹码
	IdleFiveMoney   int64 //闲5总筹码
	IdleSixMoney    int64 //闲6总筹码
	DragonMoney     int64 //龙总筹码
	TigerMoney      int64 //虎总筹码
	LHTieMoney      int64 //龙虎的和总筹码
	PlayerMoney     int64 //闲总筹码
	PlayerPairMoney int64 //闲对总筹码
	BJLTieMoney     int64 //百家乐的和总筹码
	BankerMoney     int64 //庄总筹码
	BankerPairMoney int64 //庄对总筹码
}

//龙虎 游戏记录结果推送
type LhWinner struct {
	Cmd      int
	DeskId   int
	Id       int
	Status   int
	Winner   int
	Boot_num int //靴次
	Pave_num int //铺次
}

//百家乐 游戏记录结果推送
type BjlWinner struct {
	Cmd      int
	DeskId   int
	Id       int
	Status   int
	Winner   string
	Boot_num int //靴次
	Pave_num int //铺次
}

//结算结果推送
type SettlementInfo struct {
	Cmd         int
	GameType    int         //1百家2龙虎3牛牛4三公5A89
	DeskId      interface{} //台桌id
	Desk_name   interface{} //台桌名称
	Boot_num    interface{} //靴次
	Pave_num    interface{} //铺次
	Result      interface{} //结果
	Sumgetmoney interface{} //输赢
	Balance     interface{} //剩余金额
	Todaymoney  interface{} //每日输赢
}

////牛牛游戏房间状态推送
//type NnPhase struct {
//	Cmd          int   //操作指令
//	DeskId       int   //台桌id
//	Boot_num     int   //靴次
//	Pave_num     int   //铺次
//	Phase        int   //房间状态
//	CountDown    int   //倒计时
//	GameStarTime int64 //游戏开始时间
//
//}

//牛牛游戏记录结果推送
type NnWinner struct {
	Cmd      int
	DeskId   int //台桌id
	Status   int
	Winner   string
	Boot_num int //靴次
	Pave_num int //铺次
}

//结算结果推送
type NnSettlementInfo struct {
	Cmd         int
	Boot_num    int         //靴次
	Pave_num    int         //铺次
	DeskId      int         //台桌id
	Desk_name   string      //台桌名称
	Result      string      //结果
	Sumgetmoney interface{} //输赢
	Balance     interface{} //剩余金额
	Todaymoney  interface{}
	GameType    int
}

//会员打赏主播推送
type U2LReward struct {
	Cmd         int
	UserAccount string //会员账号
	Money       int    //打赏金额
}

//台桌下注前三
type DeskTop3 struct {
	Cmd  int
	Top3 map[string]interface{}
}

//清楚房间聊天信息
type ClearTalk struct {
	Cmd   int
	Clear int
}
