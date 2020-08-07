/*
一、聊天室的基础的功能：
1、一个简单的聊天室
2、用户可以连接到这个聊天室
3、用户可以设置他们连接时的用户名
4、用户可以在里面发消息，并且消息会被广播给所有其他用户

通讯协议
客户端和服务器通信，要有一定的规则，就是相互之间要能听懂对方的话语。
一般是包头（PacketHead） + 包体（PacketLen）封装一个指令

模仿 Redis 设计消息指令，以 \n 做结束符
发送指令（SEND）：客户端可以发送聊天消息
命名指令（NAME）：客户端设置用户名
消息指令（MESSAGE）：服务端广播聊天消息给其他用户

例如，要发送一个 “Hello” 的消息，客户端会将字符串 SEND Hello\n 提交给 TCP socket，
服务端接受后会广播 MESSAGE username Hello\n 给其他用户。
*/

package protocol

// C -> S 客户端到服务器发送消息
type CS_CmdSend struct {
	Msg string // 消息
}

// C -> S
type CS_CmdName struct {
	Name string // 客户端要设置的名字
}

// S -> C || C -> S （服务器收到客户端的消息，然后广播给其他客户端）
type SCS_CmdMessage struct {
	Name string // 名字
	Msg  string // 具体消息
}
