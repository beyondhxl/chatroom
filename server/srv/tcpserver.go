package srv

import (
	"chatroom/protocol"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// TcpChatServer 上面管理了实际客户端发过来的连接
type TcpChatServer struct {
	mu       sync.Mutex   // 互斥锁，保护clients
	listener net.Listener // 侦听器
	clients  []*Client    // 连接服务器的客户端（聊天用户）
}

// 实在 TcpChatServer 上的客户端，并非实际的客户端，方便管理
type Client struct {
	conn   net.Conn            // 原生连接
	name   string              // 聊天者名字
	writer *protocol.CmdWriter // 用于发送消息
}

func NewTcpChatServer() *TcpChatServer {
	return &TcpChatServer{
		mu: sync.Mutex{},
	}
}

// 实现 ChatServer Listen() 接口
func (this *TcpChatServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err == nil {
		this.listener = l
	}
	fmt.Printf("listening on %v", addr)
	return err
}

// 实现 ChatServer Close() 接口
func (this *TcpChatServer) Close() {
	this.listener.Close() // 侦听关闭
}

// ChatServer Start() 接口
func (this *TcpChatServer) Start() {
	// 服务器启动，有个无限循环，一直接受客户端的连接
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			fmt.Printf("侦听出错 %s", err)
			return
		}
		// 将连接用客户端类包装起来
		client := this.accept(conn)
		// 单独起一个 goroutine 处理这个客户端连接
		go this.serve(client)
	}
}

// 内部方法 接受新连接
func (this *TcpChatServer) accept(conn net.Conn) *Client {
	this.mu.Lock()
	defer this.mu.Unlock()
	client := &Client{
		conn:   conn,
		writer: protocol.NewCmdWriter(conn),
	}
	this.clients = append(this.clients, client)
	log.Printf("hello,我")
	return client
}

// 内部方法 移除连接
func (this *TcpChatServer) remove(client *Client) {
	this.mu.Lock()
	defer this.mu.Unlock()
	for i, remove := range this.clients { // 比较巧妙，因为移除时，用到了下标
		if remove == client {
			this.clients = append(this.clients[:i], this.clients[i+1:]...)
		}
	}
	client.conn.Close() // 连接关闭
}

// 内部方法 服务连接
func (this *TcpChatServer) serve(client *Client) {
	cmdReader := protocol.NewCmdReader(client.conn)
	defer this.remove(client)
	for {
		msg, err := cmdReader.Read()
		if err != nil && err != io.EOF { // 不是读到了文件尾
			fmt.Printf("Read error : %s", err)
		}
		if msg != nil {
			switch cmd := msg.(type) {
			case protocol.CS_CmdName:
				client.name = cmd.Name // 设置名字
			case protocol.CS_CmdSend:
				// 单独起一个 goroutine 转发客户端的消息到其他客户端
				go this.Broadcast(protocol.SCS_CmdMessage{
					Msg:  cmd.Msg,
					Name: client.name,
				})
			}
		}
		if err == io.EOF { // 读到了消息的末尾（可能客户端 Ctrl + C）
			break
		}
	}
}

// 广播 因为 TcpChatServer 上有所有的聊天者（客户端）
func (this *TcpChatServer) Broadcast(cmd interface{}) {
	for _, client := range this.clients {
		client.writer.Write(cmd)
	}
}
