package token

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token Service", func() {

	It("Token Error", func() {
		By("getting /v1/token")
		{
			// equal
			//gomega.ExpectWithOffset(1, 1).To(gomega.Equal(2), "wah error ya")

			// not equal
			//gomega.ExpectWithOffset(1, 1).NotTo(gomega.Equal(1), "must not equal")

			// expected error
			//gomega.ExpectWithOffset(1, nil).To(gomega.HaveOccurred(), "wah error bro")

			// ExpectEmpty expects actual is empty
			//gomega.ExpectWithOffset(1, actual).To(gomega.BeEmpty(), explain...)
		}
	})

	It("Token Expired", func() {
		By("getting /v1/token2")
		{

			//gomega.Expect(true).Should(gomega.BeEquivalentTo(true))
			Expect(createError()).Should(BeEquivalentTo(ErrToShort))
			//ginkgo.Skip("expedted got this")
		}
	})
})

var (
	ErrToShort = errors.New("is to short")
)

func createError() error {
	return ErrToShort
}
