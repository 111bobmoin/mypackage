package main

import (
	"fmt"
	"net"
	"os/exec"
	"reflect"
)

// 访问连接到的虚拟机的系统调用
func accessConnectedVMSystem(vmIP, command, password string) (string, error) {
	conn, err := connectToVM(vmIP, password)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	cmd := exec.Command("sshpass", "-p", password, "ssh", "-oStrictHostKeyChecking=no", vmIP, command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// 连接到虚拟机
func connectToVM(ip, password string) (net.Conn, error) {
	conn, err := net.Dial("tcp", ip+":22")
	if err != nil {
		return nil, err
	}

	// 发送密码进行认证
	_, err = conn.Write([]byte(password + "\n"))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func system_call(vmIP, command, password string) {
	systemCallOutput, err := accessConnectedVMSystem(vmIP, command, password)
	if err != nil {
		fmt.Println("无法访问连接到的虚拟机的系统调用:", err)
		return
	}
	err := writeStringToFile("output.txt", systemCallOutput)
	if err != nil {
		fmt.Println("无法写入文件：", err)
		return
	}
}

func writeStringToFile(filename, content string) error {
	// 将字符串内容转换为字节数组
	data := []byte(content)

	// 将字节数组写入txt文件
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
