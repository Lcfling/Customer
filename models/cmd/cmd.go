package cmd

type UserEnter struct {
	Cmd int64
	Uid 	int64
	Storeid int64
	Doorid 	int64
}
//离开时关门 cmd=2
type OpenDoorInfo struct {
	Cmd 	int64
	Uid 	int64
	Storeid int64
	Doorid 	int64
}
//cmd =3
type SendOrder struct {
	Cmd int64
	Uid int64
	Storeid int64
	Ordersn string

}
//用户进入 cmd =4
type CustomerEnter struct {
	Cmd int64
	Uid int64
	Storeid int64
	Doorid 	int64
}

//cmd =5
type SendOrderStatus struct {
	Cmd int64
	Uid int64
	Storeid int64
	Ordersn string
	Status  int
}
