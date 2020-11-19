package main

import (
	"fmt"
	"ftp_upload/run"
	"github.com/SGchuyue/logger/logger"
	"time"
)

func main() {
	logger.InitLogger("../uplode.log", 3, 3, 3, false)
	starttime := time.Now()
	ch := make(chan string)
	chnum := 6 // 控制并发量
	done := make(chan bool, chnum)
	// 导入本地文件夹路径
	go run.Producer(ch, "C:\\testmp3")
	for i := 1; i < chnum; i++ {
		// ch,done,用户名，密码，IP地址，端口,远程文件存放路径
		go run.Consumer(ch, done, "", "", "", 0, "/root/testsftpstar")
	}
	for i := 1; i < chnum; i++ {
		<-done
	}
	endtime := time.Now()
	spend := endtime.Unix() - starttime.Unix()
	fmt.Printf("文件上传开始时间:%s\n上传完成时间:%s\n所有文件已全部上传完成,总共使用时间：%d秒！", starttime, endtime, spend)
}
