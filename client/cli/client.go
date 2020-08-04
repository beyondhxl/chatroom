package cli

type IChatClient interface {
	Dial(addr string) error          // 连接服务器
	Send(cmd interface{}) error      // 发送指令
	SendMessag(message string) error // 发送消息
	SetName(name string) error       //	设置名字
	Start()                          // 客户端启动
}
