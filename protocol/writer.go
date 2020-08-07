package protocol

import (
	"fmt"
	"io"
	"strings"
)

// 指令写者（负责套接字发送）
type TCmdWriter struct {
	writer io.Writer
}

// 构造指令写者
func NewCmdWriter(writer io.Writer) *TCmdWriter {
	return &TCmdWriter{
		writer: writer,
	}
}

// 内部方法
// 写字节（实际通过套接口发送字节流）
func (w *TCmdWriter) writeStr(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

// 根据指令类型打包消息并发送
func (w *TCmdWriter) Write(cmd string) error {
	slcCmd := strings.Split(cmd, " ")
	var err error
	switch slcCmd[0] {
	case "SEND": // 客户端发消息，就是往聊天框发送消息
		err = w.writeStr(fmt.Sprintf("SEND %v\n", slcCmd[1]))
	case "MESSAGE": // 服务器广播消息
		err = w.writeStr(fmt.Sprintf("MESSAGE %v %v \n", slcCmd[1], slcCmd[2]))
	case "NAME": // 客户端设置名称，以区分不同的聊天者
		err = w.writeStr(fmt.Sprintf("NAME %v\n", slcCmd[1]))
	default:
		fmt.Println("未知的消息类型")
	}
	return err
}
