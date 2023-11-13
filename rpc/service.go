package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Args struct {
	X, Y int
}

// ServiceA 自定义一个结构体类型
type ServiceA struct{}

// Add 为ServiceA类型增加一个可导出的Add方法
func (s *ServiceA) Add(args *Args, reply *int) error {
	*reply = args.X + args.Y
	return nil
}

func main() {
	service := new(ServiceA)
	rpc.Register(service) // 注册RPC服务
	l, e := net.Listen("tcp", ":9091")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		fmt.Println("Connect.....")  // 启动时会建立一次连接, 然后client请求会再次执行
		conn, _ := l.Accept()
		rpc.ServeConn(conn)
	}
}
