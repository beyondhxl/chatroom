package cli

import (
	"fmt"
	"io"
	"net"

	"chatroom/protocol"
)

// 实际的 TcpChatClient
type TcpChatClient struct {
	conn      net.Conn                     // 实际的网络连接
	cmdReader *protocol.CmdReader          // 指令读者
	cmdWriter *protocol.CmdWriter          // 指令写者
	name      string                       // 客户端名称
	msgChan   chan protocol.SCS_CmdMessage // 消息通道缓存
}

func NewTcpChatClient() *TcpChatClient {
	return &TcpChatClient{
		msgChan: make(chan protocol.SCS_CmdMessage),
	}
}

// 实现 IClient Dial() 接口
func (this *TcpChatClient) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	this.conn = conn // 保存客户端连接
	this.cmdReader = protocol.NewCmdReader(conn)
	this.cmdWriter = protocol.NewCmdWriter(conn)
	return err
}

// 实现 IClient Send() 接口
func (this *TcpChatClient) Send(cmd interface{}) error {
	return this.cmdWriter.Write(cmd)
}

// 实现 IClient Start() 接口
func (this *TcpChatClient) Start() {
	// 也是需要一直循环侦听客户端发来的消息
	for {
		msg, err := this.cmdReader.Read()
		if err == io.EOF {
			break // 读到文件尾了
		} else {
			fmt.Print("111")
			fmt.Print(err)
		}
		if msg != nil {
			switch cmd := msg.(type) {
			case protocol.SCS_CmdMessage:
				this.msgChan <- cmd // 指令缓存到消息队列里
			default:
				fmt.Printf("未知的消息类型 %s", cmd)
			}
		}
	}
}

// 实现 IClient Close() 接口
func (this *TcpChatClient) Close() {
	close(this.msgChan)
	this.conn.Close()
}

// 实现 IClient SendMessag() 接口 直接往服务器发指令
func (this *TcpChatClient) SendMessag(message string) error {
	return this.Send(protocol.SCS_CmdMessage{Msg: message})
}

// 实现 IClient SetName() 接口 直接往服务器发指令
func (this *TcpChatClient) SetName(name string) error {
	return this.Send(protocol.CS_CmdName{Name: name})
}
