package tcp_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "tcp"

    "fmt"
    "log"
    "net"
    "os"
    "time"
)

type MockBufferHandler struct {
    hasBeenCalled bool
    response string
}

func (m *MockBufferHandler) HandleBuffer([]byte) ([]byte, error) {
    m.hasBeenCalled = true
    return []byte(m.response), nil
}

var _ = Describe("TcpServer", func() {
    var (
        port string
    )

    BeforeEach(func() {
        port = "8080"
    })

    Describe("Lifecycle", func() {
        It("Works as advertised", func() {
            server := NewTcpServer(nil, port, nil)
            Expect(server.Start()).To(Succeed())
            server.Stop()
        })
    })

    Describe("Has middleware handler installed", func() {
        var (
            logger *log.Logger
            handler BufferHandler
            server TcpServer
            addr string
        )

        BeforeEach(func() {
            logger = log.New(os.Stdout, "", log.LstdFlags)
            handler = new(MockBufferHandler)
            server = NewTcpServer(logger, port, handler)
            Expect(server.Start()).To(Succeed())

            addr = "127.0.0.1"
        })

        AfterEach(func() {
            server.Stop()
            server = nil
        })

        It("Exercises the middleware correctly", func() {
            conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
            Expect(conn).ToNot(BeNil())
            Expect(err).ToNot(HaveOccurred())

            message := "Hello world!"
            response := "Recieved"
            handler.(*MockBufferHandler).response = response

            n, err := conn.Write([]byte(message))
            Expect(n).To(Equal(len(message)))
            Expect(err).ToNot(HaveOccurred())

            time.Sleep(100*time.Millisecond)

            buff := make([]byte, 4096)
            n, err = conn.Read(buff)
            Expect(err).ToNot(HaveOccurred())
            Expect(string(buff[:n])).To(Equal(response))

            conn.Close()
            Expect(handler.(*MockBufferHandler).hasBeenCalled).To(BeTrue())
        })
    })
})