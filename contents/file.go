package contents

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func ReadNode(urlpath string) (bs []byte) {
	urlpath = uriDecode(urlpath)
	filepath := path.Join(_rootDir, urlpath)
	stat, err := os.Stat(filepath)
	if err != nil {
		log.Println(err)
		return err404(urlpath)
	}
	if stat.IsDir() {
		filepath = path.Join(filepath, "_.md")
	}
	if bs, err = ioutil.ReadFile(filepath); err == nil {
		return bs
	}
	// failed to read file
	log.Printf("failed to read file %s\n%v", urlpath, err)
	return err404(urlpath)
}

func err404(urlpath string) []byte {
	return []byte(fmt.Sprintf("<h3>%s</h3>Page %s", path.Base(urlpath), urlpath))
}
