package server

import (
	"fmt"
	"net"
	"sync"
	syncatomic "test/lib/sync"
	"test/lib/wait"
	"time"
)

type Client struct {
	Conn         net.Conn
	WaitingReply wait.Wait //主要作用是等待一会，让回复处理完；
	uploading    syncatomic.AtomicBool

	// multi bulk msg lineCount - 1(first line)
	expectedArgsCount uint32

	// sent line count, exclude first line
	receivedCount uint32

	// sent lines, exclude first line
	args [][]byte

	// lock while server sending response
	mu sync.Mutex

	// subscribing channels
	subs map[string]bool
}

//新建客户连接；
func NewClient(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

//关闭客户连接
func (c *Client) Close() error {
	c.WaitingReply.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

//向客户连接写消息
func (c *Client) Write(msg []byte) error {
	if len(msg) == 0 || msg == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.Conn.Write(msg)
	if err != nil {
		fmt.Println("服务器写回消息错误！")
		return err
	}
	return nil
}
