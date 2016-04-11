package main

import (
	"bufio"
	"fmt"
	"net"
	"echo/codec"
)

// 用来记录所有的客户端连接
var ConnMap map[string]*net.TCPConn

func main() {
	var tcpAddr *net.TCPAddr
	ConnMap = make(map[string]*net.TCPConn) //初始化
	tcpAddr,_=net.ResolveTCPAddr("tcp","127.0.0.1:9999")

	tcpListener,_:=net.ListenTCP("tcp",tcpAddr) //开启tcp 服务
	//退出时关闭
	defer tcpListener.Close()
	for{
		tcpConn,err :=tcpListener.AcceptTCP()
		if err !=nil {
			continue
		}
		fmt.Println("A client connected : "+ tcpConn.RemoteAddr().String())
		// 新连接加入 map
		ConnMap[tcpConn.RemoteAddr().String()] = tcpConn

		go tcpPipe(tcpConn)
	}
}
//处理发送过来的消息
func tcpPipe(conn *net.TCPConn)  {
	ipStr :=conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected : "+ ipStr)
		conn.Close()
	}()
	//读取数据
	reader :=bufio.NewReader(conn)
	for {
		message ,err :=codec.Decode(reader)//reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Println(string(message))
		//这里返回消息改为广播
		boradcastMessage(conn.RemoteAddr().String()+":"+string(message))
	}
}
//广播给其它 
func boradcastMessage(message string)  {
	//遍历所有客户端并发消息
	for _,conn :=range ConnMap{
		b,err :=codec.Encode(message)
		if err != nil {
			continue
		}
		conn.Write(b)
	}
}
