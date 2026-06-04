package gocrud

import "regexp"

var (
	tagWithValRegexp = regexp.MustCompile(`[a-zA-Z0-9_]+:[a-zA-Z0-9_-]+`)
)
