package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ClientArgs struct {
	X, Y int
}

func main() {
	// 建立HTTP连接
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &ClientArgs{10, 20}

	// 异步调用
	var reply2 int
	divCall := client.Go("ServiceA.Add", args, &reply2, nil)
	replyCall := <-divCall.Done // 接收调用结果
	fmt.Println(replyCall.Error)
	fmt.Println(reply2)

	// 同步调用
	var reply int
	err = client.Call("ServiceA.Add", args, &reply)
	if err != nil {
		log.Fatal("ServiceA.Add error:", err)
	}
	fmt.Printf("ServiceA.Add: %d+%d=%d\n", args.X, args.Y, reply)

}
