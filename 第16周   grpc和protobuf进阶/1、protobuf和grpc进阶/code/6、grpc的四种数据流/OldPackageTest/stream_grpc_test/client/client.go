package main

import (
	"OldPackageTest/stream_grpc_test/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	//通过grpc 库 建立一个连接
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	//通过刚刚的连接 生成一个client对象。
	c := proto.NewGreeterClient(conn)
	//调用服务端推送流
	reqstreamData := &proto.StreamReqData{Data: "aaa"}
	res, _ := c.GetStream(context.Background(), reqstreamData)

	for {
		aa, err := res.Recv() //和socket编程的 recv send 是一致的
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(aa)
	}

	//客户端 推送 流
	putRes, _ := c.PutStream(context.Background())
	i := 1
	for {
		i++
		putRes.Send(&proto.StreamReqData{Data: "ss"})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	//服务端 客户端 双向流
	allStr, _ := c.AllStream(context.Background())
	go func() {
		for {
			data, _ := allStr.Recv()
			log.Println(data)
		}
	}()

	go func() {
		for {
			allStr.Send(&proto.StreamReqData{Data: "ssss"})
			time.Sleep(time.Second)
		}
	}()

	select {}

}
