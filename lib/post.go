package blog

const (
	Draft   = 0
	Publish = 1
)

type Post struct {
	Title   string
	Body    []byte
	Posted  string
	Updated string
	Status  int
}
