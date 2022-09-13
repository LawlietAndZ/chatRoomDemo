package main

import (
	"fmt"
	"net"
)

type User struct {
	//名字
	name string
	//唯一id
	id string
	//管道
	msg chan string
}

//建立全局map结构，用于保存所有用户
var allUsers = make(map[string]User)

//定义全局通道，用于接收所有用户发送过来的消息
var message = make(chan string,10)

func main()  {
	//创建服务器
	listener,err := net.Listen("tcp",":8080")
	if err != nil{
		fmt.Println("net .Listen err:",err)
		return
	}

	for {

		//监听
		fmt.Println("服务器启动，主go程监听中")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net .Listen err:", err)
			return
		}
		fmt.Println("建立连接成功")
		//启动处理业务的go程
		go handler(conn)
	}
}


//具体处理业务
func handler(conn net.Conn)  {
	fmt.Println("启动业务")

	//客户端与服务器建立链接时，会有ip和port，可以当初user的id

	clientAddr :=conn.RemoteAddr().String()





	for  {
		fmt.Println("启动业务")

		newUser := User{
			name: clientAddr,//后续会修改，初始值与id相同
			id: clientAddr,
			msg: make(chan string),
		}
		//添加User到map
		allUsers[newUser.id] = newUser

		//向message写入数据,当用户上线的时，将消息通知所有人
		loginInfo := fmt.Sprintf("[%s]:[%s]=====>上线啦",newUser.id,newUser.name)
		message <- loginInfo



		buf :=make([]byte,1024)
		cnt,err:=conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read Err:",err)

		}
		fmt.Println("服务端接收到客户发送过来的数据为：",string(buf[:cnt]),"，长度为：",cnt)
	}
}
