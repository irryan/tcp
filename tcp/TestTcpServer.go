package tcp

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

type ExpectationHandler func([]byte) error

type TestTcpServer struct {
    port string
    quit chan struct{}
    wg sync.WaitGroup

    ln net.Listener

    expectations []ExpectationHandler
    errs []error
}

func NewTestTcpServer(port string) *TestTcpServer {
    return &TestTcpServer{
        port: port,
        quit: make(chan struct{}),
    }
}

func (s *TestTcpServer) Start() (err error) {
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

                        log.Println(err)
                        return
                    }

                    if len(s.expectations) > 0 {
                        next := s.expectations[0]
                        s.expectations = s.expectations[1:]
                        if err = next(buff); err != nil {
                            s.errs = append(s.errs, err)
                        }
                    }
                }
            }(conn, quit)
        }
    }(s.quit)

    return nil
}

func (s *TestTcpServer) Stop() {
    s.ln.Close()
    close(s.quit)
    s.wg.Wait()
}

func (s *TestTcpServer) AddExpectation(expectation ExpectationHandler) {
    s.expectations = append(s.expectations, expectation)
}

func (s *TestTcpServer) HasFailedExpectations() bool {
    return len(s.errs) > 0
}

func (s *TestTcpServer) HasRemainingExpectations() bool {
    return len(s.expectations) > 0
}