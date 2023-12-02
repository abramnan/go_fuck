package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("client connect failed:", err)
		return
	}
	defer conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取输入
		inputInfo := strings.Trim(input, "\r\n")
		// 输入Q退出，关闭链接
		if strings.ToUpper(inputInfo) == "Q" {
			return
		}
		_, err := conn.Write([]byte(inputInfo)) // 发送数据
		if err != nil {
			fmt.Println("client sent data failed:", err)
			return
		}
		buf := [512]byte{}
		n, err := conn.Read(buf[:]) // 接收数据
		if err != nil {
			fmt.Println("client recv data failed:", err)
			return
		}
		fmt.Println("client recv data:" + string(buf[:n]))
	}
}
