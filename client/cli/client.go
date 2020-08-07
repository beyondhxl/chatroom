package cli

type IChatClient interface {
	Dial(addr string) error // 连接服务器
	Start()                 // 客户端启动
	Close()                 // 客户端关闭
	SetName(name string)    // 设置名字
	Send(message string)    // 发送指令
}
