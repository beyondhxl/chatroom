package protocol

import (
	"fmt"
	"io"
)

type CmdWriter struct {
	writer io.Writer
}

// 构造指令写者
func NewCmdWriter(writer io.Writer) *CmdWriter {
	return &CmdWriter{
		writer: writer,
	}
}

// 写字符串
func (w *CmdWriter) writeStr(msg string) error {
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CmdWriter) Write(cmd interface{}) error {
	var err error
	switch msg := cmd.(type) {
	case CS_CmdSend: // 客户端发消息，就是往聊天框发送消息
		err = w.writeStr(fmt.Sprintf("SEND %v\n", msg.Msg))
	case SCS_CmdMessage: // 服务器广播消息
		err = w.writeStr(fmt.Sprintf("MESSAGE %v %v \n", msg.Name, msg.Msg))
	case CS_CmdName: // 客户端设置名称，以区分不同的聊天者
		err = w.writeStr(fmt.Sprintf("NAME %v\n", msg.Name))
	default:
		fmt.Print("未知的消息类型")
	}
	return err
}
