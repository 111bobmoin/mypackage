package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"


	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// 解析 IP 层数据，返回源 IP 和目的 IP
func parseIP(packet gopacket.Packet) (string, string) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		return ip.SrcIP.String(), ip.DstIP.String()
	}
	return "", ""
}

// 解析传输层数据，返回源端口、目的端口、传输数据和协议类型
func parseTransport(packet gopacket.Packet) (string, string, string, string) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		return tcp.SrcPort.String(), tcp.DstPort.String(), string(tcp.Payload), "TCP"
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		return udp.SrcPort.String(), udp.DstPort.String(), string(udp.Payload), "UDP"
	}

	return "", "", "", ""
}

// 获取时间戳
func getTimestamp(packet gopacket.Packet) string {
	return packet.Metadata().Timestamp.String()
}

// 封装的函数，根据指定的 wireshark 过滤语句返回信息
func processPcapFile(fileName string, filterExpression string) error {
	// 打开 pcap 文件
	handle, err := pcap.OpenOffline(fileName)
	if err != nil {
		return err
	}
	defer handle.Close()

	// 设置过滤器
	if err := handle.SetBPFFilter(filterExpression); err != nil {
		return err
	}

	// 统计 srcip 和 dstport 分类数量的 map
	srcIPCount := make(map[string]int)
	dstPortCount := make(map[string]int)

	// 创建目标文件
	files := map[string]*os.File{
		"srcip.txt":   nil,
		"dstport.txt": nil,
		"time.txt":    nil,
	}

	for file, _ := range files {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()
		files[file] = f
	}

	// 循环读取数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		srcIP, dstIP := parseIP(packet)
		srcPort, dstPort, payload, protocol := parseTransport(packet)
		timestamp := getTimestamp(packet)

		// 写入源IP到文件
		files["srcip.txt"].WriteString("源 IP(分类): " + srcIP + "\n")

		// 写入目的IP到文件
		files["srcip.txt"].WriteString("目的 IP: " + dstIP + "\n")

		// 统计 srcIP 和 dstPort 的分类数量
		srcIPCount[srcIP]++
		dstPortCount[dstPort]++

		// 写入目的端口到文件
		files["dstport.txt"].WriteString("目的端口(分类): " + dstPort + "\n")

		// 写入源端口到文件
		files["dstport.txt"].WriteString("源端口: " + srcPort + "\n")

		// 写入时间戳到文件
		files["time.txt"].WriteString("时间戳(分类): " + timestamp + "\n")

		// 写入其他信息到文件srcip.txt
		dst1 := "srcip.txt"
		files[dst1].WriteString("源端口: " + srcPort + "\n")
		files[dst1].WriteString("目的端口: " + dstPort + "\n")
		files[dst1].WriteString("传输数据信息:\n" + payload + "\n")
		files[dst1].WriteString("协议类型: " + protocol + "\n")
		files[dst1].WriteString("时间戳: " + timestamp + "\n")
		files[dst1].WriteString("-------------------------------------\n")

		// 写入其他信息到文件dstport.txt
		dst2 := "dstport.txt"
		files[dst2].WriteString("源 IP: " + srcIP + "\n")
		files[dst2].WriteString("目的 IP: " + dstIP + "\n")
		files[dst2].WriteString("传输数据信息:\n" + payload + "\n")
		files[dst2].WriteString("协议类型: " + protocol + "\n")
		files[dst2].WriteString("时间戳: " + timestamp + "\n")
		files[dst2].WriteString("-------------------------------------\n")

		// 写入其他信息到文件time.txt
		dst3 := "time.txt"
		files[dst3].WriteString("源 IP: " + srcIP + "\n")
		files[dst3].WriteString("目的 IP: " + dstIP + "\n")
		files[dst3].WriteString("源端口: " + srcPort + "\n")
		files[dst3].WriteString("目的端口: " + dstPort + "\n")
		files[dst3].WriteString("传输数据信息:\n" + payload + "\n")
		files[dst3].WriteString("协议类型: " + protocol + "\n")
		files[dst3].WriteString("-------------------------------------\n")

		// 终端输出
		fmt.Println("源 IP:", srcIP)
		fmt.Println("目的 IP:", dstIP)
		fmt.Println("源端口:", srcPort)
		fmt.Println("目的端口:", dstPort)
		fmt.Println("传输数据信息:", payload)
		fmt.Println("协议类型:", protocol)
		fmt.Println("时间戳:", timestamp)
		fmt.Println("-------------------------------------")
	}

	// 在文件头部插入统计信息和分类后的 IP 和端口
	srcIPCountHeader := fmt.Sprintf("源 IP 分类数量: %d\n", len(srcIPCount))
	dstPortCountHeader := fmt.Sprintf("目的端口 分类数量: %d\n", len(dstPortCount))

	insertAndSortIPPort := func(countMap map[string]int, fileName string) {
		output := fmt.Sprintf("%s\n分类后的 IP 和端口:\n", countMap)
		var items []string
		for item := range countMap {
			items = append(items, item)
		}
		sort.Strings(items)
		for _, item := range items {
			output += fmt.Sprintf("%s: %d\n", item, countMap[item]) // 注意这里使用了 %d
		}
		output += "-------------------------------------\n"
		file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		content, _ := io.ReadAll(file)
		file.Seek(0, 0)
		file.WriteString(srcIPCountHeader + dstPortCountHeader + output + string(content))
	}


	insertAndSortIPPort(srcIPCount, "srcip.txt")
	insertAndSortIPPort(dstPortCount, "dstport.txt")

	return nil
}

func main() {
	var fileName2 string
	
	fileName1 := "/home/xiaowang/Desktop/riotpot/tcpdump/"//根据系统配置更改
	fmt.Print("请输入pcap文件名：")
	fmt.Scanln(&fileName2)
	fileName := fileName1 + fileName2 
	filterExpression := "(tcp dst port 22 or tcp dst port 7 or tcp dst port 23 or tcp dst port 1883 or tcp dst port 80 or tcp dst port 502) and ip dst host 192.168.75.128"//过滤条件语句不建议更改，但是也可以
	err := processPcapFile(fileName, filterExpression)
	if err != nil {
		log.Fatal(err)
	}
}
