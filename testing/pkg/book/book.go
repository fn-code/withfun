package book

const (
	CategoryNovel = iota
	CategoryShortStory
)

type Book struct {
	Title  string
	Author string
	Pages  int
}

func (b *Book) GetCategory() int {
	if b.Pages <= 300 {
		return CategoryShortStory
	}

	return CategoryNovel
}
