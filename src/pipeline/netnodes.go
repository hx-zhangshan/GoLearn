package pipeline

import (
	"net"
	"bufio"
)

/**
	把数据写入 网络中的主机当中就像写进channel一样
 */
func NetWorkSink(addr string, in <-chan int){
	listenr,err:=net.Listen("tcp",addr)
	if err!=nil {
		panic(err)
	}
	go func() {
		defer listenr.Close()
		conn,err:=listenr.Accept()
		if err!=nil {
			panic(err)
		}
		defer  conn.Close()
		writer:=bufio.NewWriter(conn)
		defer writer.Flush()
		WriteSink(writer,in)
	}()
}
/*

 */
func NetWorkResource(addr string) <-chan int{
	out:=make(chan int)
	go func() {
		//连接 网络 从中读取数据
		conn,err:=net.Dial("tcp",addr)
		if err!=nil {
			panic(err)
		}
		defer  conn.Close()
		v:=ReaderSource(bufio.NewReader(conn),-1)
		for  r:=range v {
			out<-r
		}
		defer close(out)
	}()
	return out
}