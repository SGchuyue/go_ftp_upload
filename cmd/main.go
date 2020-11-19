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
	chnum := 2 // 控制并发量
	done := make(chan bool, chnum)
	// 导入本地文件夹路径
	go run.Producer(ch, "C:\\test200k")
	for i := 1; i < chnum; i++ {
		// ch,done,用户名，密码，IP地址，端口,远程文件存放路径
		go run.Consumer(ch, done, "", "", "", 0, "/root/testsftpstar")
	}
	for i := 1; i < chnum; i++ {
		<-done
	}
	endtime := time.Now()
	fmt.Printf("一共上传了%d个文件，总共文件大小为%dM\n", run.Num, run.Allsize)
	spend := endtime.Unix() - starttime.Unix()
	spent := int(spend)
	vm := run.Allsize / spent
	vk := run.Allsize * 1024 / spent
	fmt.Printf("文件上传开始时间:%s\n文件上传完成时间:%s\n所有文件已全部上传完成!\n总共使用时间：%d秒\n总共上传大小%dM\n平均上传速度%dM/s;%dKB/s\n", starttime, endtime, spend, run.Allsize, vm, vk)
}
