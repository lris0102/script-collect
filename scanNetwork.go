package main

import (
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
)

// 定义常见端口和服务映射
var commonPorts = map[int]string{
	21:   "FTP",
	22:   "SSH",
	23:   "Telnet",
	80:   "HTTP",
	443:  "HTTPS",
	3306: "MySQL",
	3389: "RDP",
}

// ping 主机以确定其是否在线
func ping(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	return err == nil
}

// scanPort 扫描指定IP的端口并推测服务
func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done() // 在函数结束时减少WaitGroup计数
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	conn.Close()

	service, exists := commonPorts[port]
	if !exists {
		service = "Unknown"
	}

	fmt.Printf("Host %s has port %d open (%s)\n", ip, port, service)
}

// scanNetwork 扫描整个子网并进行端口扫描和日志分析
func scanNetwork(subnet string) {
	var wg sync.WaitGroup // 创建WaitGroup来同步所有并发操作

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", subnet, i)
		wg.Add(1) // 增加WaitGroup计数
		go func(ip string) {
			defer wg.Done() // 在goroutine结束时减少WaitGroup计数
			if ping(ip) {
				fmt.Printf("Host %s is up\n", ip)
				for port := range commonPorts {
					wg.Add(1) // 对每个端口扫描任务增加计数
					go scanPort(ip, port, 1*time.Second, &wg)
				}
				// 检查是否有入侵迹象
				checkLogsForIntrusion(ip, "root", "password") // 需要配置实际的SSH登录信息
			}
		}(ip)
	}

	wg.Wait() // 等待所有的扫描和日志分析任务完成
}

// SSH 连接并分析远程主机的日志
func getSSHClient(ip, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func checkLogsForIntrusion(ip, user, password string) {
	client, err := getSSHClient(ip, user, password)
	if err != nil {
		log.Println("SSH connection failed:", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	output, err := session.CombinedOutput("grep 'Failed password' /var/log/auth.log")
	if err != nil {
		log.Println("Log file access error:", err)
		return
	}

	if strings.Contains(string(output), "Failed password") {
		fmt.Printf("Potential intrusion detected on %s: Failed login attempts found.\n", ip)
	} else {
		fmt.Printf("No intrusion signs found on %s.\n", ip)
	}
}

func main() {
	subnet := "192.168.1" // 子网IP
	scanNetwork(subnet)    // 开始扫描网络
}
