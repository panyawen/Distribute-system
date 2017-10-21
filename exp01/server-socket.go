
package main

import (
	"os"
    "net"
    "log"
    "time"
	"fmt"
	"strings"
	"strconv"
	"encoding/binary"
)

var client_time uint64

func Authorized(conn net.Conn) (bool){ //认证用户
	remoteAddr := conn.RemoteAddr() //获取连接到的对象的IP:Port地址
	remoteAddr_split := strings.Split(remoteAddr.String(), ":") //获取远端IP
	remote_ip := strings.Split(remoteAddr_split[0], ".")
	
	localAddr := conn.LocalAddr()  //本地IP：Port
	localAddr_split := strings.Split(localAddr.String(), ":") //本地IP
	local_ip := strings.Split(localAddr_split[0], ".")
	
	//remote_head_ip := strconv.FormatUint(remote_ip[0], 10)
	local_head_ip, _ := strconv.ParseUint(local_ip[0], 10, 8)
	if local_head_ip == 127{  //localhost
		if remoteAddr_split[0] != "127.0.0.1"{
			return false
		}
	}else if local_head_ip < 127{  //A类
		if local_ip[0] != remote_ip[0]{
			return false
		}
	}else if local_head_ip < 192{ //B类
		if local_ip[0] != remote_ip[0] || local_ip[1] != local_ip[1]{
			return false
		}
	}else if local_head_ip < 224{ //C类
		if local_ip[0] != remote_ip[0] || local_ip[1] != remote_ip[1] || 
			local_ip[2] != remote_ip[2]{
			
			return false
		}
	}else{
		return false
	}
	fmt.Println("远程IP", remote_ip[:])
//	fmt.Println(local_ip[:])
	//	fmt.Println(addr[0])
	return true
}

func Handle_delay(conn net.Conn){ //估算传输时延
	buf := make([]byte, 20)
	tb := make([]byte, 20)
	for true{
		n, err := conn.Read(buf)
		if err != nil{
			log.Fatal(err)
			//break
		}
		client_time = binary.BigEndian.Uint64(buf[0:n]) //客户端发送时间，用于计算传输时延
		binary.BigEndian.PutUint64(tb, client_time)
		_, err = conn.Write(tb)  //返回客户端发来的时间，用于估计传输时延
		//println("测试:", n , "  ", client_time)
		if err != nil{
			log.Fatal(err)
		//	break
		}

	}
}

func Handle_conn(conn net.Conn) { //这个是在处理客户端会阻塞的代码。
	tb := make([]byte,20)
	for true{
		t1 := time.Now().UnixNano()

	//	fmt.Println(time.Now().UnixNano())

		binary.BigEndian.PutUint64(tb, uint64(t1))
		var t2 = time.Unix(t1/1e9, 0)
		var _, err = conn.Write(tb)
		if err != nil{
			//log.Fatal(err)
			break
		}
		//println([]byte(t2.Format("2006-01-02 03:04:05")))  //通过conn的wirte方法将这些数据返回给客户端。
		println(t2.Format("2006-01-02 03:04:05 PM"))
		time.Sleep(time.Second)
	}
    conn.Close() //与客户端断开连接。
}

func main() {
    args := os.Args //获取用户输入的所有参数
    if args == nil || len(args) <2{
        fmt.Println("Please input ip and port!")//如果用户没有输入,或参数个数不够,则调用该函数提示用户
        return
    }

	addr := args[1] + ":" + args[2] //ip + port
	fmt.Println("Server start at ", addr)    
	listener,err := net.Listen("tcp",addr)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    for true {
        conn,err := listener.Accept() //用conn接收链接
		println("conn success")
		if err != nil {
            log.Fatal(err)
			continue
        }
		if Authorized(conn){
			go Handle_conn(conn)  //开启多个协程。
			go Handle_delay(conn)
		}else{
			conn.Write([]byte("You are not authorized!"));
			conn.Close()
		}
    }
}
