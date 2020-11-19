// 与ftp进行连接及上传文件核心功能实现
package run

import (
	"fmt"
	"github.com/SGchuyue/logger/logger"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var Allsize int // 总文件大小
var Num int     // 上传文件数

// connect 建立本地与远程的连接，提供用户名和密码，ip和端口号
func Connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// UploadFile 选择文件上传到服务器中
func UploadFile(sftpClient *sftp.Client, filename string, remotePath string) (size int) {
	logger.Debug("本地上传文件路径为：", filename)
	srcFile, err := os.Open(filename)
	if err != nil {
		logger.Error("打开本地文件路径失败: ", err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(filename)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		logger.Error("远程文件创建失败: ", err)
	}
	defer dstFile.Close()
	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		logger.Error("读取文件%s全部内容失败: ", err)
	}
	dstFile.Write(ff)
	size = len(ff)
	logger.Debug(filename + "已经成功上传到服务器中.")
	return
}

// CountSizeTime 统计上传一个文件的大小和时间
func CountSizeTime(filename, remotepath, user, password, host string, port int) (int, string) {
	var result int
	conn, err := Connect(user, password, host, port)
	defer conn.Close()
	if err != nil {
		logger.Error("与远程建立连接失败：", err)
	}
	start := time.Now()
	result = UploadFile(conn, filename, remotepath)
	spend := time.Since(start).String()
	return result, spend
}

// GetAllFiles 扫描本地文件夹，将.wav文件存到一个有序数组中
func GetAllFiles(hostdir string) (files []string, err error) {
	var dirs []string
	allfile, err := ioutil.ReadDir(hostdir)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator) // 获取文件分隔符
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range allfile {
		if fi.IsDir() {
			dirs = append(dirs, hostdir+PthSep+fi.Name())
			GetAllFiles(hostdir + PthSep + fi.Name())
		} else { // 获取指定格式
			ok := strings.HasSuffix(fi.Name(), ".mp3")
			if ok {
				files = append(files, hostdir+PthSep+fi.Name())
			}
		}
	}
	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}
	return files, nil
}

// WavToMp3 实现将wav格式转化为MP3模式
/*func WavToMp3(wavname string) {
	wav_file := wavname // 需要转换的wav文件
	mp3_file := "C:\\testftp"  // 转换后mp3文件存放路径
	cmd := exec.Command("C:\\Users\\12962\\Desktop\\lame.exe", wav_file, mp3_file)
	err := cmd.Run()
	if err != nil {
	      log.Fatal("Convert wav to mp3 error. ", err)
	     return
	}
	// wav转mp3成功后，如有必要则可删除wav原文件
	//os.Remove(wav_file)
}*/

// Producer 生产者管道
func Producer(ch chan string, localdir string) {
	Filename, err := GetAllFiles(localdir)
	if err != nil {
		logger.Error("获取所有的.wav文件失败：", err)
	}
	for _, file := range Filename {
		ch <- file
	}
	close(ch)
}

// Consumer 消费者管道
func Consumer(ch chan string, done chan bool, user, password, host string, port int, remotepath string) {
	var kb int
	var akb int
	for {
		filename, ok := <-ch
		if ok {
			logger.Debug("成功读取到文件：", filename)
		} else {
			logger.Error("从管道读取文件发生错误。")
			break
		}
		size, spend := CountSizeTime(filename, remotepath, user, password, host, port)
		Num++
		kb = size / 1024
		fmt.Printf("成功上传文件：%s\n使用了：%s\n文件大小为%dKB\n", filename, spend, kb)
		akb += kb
		Allsize = akb / 1024
	}
	done <- true
}

//------强转类型不适合当前场景，精度损失太大，还是很感谢同事的建议-----
/*spend = spend[:len(spend)-2]
spend_float,_ := strconv.ParseFloat(spend,64)
spend_int := int(spend_float)
*/
