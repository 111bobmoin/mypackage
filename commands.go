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
)
var vmIP = "192.168.000.000"//虚拟机ip
var password = "xxxxxx"//根据访问的虚拟机密码决定
var username = "xxxxxxx"//根据访问虚拟机的用户名决定
// Parses any other command
func c_default(com string, conn io.Writer) (err error) {
	output, err := executeRemoteCommandWithPassword(vmIP, username, password, com)
	if err != nil {
		_, err = fmt.Fprintf(conn, "command not found: %s\n", com)
		//fmt.Fprintf("执行远程命令时发生错误:", err)
		return
	}
	_, err = fmt.Fprintf(conn, " %s\n", output)
	//_, err = fmt.Fprintf(conn, "command not found: %s\n", com)
	return
}

// Parses the `enable` command.
func c_enable(com string, conn io.Writer) (err error) {
	// split the string into 2 parts, after the `enable` string
	sp := strings.SplitAfter(com, "enable ")
	rsp := fmt.Sprintf("%s\n%s\n%s\n",
		"cd",
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

	_, err = fmt.Fprintf(conn, rsp)
	return
}

// Parses the exit command
func c_exit(com string, conn io.ReadWriteCloser) (err error) {
	_, err = fmt.Fprintf(conn, "bye\n")
	conn.Close()
	return
}

// parses an attempt to execute a file
func c_exec(com string, conn io.Writer) (err error) {
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
