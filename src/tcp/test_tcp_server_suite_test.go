package tcp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestTcpServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestTcpServer Suite")
}
