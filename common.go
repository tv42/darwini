package darwini

import "strings"

// segment splits the url path, separating the first segment.
func segment(url string) (string, string) {
	seg := url[1:]
	idx := strings.IndexByte(seg, '/')
	var rest string
	if idx >= 0 {
		seg = seg[:idx]
		rest = url[idx+1:]
	}
	return seg, rest
}
