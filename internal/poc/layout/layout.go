package layout

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"text/template"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

var ExecuteTemplateError = errors.New("error executing template")

type Layout struct {
	HTML     embed.FS
	sitemaps []*Sitemap
}

func (l *Layout) Render(uri string, userID, userName string, content string) ([]byte, error) {
	configCss, stylesCss := l.css()
	pageHTMLTemplate := l.pageHTML("layout")

	// path, title, path, title, ..
	var pages []*Page

	for _, sitemap := range l.sitemaps {
		sortedPageKeys := []string{}
		for page := range sitemap.Pages {
			sortedPageKeys = append(sortedPageKeys, page)
		}
		sort.Strings(sortedPageKeys)

		for _, page := range sortedPageKeys {
			pages = append(pages, sitemap.Pages[page])
		}
	}

	tplObj := struct {
		URI                           string
		ConfigCSS                     string
		StylesCSS                     string
		Username                      string
		UserID                        string
		Pages                         []*Page
		Content                       string
		PageAuthorizedAndUnauthorized int
		PageAuthorizedOnly            int
		PageUnauthorizedOnly          int
	}{
		URI:                           uri,
		ConfigCSS:                     string(configCss),
		StylesCSS:                     string(stylesCss),
		Username:                      userName,
		UserID:                        userID,
		Pages:                         pages,
		Content:                       content,
		PageAuthorizedAndUnauthorized: AuthorizedAndUnauthorized,
		PageAuthorizedOnly:            AuthorizedOnly,
		PageUnauthorizedOnly:          UnauthorizedOnly,
	}

	buf := &bytes.Buffer{}
	t := template.Must(template.New("layout").Parse(string(pageHTMLTemplate)))
	err := t.Execute(buf, &tplObj)
	if err != nil {
		slog.Error("error executing template for home page", logger.AttrError(err))
		return []byte{}, ExecuteTemplateError
	}

	return buf.Bytes(), nil
}

func (l *Layout) AddSitemap(sitemap *Sitemap) {
	l.sitemaps = append(l.sitemaps, sitemap)
}

func (l *Layout) css() ([]byte, []byte) {
	configCSS, _ := embed.FS.ReadFile(l.HTML, "html/config.css")
	stylesCSS, _ := embed.FS.ReadFile(l.HTML, "html/styles.css")

	return configCSS, stylesCSS
}

func (l *Layout) pageHTML(page string) []byte {
	pageTemplate, _ := embed.FS.ReadFile(l.HTML, fmt.Sprintf("html/%s.html", page))

	return pageTemplate
}
