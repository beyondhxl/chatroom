package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// 读指令者
type TCmdReader struct {
	reader *bufio.Reader
}

// 构造函数
func NewCmdReader(r io.Reader) *TCmdReader {
	return &TCmdReader{
		reader: bufio.NewReader(r),
	}
}

// 读消息
func (r *TCmdReader) Read() (string, error) {
	strCmd, err := r.reader.ReadString(' ')
	cmdName := strings.TrimSpace(strCmd)
	if err != nil {
		return "", err
	}
	if strCmd == "" {
		return "", nil
	}
	switch cmdName {
	case "MESSAGE":
		user, err := r.reader.ReadString(' ')
		if err != nil {
			return "", err
		}
		msg, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return cmdName + user + msg, nil
	case "SEND":
		msg, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return cmdName + " " + msg, nil
	case "NAME":
		name, err := r.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return cmdName + name, nil
	default:
		fmt.Printf("未知的指令类型 %s", cmdName)
		return "", err
	}
}
