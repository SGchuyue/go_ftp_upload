package main

import (
	"fmt"
	"ftp_upload/run"
	"github.com/SGchuyue/logger/logger"
	"time"
)

var user string       // 用户名
var password string   // 密码
var host string       // 远程Ip地址
var port int          // 端口号
var localpath string  // 本地录音文件路径
var Remotepath string // 远程存放文件地址

func main() {
	logger.InitLogger("../uplode.log", 3, 3, 3, false)
	ch := make(chan string)
	chnum := 2 // 控制并发量
	done := make(chan bool, chnum)
	fmt.Println("您好，非常感谢您对我们的信赖和使用我们的产品！")
	fmt.Println("请输入本地包含录音文件的文件名！")
	fmt.Scanln(&localpath)
	// run.Localpath = "C:\\testftp\\record"
	fmt.Println("请输入存放Ffpeg的地址文件：")
	fmt.Scanln(&run.Ffmpeg)
	// run.Ffmpeg = "c:/testftp/ffmpeg.exe"
	fmt.Println("请依次输入用户名和密码和Ip地址和端口！(相邻字符用空行隔开)")
	fmt.Scanf("%s %s %s %d", &user, &password, &host, &port)
	starttime := time.Now()
	// 导入本地文件夹路径
	go run.Producer(ch, localpath)
	for i := 1; i < chnum; i++ {
		// ch,done,用户名，密码，IP地址，端口,远程文件存放路径
		go run.Consumer(ch, done, user, password, host, port, Remotepath)
	}
	for i := 1; i < chnum; i++ {
		<-done
	}
	endtime := time.Now()
	fmt.Printf("一共上传了%d个文件，总共文件大小为%dM\n", run.Num, run.Allsize)
	spend := endtime.Unix() - starttime.Unix()
	spent := float64(spend)
	vm := float64(run.Allsize) / spent
	vk := float64(run.Allsize*1024) / spent
	fmt.Printf("文件上传开始时间:%s\n文件上传完成时间:%s\n所有文件已全部上传完成!\n总共使用时间：%d秒\n总共上传大小%dM\n平均上传速度为%fM/s和%fKB/s。\n感谢您的使用！", starttime, endtime, spend, run.Allsize, vm, vk)
	for {
		fmt.Scan()
	}
}
