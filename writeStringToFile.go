package main

import (
	"fmt"
	"io/ioutil"
)

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
