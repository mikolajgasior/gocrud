package layout

const (
	AuthorizedAndUnauthorized = iota
	AuthorizedOnly
	UnauthorizedOnly
)

type Page struct {
	Path  string
	Title string
	Auth  int
}

type XPageGroup struct {
	Title string
	Pages []*Page
}
