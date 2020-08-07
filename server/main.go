package main

import "chatroom/server/srv"

// 启动服务
func main() {
	var chatserver *srv.TTCPChatServer
	chatserver = srv.NewTcpChatServer()
	chatserver.Listen(":8080")
	chatserver.Start()
}
