package contents

import (
	"log"
)

var escape map[rune]string = map[rune]string{
	' ': "%20",
	'#': "%23",
	'%': "%25",
	'&': "%26",
	'+': "%2B",
	'=': "%3D",
	'?': "%3F",
	// '/': "%2F", // not valid in Linux
}

var parse map[string]string = map[string]string{
	"%20": " ",
	"%23": "#",
	"%25": "%",
	"%26": "&",
	"%2B": "+",
	"%3D": "=",
	"%3F": "?",
	// "%2F": "/", // not valid in Linux
}

func uriEncode(raw string) (uri string) {
	for _, v := range raw {
		if s, ok := escape[v]; ok {
			uri += s
		} else {
			uri += string(v)
		}
	}
	return
}

func uriDecode(uri string) (raw string) {
	for i := 0; i < len(uri); i++ {
		if uri[i] != '%' {
			raw += string(uri[i])
			continue
		}
		if v, ok := parse[uri[i:i+3]]; ok {
			raw += v
			i += 2
		} else {
			log.Printf("bad encoded uri: %s", uri)
		}
	}
	return
}
