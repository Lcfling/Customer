package initial

func init() {
	InitSql()
	InitTemplate()
	InitCache()
	go Intconfig()


}
