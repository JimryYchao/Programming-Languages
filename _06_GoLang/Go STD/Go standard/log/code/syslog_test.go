package gostd

// TODO : syslog 在 windows 上未实现

/*
! Dial 通过连接到指定网路的地址`raddr` 来建立与日志守护程序的连接，对其返回 `Writer` 的每次写入都会发送一条日志消息（包含 facility、severity 和 tag）
		`tag` 为空时使用 os.Args[0]
		`network` 为空时，Dial 将连接到本地系统日志服务器。
		`network`, `raddr` 参阅 [net.Dial]
*/
