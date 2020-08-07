package srv

type IChatServer interface {
	Listen(addr string) error // 侦听
	Start()                   // 启动
	Close()                   // 关闭
	Broadcast(message string) // 广播
}
