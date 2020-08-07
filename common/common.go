package common

import (
	"bufio"
	"os"
	"strings"
)

// ScanLine 读取命令行整行
func ScanLine() string {
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	return strings.Replace(input, "\n", "", -1)
}
