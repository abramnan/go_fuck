package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// 验证点：
// 1. 在server启动listen后，sleep，是否可以建立链接，并且发送数据(服务端没有执行accept，是否可以建立tcp链接) -> 可以
// 2. 如果两边建立了链接，然后有一端sleep休眠了，请问write和send能不能成功 -> 可以

// 验证过程
// 启动server，执行到listen后，陷入睡眠；启动client，这个时候没有链接断开，而是正常链接，切可以发送数据，只是服务端没有收到
// 唤醒sever，会接收到刚刚客户端发送的数据

// 一些疑问：
// 什么是粘包
//

func main() {
	// 包含了bind和listen的过程
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen failed:", err)
		return
	}
	var listenAfterSleep bool
	flag.BoolVar(&listenAfterSleep, "lsa", false, "listen after sleep")
	flag.Parse()
	fmt.Printf("server start %d...\n", listenAfterSleep)
	// 验证服务端listen后睡眠，是否可以正常链接
	if listenAfterSleep {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			// 没有输入的时候，会一直阻塞在这里
			fmt.Println("server is sleeping.......")
			input, _ := inputReader.ReadString('\n') // 读取输入
			inputInfo := strings.Trim(input, "\r\n")
			// 输入up唤醒
			if strings.ToUpper(inputInfo) == "UP" {
				break
			}
			//time.Sleep(60 * time.Minute)
		}
	}
	for {
		fmt.Println("before accept")
		// 验证accept后（即成功建立链接后，另一方sleep了，write和send能够成功）
		// 客户端没启动的时候，会阻塞在accept处
		// 这里使用循环的原因是，如果不使用循环，你这里建立一个链接后，服务端就直接退出了
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("server accept failed:", err)
			return
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	fmt.Println("server ready")
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client falied:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("server recv data:", recvStr)
		conn.Write([]byte(recvStr + "server")) // 发送数据
	}
}
