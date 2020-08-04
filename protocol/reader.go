package protocol

import (
	"bufio"
	"fmt"
	"io"
)

type CmdReader struct {
	reader *bufio.Reader
}

func NewCmdReader(r io.Reader) *CmdReader {
	return &CmdReader{
		reader: bufio.NewReader(r),
	}
}

// 读消息
func (r *CmdReader) Read() (interface{}, error) {
	cmdName, err := r.reader.ReadString(' ') // 这里不能用""，字符用''，字符串用""
	if err != nil {
		return nil, err
	}
	switch cmdName {
	case "MESSAGE":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return nil, err
		}
		msg, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		// 构造消息
		return SCS_CmdMessage{
			user[:len(user)-1],
			msg[:len(msg)-1],
		}, nil
	case "SEND":
		msg, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return CS_CmdSend{
			Msg: msg[:len(msg)-1],
		}, nil
	case "NAME":
		name, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return CS_CmdName{
			Name: name[:len(name)-1],
		}, nil
	default:
		fmt.Printf("未知的指令类型 %s", cmdName)
		return nil, err
	}
}
