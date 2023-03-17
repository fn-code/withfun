package e2e_test

import (
	"testing"

	_ "github.com/fn-code/withfun/testing/test/e2e/auth"
	_ "github.com/fn-code/withfun/testing/test/e2e/book"
	_ "github.com/fn-code/withfun/testing/test/e2e/service"
	_ "github.com/fn-code/withfun/testing/test/e2e/token"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}
