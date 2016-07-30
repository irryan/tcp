package tcp_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "tcp"

    "fmt"
    "net"
)

var _ = Describe("TestTcpServer", func() {
    var (
        port string
    )

    BeforeEach(func() {
        port = "8080"
    })

    Describe("Lifecycle", func() {
        It("Works as advertised", func() {
            server := NewTestTcpServer(port)
            Expect(server.Start()).To(Succeed())
            server.Stop()
        })
    })

    Describe("Receives bytes", func() {
        var (
            addr string
            server *TestTcpServer
        )

        BeforeEach(func() {
            addr = "127.0.0.1"

            server = NewTestTcpServer(port)
            Expect(server.Start()).To(Succeed())
        })

        AfterEach(func() {
            server.Stop()
        })

        It("Receives bytes correctly", func() {
            conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
            Expect(err).ToNot(HaveOccurred())

            n, err := conn.Write([]byte(`Hello world!`))
            Expect(err).ToNot(HaveOccurred())
            fmt.Println(n)
        })
    })
})