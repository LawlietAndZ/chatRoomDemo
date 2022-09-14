package main

import (
	"fmt"
	"net"
	"strings"
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

	//启动全局唯一的go程，负责监听message 通道，写给所有用户
	go broadcast()

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

	//客户端与服务器建立链接时，会有ip和port，可以当user的id
	clientAddr :=conn.RemoteAddr().String()

	fmt.Println("启动业务")

	newUser := User{
		name: clientAddr,//后续会修改，初始值与id相同
		id: clientAddr,
		msg: make(chan string,10),
	}
	//添加User到map
	allUsers[newUser.id] = newUser

	//启动go程，将message 信息返回给客户端
	go writeBackToCliernt(&newUser,conn)

	//向message写入数据,当用户上线的时，将消息通知所有人
	loginInfo := fmt.Sprintf("[%s]:[%s]=====>login now!!!\n",newUser.id,newUser.name)
	message <- loginInfo
	for  {
		//具体业务逻辑
		buf :=make([]byte,1024)
		cnt,err:=conn.Read(buf)
		if err != nil {
			fmt.Println("conn.Read Err:",err)
		}
		fmt.Println("服务端接收到客户发送过来的数据为：",string(buf[:cnt]),"，长度为：",cnt)
		//-------业务逻辑处理，开始-------
		//1、查询当前所有用户 who
		//a.判断接收的数据是不是who,最后一个是回车
		userInpt :=string(buf[:cnt-1])
		//b.遍历map,将id和username拼接成字符串，返回给客户端
		if len(userInpt) == 4 && userInpt == "\\who"{
			//todo
			fmt.Println("即将查询所有用户信息")
			//这个切片包含所有的用户
			var  userInfos []string
			for _,user := range allUsers{
				userInfo := fmt.Sprintf("userid:%s,username:%s",user.id,user.name)
				userInfos = append(userInfos, userInfo)
			}
			//将userInfos 转化为字符串，写入自身管道中
			newUser.msg <- strings.Join(userInfos,"\n")
		}else if len(userInpt) > 7 && userInpt[0:7] == "\\rename"{
			//重命名
			newUser.name = strings.Split(userInpt,"|")[1]
			fmt.Println(newUser.name+"======================>")
			allUsers[newUser.id] = newUser //更新
			newUser.msg<- "rename sucess"

		}else{
			//如果不是命令，往message中写数据
			message <- userInpt
			//通知客户端

		}
		//-------业务逻辑处理，结束-------

	}
}


func broadcast(){
	fmt.Println("广播go程启动成功")

	for{
		//从message 中读数据
		info := <- message
		fmt.Println(info)
		for _,user := range allUsers{
			user.msg <-info
		}
	}
}
func writeBackToCliernt(user *User,conn net.Conn){

	fmt.Println("用户监听自己的管道")
	for date := range user.msg{
		fmt.Printf("user : %s 写回给客户端的数据为  %s\n",user.name,date)
		_ , _ = conn.Write([]byte(date))
	}

}
