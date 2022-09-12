package main

import (
	"fmt"
	"net"
)

func main()  {
	//创建服务器
	listener,err := net.Listen("tcp",":8080")
	if err != nil{
		fmt.Println("net .Listen err:",err)
		return
	}

	//监听
	conn,err := listener.Accept()
	if err != nil {
		fmt.Println("net .Listen err:",err)
		return
	}
	fmt.Println("建立连接成功")

}
