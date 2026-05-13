package udp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"wails-router-link/service/api/udp/utils/gen"
)

type UDPClient struct {
	conn              *net.UDPConn
	done              chan struct{}
	heartbeatInterval time.Duration

	lastMsg string
	mu      sync.RWMutex
}

func NewUDPClient(serverAddr string, heartbeatInterval time.Duration) (*UDPClient, error) {
	// 解析服务端 UDP 地址
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("解析 UDP 地址失败: %v", err.Error())
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("连接服务端失败: %v", err.Error())
	}

	client := &UDPClient{
		conn:              conn,
		done:              make(chan struct{}),
		heartbeatInterval: heartbeatInterval,
	}
	login_buf := udp_utils_gen.Login()
	fmt.Println("login_buf:", login_buf)

	res, ok := conn.Write(login_buf)
	fmt.Println("res:", res, " ok:", ok)

	// 启动协程
	go client.recvLoop()
	// go client.heartbeatLoop()

	return client, nil
}

// 关闭 UDP 客户端，停止协程并关闭连接
func (c *UDPClient) Close() {
	select {
	case <-c.done:
		// 已经关闭
	default:
		close(c.done)
	}
	_ = c.conn.Close()
}

// 发送一条消息给 UDP 服务端
func (c *UDPClient) Send(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

// 返回最近接收到的一条 UDP 消息
func (c *UDPClient) LastMsg() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastMsg
}

// 通过for循环持续接收服务端消息
func (c *UDPClient) recvLoop() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-c.done:
			fmt.Println("用于消息接收的协程退出")
			return
		default:
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(2 * time.Second))

		n, addr, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			// TODO：读取超时等错误是否需要处理，暂不清楚，此处不做处理，继续循环
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			fmt.Println("接收出错:", err)
			continue
		}

		msg := buf[:n]

		c.mu.Lock()
		c.lastMsg = string(msg)
		c.mu.Unlock()

		fmt.Printf("收到来自 %s 的消息: %+v\n", addr.String(), msg)
	}
}

// // 发送心跳
// func (c *UDPClient) heartbeatLoop() {
// 	ticker := time.NewTicker(c.heartbeatInterval)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-c.done:
// 			fmt.Println("心跳协程退出")
// 			return
// 		case <-ticker.C:
// 			heartbeatMsg := []byte("heartbeat")
// 			_, err := c.conn.Write(heartbeatMsg)
// 			if err != nil {
// 				fmt.Println("发送心跳失败:", err)
// 				continue
// 			}
// 			fmt.Println("已发送心跳包")
// 		}
// 	}
// }
