package srv

type IChatServer interface {
	Listen(addr string) error        // 侦听
	Broadcast(cmd interface{}) error // 广播
	Start()                          // 启动
	Close()                          // 关闭
}
