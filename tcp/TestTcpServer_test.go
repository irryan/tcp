package tcp_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "tcp"

    "fmt"
    "net"
    "reflect"
    "time"
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
            server = nil
        })

        Context("There are failed expectations", func() {
            BeforeEach(func() {
                server.AddExpectation(func([]byte) error {
                    return fmt.Errorf("Error occurred!")
                })
            })

            It("Reports failed expectations correctly", func() {
                conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
                Expect(err).ToNot(HaveOccurred())

                b := []byte(`Hello world!`)

                n, err := conn.Write(b)
                Expect(n).To(BeNumerically("==", len(b)))
                Expect(err).ToNot(HaveOccurred())

                time.Sleep(100*time.Millisecond)

                Expect(server.HasFailedExpectations()).To(BeTrue())
                Expect(server.HasRemainingExpectations()).To(BeFalse())
            })
        })

        Context("There are remaining expectations", func() {
            BeforeEach(func() {
                server.AddExpectation(func([]byte) error {
                    return nil
                })
            })

            It("Reports remaining expectations correctly", func() {
                Expect(server.HasFailedExpectations()).To(BeFalse())
                Expect(server.HasRemainingExpectations()).To(BeTrue())
            })
        })

        It("Receives bytes correctly", func() {
            conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
            Expect(err).ToNot(HaveOccurred())

            b := []byte(`Hello world!`)

            server.AddExpectation(func(buff []byte) error {
                if reflect.DeepEqual(buff[:len(b)], b) {
                    return nil
                }

                return fmt.Errorf("Values are not equal.")
            })

            n, err := conn.Write(b)
            Expect(n).To(BeNumerically("==", len(b)))
            Expect(err).ToNot(HaveOccurred())

            time.Sleep(100*time.Millisecond)

            Expect(server.HasFailedExpectations()).To(BeFalse())
            Expect(server.HasRemainingExpectations()).To(BeFalse())
        })

        Context("With multiple expectations", func() {
            var (
                conn net.Conn
                b []byte
            )

            BeforeEach(func() {
                var err error
                conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
                Expect(err).ToNot(HaveOccurred())

                b = []byte(`Hello world!`)

                server.AddExpectation(func(buff []byte) error {
                    if reflect.DeepEqual(buff[:len(b)], b) {
                        return nil
                    }

                    return fmt.Errorf("Values are not equal.")
                })

                server.AddExpectation(func(buff []byte) error {
                    if reflect.DeepEqual(buff[:len(b)], b) {
                        return nil
                    }

                    return fmt.Errorf("Values are not equal.")
                })
            })

            FIt("Works", func() {
                n, err := conn.Write(b)
                Expect(n).To(BeNumerically("==", len(b)))
                Expect(err).ToNot(HaveOccurred())

                n, err = conn.Write(b)
                Expect(n).To(BeNumerically("==", len(b)))
                Expect(err).ToNot(HaveOccurred())

                time.Sleep(100*time.Millisecond)

                Expect(server.HasFailedExpectations()).To(BeFalse())
                Expect(server.HasRemainingExpectations()).To(BeFalse())
            })
        })
    })
})