package main;
 
import (
	"os"
	"fmt"
    "net"
    "log"
	"time"
    "net/rpc"
)
 
var Key string = "key"
var Authorized bool = false

//注意字段必须是导出
type Params struct {
    ServerTime int64
	ClientTime int64
}

type Time struct{}

func chkError(err error) {
	if err != nil{
		log.Fatal(err)
	}
}

func (t *Time) Authorize(k string, ret *bool) error {  //认证
	if k == Key{
		*ret = true
		Authorized = true
	}else{
		*ret = false
		Authorized = false
	}
	return nil
}

func (t *Time) GetTime(ct int64, ret *Params) error {  //获取时间
	if Authorized == false{
		(*ret).ClientTime = ct
		(*ret).ServerTime = 0
		return nil
	}
    (*ret).ClientTime = ct
	(*ret).ServerTime = time.Now().Unix()
	return nil;
}

func main() {
    args := os.Args //获取用户输入的所有参数
    if args == nil || len(args) < 3{
        fmt.Println("Please input ip, port and key!")//如果用户没有输入,或参数个数不够,则调用该函数提示用户
        return
    }

	addr := args[1] + ":" + args[2] //ip + port
	Key = args[3]
	fmt.Println("Server start at ", addr)    
    t := new(Time);
    //注册rpc服务
    rpc.Register(t);
    //获取tcpaddr
    tcpaddr, err := net.ResolveTCPAddr("tcp4", addr);
    chkError(err);
    //监听端口
    tcplisten, err2 := net.ListenTCP("tcp", tcpaddr);
    chkError(err2);
    //死循环处理连接请求
    for {
        conn, err3 := tcplisten.Accept();
		//println("接受连接")
        if err3 != nil {
            continue;
        }
	    //使用goroutine单独处理rpc连接请求
		go rpc.ServeConn(conn);
    }
}
