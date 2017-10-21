package main

import (
    "os"
    "net"
	"time"
	"log"
    "fmt"
	"encoding/binary"
)

var delay int64

func Handle_delay(conn net.Conn){
	tb := make([]byte, 8)
	for true{
		int_t := time.Now().UnixNano() - 100000000000  //减去1000...，为了区分返回的服务器时间
//		println("测试：", int_t)
		binary.BigEndian.PutUint64(tb, uint64(int_t))
		// string_t := time.Unix(int
		_, err := conn.Write(tb)
		if err != nil{
			break
		}
		time.Sleep(time.Second)
	}
}

func main() {
    args := os.Args //获取用户输入的所有参数
    if args == nil || len(args) <2{
        fmt.Println("Please input server's ip and port!")//如果用户没有输入,或参数个数不够,则调用该函数提示用户
        return
    }
    addr := args[1] + ":" + args[2]
    conn,err := net.Dial("tcp",addr) //拨号操作，需要指定协议。
    if err != nil {
        log.Fatal(err)
    }
	go Handle_delay(conn) //处理传输时延
	//println("此处指出现一次")
	buf := make([]byte, 20)
	for true{
		n,err := conn.Read(buf) //n接受返回的数据大小，用err接受错误信息。
	    if err != nil {
		    log.Fatal(err)
	    }
		var int_t = int64(binary.BigEndian.Uint64(buf[0:n]))
		if(int_t < time.Now().UnixNano() - 100000000000){ //用于估计传输时延
			delay = time.Now().UnixNano() - int_t - 100000000000
			println("传输时延：", delay / 2, "ns")
		}else{	//服务器时间
			fmt.Println(int_t / 1e9, "s") //将接受的内容都读取出来。
			var string_t = time.Unix((int_t + delay / 2) / 1e9, 0)
			fmt.Println(string_t.Format("2006-01-01 03:04:05 AM"))
			println()
		}
		//println()
	}
    conn.Close()  //断开TCP链接。
}
