package hander

import (
	"context"
	"fmt"
	"net"
)

type HandleFunc func(ctx context.Context, conn net.Conn)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

type Node struct {
	name string
}

func NewNode() *Node {
	return &Node{name: "123"}
}
func (n Node) Handle(ctx context.Context, conn net.Conn) {
	fmt.Println("yunxing OK")
	return
}

func (n Node) Close() error {
	return nil
}
