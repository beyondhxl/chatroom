package cli

import (
	"chatroom/common"
	"chatroom/protocol"
	"fmt"
	"net"
)

// 实际的 TCPChatClient
type TTCPChatClient struct {
	conn    net.Conn             // 实际的网络连接
	strName string               // 客户端名称
	tReader *protocol.TCmdReader // 读指令器
	tWriter *protocol.TCmdWriter // 写指令器
}

// 构造聊天客户端
func NewTCPChatClient() *TTCPChatClient {
	return &TTCPChatClient{}
}

// 实现 IClient Dial() 接口
func (this *TTCPChatClient) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	this.conn = conn                           // 保存客户端连接
	this.tWriter = protocol.NewCmdWriter(conn) // 构造读写器
	this.tReader = protocol.NewCmdReader(conn)
	return err
}

// 实现 IClient Start() 接口
func (this *TTCPChatClient) Start() {
	// 启动消息循环 goroutine
	go this.recvLoop()

	// 发送消息循环
	for {
		message := common.ScanLine() // 从命令行获取用户输入
		// 消息格式
		// MESSAGE username msg\n
		// SEND msg\n
		// NAME username\n
		this.tWriter.Write(message)

		/*
			this.Send(message) // 发送消息
		*/
	}

}

// 实现 IClient Close() 接口
func (this *TTCPChatClient) Close() {
	this.conn.Close()
}

// 实现 IClient SetName() 接口
func (this *TTCPChatClient) SetName(name string) {
	this.strName = name
}

// 实现 IClient Send() 接口
func (this *TTCPChatClient) Send(message string) {
	this.conn.Write([]byte(message))
}

// 接收循环
func (this *TTCPChatClient) recvLoop() {
	// 也是需要一直循环侦听服务端发来的消息
	for {
		/*
			data := make([]byte, 1024)
			//读取消息
			this.conn.Read(data)
			fmt.Println(string(data))
		*/
		msg, _ := this.tReader.Read()
		fmt.Println(msg)
	}
}
