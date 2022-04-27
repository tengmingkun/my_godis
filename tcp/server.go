package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test/interface/hander"
	syncatomic "test/lib/sync"
	"time"
)

type Config struct {
	Address string        `yaml:"address"`
	MaxConn uint32        `yaml:"maxConn"`
	TimeOut time.Duration `yaml:"timeout"`
}

func NewConfig() *Config {
	return &Config{Address: ":8080"}
}

func ListenAndServe(cfg *Config, handler hander.Handler) {
	listen, err := net.Listen("tcp", cfg.Address)
	fmt.Println("ip:120.76.99.69  port:8080")
	if err != nil {
		fmt.Println("服务建立错误")
	}

	//监听信号，实现优雅退出；
	var closeflag syncatomic.AtomicBool
	var channel = make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		msg := <-channel
		switch msg {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			handler.Close()
			listen.Close()
			closeflag.Set(true)
		}
	}()

	//接收事件,并开启协程进行处理
	defer func() {
		handler.Close()
		listen.Close()
	}()
	var waitDone sync.WaitGroup
	ctx, _ := context.WithCancel(context.Background())
	for {
		conn, err := listen.Accept()
		if err != nil {
			if closeflag.Get() {
				waitDone.Wait()
				return
			} else {
				continue
			}
		}
		waitDone.Add(1)
		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn)
		}()
	}
}
