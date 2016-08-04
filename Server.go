package tcp

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

type TcpServer interface {
    Start() error
    Stop()
}

type BufferHandler interface {
    HandleBuffer([]byte) ([]byte, error)
}

func NewTcpServer(logger *log.Logger, port string, handler BufferHandler) TcpServer {
    return &tcpServer{
        port: port,
        handler: handler,
        quit: make(chan struct{}),
    }
}

type tcpServer struct {
    logger log.Logger
    port string
    handler BufferHandler

    ln net.Listener
    quit chan struct{}
    wg sync.WaitGroup
}

func (s *tcpServer) Start() (err error) {
    s.ln, err = net.Listen("tcp", fmt.Sprintf(":%s", s.port))
    if err != nil {
        return
    }

    s.wg.Add(1)
    go func(quit chan struct{}) {
        defer s.wg.Done()
        for {
            select {
                case <-quit:
                    return

                default:
            }

            conn, err := s.ln.Accept()
            if err != nil {
                continue
            }

            s.wg.Add(1)
            go func(conn net.Conn, quit chan struct{}) {
                defer s.wg.Done()
                for {
                    select {
                        case <-quit:
                            return

                        default:
                    }

                    conn.SetDeadline(time.Now().Add(1e9))
                    buff := make([]byte, 4096)
                    if _, err := conn.Read(buff); err != nil {
                        if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                            continue
                        }

                        s.logger.Println(err)
                        return
                    }

                    resp, err := s.handler.HandleBuffer(buff)
                    if err != nil {
                        s.logger.Println(err)
                        continue
                    }

                    if _, err := conn.Write(resp); err != nil {
                        s.logger.Println(err)
                    }
                }
            }(conn, quit)
        }
    }(s.quit)

    return nil
}

func (s *tcpServer) Stop() {
    s.ln.Close()
    close(s.quit)
    s.wg.Wait()
}