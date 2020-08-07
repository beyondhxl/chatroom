package srv

import (
	"chatroom/common"
	"chatroom/protocol"
	"fmt"
	"net"
	"sync"
)

// TcpChatServer 上面管理了实际客户端发过来的连接
type TTCPChatServer struct {
	tMutex     sync.Mutex   // 互斥锁，保护 clients
	listener   net.Listener // 侦听器
	slcClients []*TClient   // 连接服务器的客户端（聊天用户）
}

// TCPChatServer 上的客户端，并非实际的客户端，方便管理
type TClient struct {
	conn    net.Conn             // 原生连接
	strName string               // 聊天者名字
	tWriter *protocol.TCmdWriter // 写指令器
}

func NewTcpChatServer() *TTCPChatServer {
	return &TTCPChatServer{
		tMutex: sync.Mutex{},
	}
}

// 实现 ChatServer Listen() 接口
func (this *TTCPChatServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err == nil {
		this.listener = l
	}
	fmt.Printf("listening on %v\n", addr)
	return err
}

// 实现 ChatServer Close() 接口
func (this *TTCPChatServer) Close() {
	this.tMutex.Lock()
	defer this.tMutex.Unlock()

	// 移除所有客户端
	for _, cli := range this.slcClients {
		this.remove(cli)
	}

	this.listener.Close() // 侦听关闭
}

// ChatServer Start() 接口
func (this *TTCPChatServer) Start() {
	// 同时启动协程处理服务端消息及命令
	go this.sendLoop()

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
		go this.recvLoop(client)
	}

}

// 内部方法 接受新连接
func (this *TTCPChatServer) accept(conn net.Conn) *TClient {
	this.tMutex.Lock()
	defer this.tMutex.Unlock()
	client := &TClient{
		conn:    conn,
		strName: "",
		tWriter: protocol.NewCmdWriter(conn),
	}
	this.slcClients = append(this.slcClients, client)
	return client
}

// 内部方法 移除连接
func (this *TTCPChatServer) remove(client *TClient) {
	this.tMutex.Lock()
	defer this.tMutex.Unlock()

	for i, remove := range this.slcClients { // 比较巧妙，因为移除时，用到了下标
		if remove == client {
			this.slcClients = append(this.slcClients[:i], this.slcClients[i+1:]...)
		}
	}
	client.conn.Close() // 连接关闭
}

// 内部方法 读消息循环
func (this *TTCPChatServer) recvLoop(client *TClient) {
	for {

		data := make([]byte, 1024)
		//读取消息
		client.conn.Read(data)

		fmt.Println(string(data))

		// 广播
		this.Broadcast(string(data))
	}
}

func (this *TTCPChatServer) sendLoop() {
	for {
		cmdLine := common.ScanLine() // 读取命令行输入
		fmt.Println(cmdLine)

		// 直接广播
		this.Broadcast(cmdLine)
	}
}

// 广播 因为 TcpChatServer 上有所有的聊天者（客户端）
func (this *TTCPChatServer) Broadcast(message string) {
	for _, client := range this.slcClients {
		/*
			client.conn.Write([]byte(message))
		*/
		client.tWriter.Write(message)
	}
}
