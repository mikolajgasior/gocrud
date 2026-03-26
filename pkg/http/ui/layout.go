package ui

type Layout interface {
	Render(uri string, userID, userName string, content string) ([]byte, error)
}
