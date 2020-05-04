package contents

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func ReadTree(urlpath string) []byte {
	urlpath = path.Clean(urlpath)
	nodes := strings.Split(urlpath, "/")
	// remove empty node from split "/foo/bar"
	if len(nodes) > 0 && nodes[0] == "" {
		nodes = nodes[1:]
	}

	links := readNode("", nodes)
	return []byte(strings.Join(links, ""))
}

func readNode(curDir string, nodes []string) []string {
	filepath := path.Join(_rootDir, curDir)
	if stat, err := os.Stat(filepath); err != nil {
		log.Fatalf("file path: %s; %v", filepath, err)
	} else if !stat.IsDir() {
		return []string{}
	}
	files, err := ioutil.ReadDir(filepath)
	if err != nil {
		log.Fatalf(`dirname: %s; %v`, curDir, err)
	}
	// range []os.FileInfo
	var links, sublinks []string
	subnode := ""
	if len(nodes) > 0 {
		subnode = nodes[0]
		sublinks = readNode(path.Join(curDir, nodes[0]), nodes[1:])
	}
	isActive := false
	if len(nodes) == 0 {
		isActive = true
	}
	for _, file := range files {
		filename := file.Name()
		if filename[0] == '.' || filename == "_.md" {
			continue
		}
		links = append(links, linkedListItem(curDir, filename, file.IsDir(), isActive))
		if filename == subnode {
			links = append(links, "<ul>")
			links = append(links, sublinks...)
			links = append(links, "</ul>")
		}
	}
	return links
}

func linkedListItem(dirname, basename string, isDir, isActive bool) string {
	filename := basename
	if isDir {
		filename += "/"
	} else {
		filename, _ = fileNameExt(filename)
	}
	class := ""
	if isActive {
		class = `class="active"`
	}
	return fmt.Sprintf(`<a href="/page/%s"><li %s>%s</li></a>`, uriEncode(path.Join(dirname, basename)), class, filename)
}

func fileNameExt(filename string) (readname, ext string) {
	delim := strings.LastIndex(filename, ".")
	if delim == -1 {
		return filename, ""
	}
	return filename[:delim], filename[delim+1:]
}
