package main

import (
	"chatroom/client/cli"
	"fmt"
)

// 客户端启动入口
func main() {
	client := cli.NewTCPChatClient() // 新建
	err := client.Dial(":8080")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer client.Close()

	client.Start()
}
