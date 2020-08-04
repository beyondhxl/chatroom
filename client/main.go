package main

import (
	"chatroom/client/cli"
	"fmt"
)

func main() {
	client := cli.NewTcpChatClient()
	err := client.Dial(":8080")
	if err != nil {
		fmt.Print(err)
		return
	}
	defer client.Close()

	go client.Start()
}
