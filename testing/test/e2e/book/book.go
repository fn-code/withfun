package book

import (
	"github.com/fn-code/withfun/testing/pkg/book"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Books", func() {
	var foxInSocks, lesMis *book.Book

	ginkgo.BeforeEach(func() {
		lesMis = &book.Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  2783,
		}

		foxInSocks = &book.Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  24,
		}
	})

	ginkgo.Describe("Categorizing books", func() {
		ginkgo.Context("with more than 300 pages", func() {
			ginkgo.It("should be a novel", func() {
				gomega.Expect(lesMis.GetCategory()).To(gomega.Equal(book.CategoryNovel))
			})
		})

		ginkgo.Context("with fewer than 300 pages", func() {
			ginkgo.It("should be a short story", func() {
				gomega.Expect(foxInSocks.GetCategory()).To(gomega.Equal(book.CategoryShortStory))
			})
		})
	})
})
