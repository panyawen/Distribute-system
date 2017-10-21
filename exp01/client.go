package main;
 
import (
    "net/rpc"
	"time"
	"fmt"
    "log"
	"os"
)
 
type Params struct {
   ServerTime int64
   ClientTime int64
}

func main() {
    args := os.Args //获取用户输入的所有参数
    if args == nil || len(args) <3{
        fmt.Println("Please input ip，port and key!")//如果用户没有输入,或参数个数不够,则调用该函数提示用户
        return
    }

	addr := args[1] + ":" + args[2] //ip + port
	fmt.Println("Server:", addr)    
    //连接远程rpc服务
    //这里使用Dial，http方式使用DialHTTP，其他代码都一样
    rpc, err := rpc.Dial("tcp", addr)
	if err != nil {
        log.Fatal(err)
    }

	//调用远程方法
    //注意第三个参数是指针类型
	var authorized bool
	err1 := rpc.Call("Time.Authorize", args[3], &authorized)
	if authorized == false{
		println("认证失败")
		return
	}else{
		println("认证成功")
	}
	
	var ret = Params{}
	if err1 != nil{
		log.Fatal(err1)
	}

	for true{
	    err2 := rpc.Call("Time.GetTime", time.Now().UnixNano(), &ret)
		if err2 != nil {
			log.Fatal(err2)
		}
		delay := time.Now().UnixNano() - ret.ClientTime
		delay /= 2
		ret.ServerTime += delay / 1e9  //用时延更新，得到最新时间
	//	fmt.Print("\r", ret.ServerTime)

		format_time := time.Unix(ret.ServerTime, 0)
		fmt.Print(format_time.Format("\r2006-01-01 03:04:05 PM   "), format_time.Weekday().String())
		fmt.Print("   ", ret.ServerTime)
		fmt.Print("   Delay:", delay, "ns")
		time.Sleep(time.Second)
		
	}
}
