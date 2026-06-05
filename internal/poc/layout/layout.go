package layout

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	htmltemplate "html/template"
	"log/slog"

	"codeberg.org/mikolajgasior/gocrud/pkg/logger"
)

var ExecuteTemplateError = errors.New("error executing template")

//go:embed html
var embedHTML embed.FS

type Layout struct {
	sitemaps []*Sitemap
}

func (l *Layout) Render(uri string, userID, userName string, content string) ([]byte, error) {
	styleCSS, scriptJS := l.assets()
	pageHTMLTemplate := l.pageHTML("layout")

	var xpageGroups []*XPageGroup

	for _, sitemap := range l.sitemaps {
		if len(sitemap.XPageGroups) > 0 {
			xpageGroups = append(xpageGroups, sitemap.XPageGroups...)
		}
	}

	tplObj := struct {
		URI                           string
		StyleCSS                      htmltemplate.CSS
		ScriptJS                      htmltemplate.JS
		Username                      string
		UserID                        string
		Content                       htmltemplate.HTML
		PageAuthorizedAndUnauthorized int
		PageAuthorizedOnly            int
		PageUnauthorizedOnly          int
		XPageGroups                   []*XPageGroup
	}{
		URI:                           uri,
		StyleCSS:                      htmltemplate.CSS(styleCSS),
		ScriptJS:                      htmltemplate.JS(scriptJS),
		Username:                      userName,
		UserID:                        userID,
		Content:                       htmltemplate.HTML(content),
		PageAuthorizedAndUnauthorized: AuthorizedAndUnauthorized,
		PageAuthorizedOnly:            AuthorizedOnly,
		PageUnauthorizedOnly:          UnauthorizedOnly,
		XPageGroups:                   xpageGroups,
	}

	buf := &bytes.Buffer{}
	t := htmltemplate.Must(htmltemplate.New("layout").Parse(string(pageHTMLTemplate)))
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

func (l *Layout) assets() ([]byte, []byte) {
	scriptJS, _ := embed.FS.ReadFile(embedHTML, "html/script.js")
	styleCSS, _ := embed.FS.ReadFile(embedHTML, "html/style.css")

	return styleCSS, scriptJS
}

func (l *Layout) pageHTML(page string) []byte {
	pageTemplate, _ := embed.FS.ReadFile(embedHTML, fmt.Sprintf("html/%s.html", page))

	return pageTemplate
}
