package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"test/db"
	syncatomic "test/lib/sync"
	"test/redis/reply"
)

var (
	UnknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type Handler struct {
	ActiveConn sync.Map
	db         *db.DB
	closing    syncatomic.AtomicBool //控制分配器的开关
}

func NenHandler() *Handler {
	db := db.NewDB()
	return &Handler{db: db}
}

func (h *Handler) CloseConn(cli *Client) {
	cli.Close()
	h.ActiveConn.Delete(cli)
}
func (h *Handler) Close() error {
	//关闭db；

	h.closing.Set(true)
	h.ActiveConn.Range(func(key, value interface{}) bool {
		client := key.(*Client)
		_ = client.Close()
		return true
	})
	return nil
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		return
	}
	client := NewClient(conn)
	h.ActiveConn.Store(client, 1)
	reader := bufio.NewReader(conn)
	var fixedLen int64 = 0
	var err error
	var msg []byte
	for {
		//这个if用来进行读取；
		if fixedLen == 0 {
			msg, err = reader.ReadBytes('\n') //按照行读取 \r\n
			if err != nil {                   //错误处理
				if err == io.EOF ||
					err == io.ErrUnexpectedEOF ||
					strings.Contains(err.Error(), "use of closed network connection") {
					fmt.Println("connection close")
				} else {
					fmt.Println(err)
				}

				// after client close
				h.CloseConn(client)
				return // io error, disconnect with client
			}
			if len(msg) == 0 || msg[len(msg)-2] != '\r' { //读到的消息有错
				errReply := &reply.ProtocolErrReply{Msg: "invalid multibulk length"}
				_, _ = client.Conn.Write(errReply.ToBytes())
			}
		} else {
			msg = make([]byte, fixedLen+2)
			_, err = io.ReadFull(reader, msg)
			if err != nil {
				if err == io.EOF ||
					err == io.ErrUnexpectedEOF ||
					strings.Contains(err.Error(), "use of closed network connection") {
					fmt.Println("connection close")
				} else {
					fmt.Println(err)
				}

				// after client close
				h.CloseConn(client)
				return // io error, disconnect with client
			}
			if len(msg) == 0 ||
				msg[len(msg)-2] != '\r' ||
				msg[len(msg)-1] != '\n' {
				errReply := &reply.ProtocolErrReply{Msg: "invalid multibulk length"}
				_, _ = client.Conn.Write(errReply.ToBytes())
			}
			fixedLen = 0
		}
		//这个if用来进行解析
		if !client.uploading.Get() {
			// new request
			if msg[0] == '*' { //是最开头；
				// bulk multi msg  //解析出数量
				expectedLine, err := strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
				if err != nil {
					_, _ = client.Conn.Write(UnknownErrReplyBytes)
					continue
				}
				client.WaitingReply.Add(1)
				client.uploading.Set(true)
				client.expectedArgsCount = uint32(expectedLine)
				client.receivedCount = 0
				client.args = make([][]byte, expectedLine)
			} else {
				// text protocol
				// remove \r or \n or \r\n in the end of line
				str := strings.TrimSuffix(string(msg), "\n")
				str = strings.TrimSuffix(str, "\r")
				strs := strings.Split(str, " ")
				args := make([][]byte, len(strs))
				for i, s := range strs {
					args[i] = []byte(s)
				}

				// send reply
				result := h.db.Exec(client, args)
				if result != nil {
					_ = client.Write(result.ToBytes())
				} else {
					_ = client.Write(UnknownErrReplyBytes)
				}
			}
		} else {
			// receive following part of a request
			line := msg[0 : len(msg)-2]
			if line[0] == '$' {
				fixedLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
				if err != nil {
					errReply := &reply.ProtocolErrReply{Msg: err.Error()}
					_, _ = client.Conn.Write(errReply.ToBytes())
				}
				if fixedLen <= 0 {
					errReply := &reply.ProtocolErrReply{Msg: "invalid multibulk length"}
					_, _ = client.Conn.Write(errReply.ToBytes())
				}
			} else { //加入参数中
				client.args[client.receivedCount] = line
				client.receivedCount++
			}

			// if sending finished   读完了
			if client.receivedCount == client.expectedArgsCount {
				client.uploading.Set(false) // finish sending progress

				// send reply   发送回复
				result := h.db.Exec(client, client.args)
				if result != nil {
					_ = client.Write(result.ToBytes())
				} else {
					_ = client.Write(UnknownErrReplyBytes)
				}

				// finish reply
				client.expectedArgsCount = 0
				client.receivedCount = 0
				client.args = nil
				client.WaitingReply.Done()
			}
		}

	}
}
