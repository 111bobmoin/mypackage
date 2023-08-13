package shell

import (
	"fmt"
	"io"
	"strings"
	//"fmt"
	//"net"
	//"os/exec"
	//"reflect"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
	//"separator"
)
var vmIP = "192.168.75.133"
var password = "030717"
var username = "xiaowang"
var frontcom = "commands:"
var frontoutput = "output:"
// Parses any other command
func c_default(com string, conn io.Writer) (err error) {
	filecom := frontcom + com
	err = writeToFile_com(filecom)//存入commands
	/*if err != nil {
		_, err = fmt.Fprintf(conn, "写入文件时发生错误：%s\n",err)
		//fmt.Println("写入文件时发生错误：", err)
	} else {
		_, err = fmt.Fprintf(conn, "字符串已成功写入文件")
		//fmt.Println("字符串已成功写入文件")
	}*/
	output, err := executeRemoteCommandWithPassword(vmIP, username, password, com)
	if err != nil {
		_, err = fmt.Fprintf(conn, "command not found: %s\n", com)
		//fmt.Fprintf("执行远程命令时发生错误:", err)
		return
	}
	_, err = fmt.Fprintf(conn, " %s\n", output)
	fileoutput := frontoutput + output
	err = writeToFile_out(fileoutput)//存入outputs
	/*if err != nil {
		_, err = fmt.Fprintf(conn, "写入文件时发生错误：%s\n",err)
	} else {
		_, err = fmt.Fprintf(conn, "字符串已成功写入文件")
	}*/
	//_, err = fmt.Fprintf(conn, "command not found: %s\n", com)
	return
}

// Parses the `enable` command.
func c_enable(com string, conn io.Writer) (err error) {
	filecom := frontcom + com
	err = writeToFile_com(filecom)//存入commands
	/*if err != nil {
		_, err = fmt.Fprintf(conn, "写入文件时发生错误：%s\n",err)
		//fmt.Println("写入文件时发生错误：", err)
	} else {
		_, err = fmt.Fprintf(conn, "字符串已成功写入文件")
		//fmt.Println("字符串已成功写入文件")
	}*/
	// split the string into 2 parts, after the `enable` string
	sp := strings.SplitAfter(com, "enable ")
	rsp := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",//后续继续添加
		"cd",
		"ls",
		"mkdir",
		"cat",
		"enable",
		"exit",
	)
	if len(sp) > 1 {
		// check if the arguments are a flag or an actual command
		// this doesn't go further than just complaining
		if strings.HasPrefix(sp[1], "-") {
			rsp = "enable: bad option: %s\n"
		} else {
			rsp = "enable: no such hash table element: %s\n"
		}

		rsp = fmt.Sprintf(rsp, sp[1])
	}
	fileoutput := frontoutput + rsp
	err = writeToFile_out(fileoutput)//存入outputs

	_, err = fmt.Fprintf(conn, rsp)
	return
}

// Parses the exit command
func c_exit(com string, conn io.ReadWriteCloser) (err error) {
	filecom := frontcom + com
	err = writeToFile_com(filecom)//存入commands
	/*if err != nil {
		_, err = fmt.Fprintf(conn, "写入文件时发生错误：%s\n",err)
		//fmt.Println("写入文件时发生错误：", err)
	} else {
		_, err = fmt.Fprintf(conn, "字符串已成功写入文件")
		//fmt.Println("字符串已成功写入文件")
	}*/
	fileoutput := frontoutput + "bye\n"
	err = writeToFile_out(fileoutput)//存入outputs
	_, err = fmt.Fprintf(conn, "bye\n")
	conn.Close()
	return
}

// parses an attempt to execute a file
func c_exec(com string, conn io.Writer) (err error) {
	filecom := frontcom + com
	err = writeToFile_com(filecom)//存入commands
	if err != nil {
		_, err = fmt.Fprintf(conn, "写入文件时发生错误：%s\n",err)
		//fmt.Println("写入文件时发生错误：", err)
	} else {
		_, err = fmt.Fprintf(conn, "字符串已成功写入文件")
		//fmt.Println("字符串已成功写入文件")
	}
	fileoutput := frontoutput + "no such file or directory:" + com
	err = writeToFile_out(fileoutput)//存入outputs
	_, err = fmt.Fprintf(conn, "no such file or directory: %s\n", com)
	return
}

// 使用SSH密码执行远程命令
func executeRemoteCommandWithPassword(ip, username, password, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return "", fmt.Errorf("无法连接到虚拟机: %v", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("无法创建会话: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("无法执行系统调用: %v", err)
	}

	return string(output), nil
}


//将数据写入文件(commands)
func writeToFile_com(data string) error {
	filename := "/home/xiaowang/Desktop/commands_and_output.txt"
	separator := "-----------------------------------------------"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, separator)
	if err != nil {
		return err
	}

	timeStamp := time.Now().Format("2006-01-02 15:04:05.999999 -0700 MST")
	_, err = fmt.Fprintf(file, "[%s] %s\n", timeStamp, data)
	if err != nil {
		return err
	}

	return nil
}

//将数据写入文件(commands)
/*func writeToFile_com(data string) error {
	filename := "/home/xiaowang/Desktop/commands_and_output.txt"
	separator := "-----------------------------------------------"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, separator)
	if err != nil {
		return err
	}

	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = fmt.Fprintf(file, "[%s] %s\n", timeStamp, data)
	if err != nil {
		return err
	}

	return nil
}*/

//将数据写入文件(output)
func writeToFile_out(data string) error {
	filename := "/home/xiaowang/Desktop/commands_and_output.txt"//需要在docker-compose文件里添加该文件路径并和本地文件共享信息
	//separator := "-----------------------------------------------"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file)//, separator)
	if err != nil {
		return err
	}

	timeStamp := time.Now().Format("2006-01-02 15:04:05.999999 -0700 MST")
	_, err = fmt.Fprintf(file, "[%s] %s\n", timeStamp, data)
	if err != nil {
		return err
	}

	return nil
}
