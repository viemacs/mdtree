package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	blackfriday "github.com/russross/blackfriday/v2"

	"viemacs/mdtree/contents"
	_ "viemacs/mdtree/statik"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	const port int = 8080
	const version string = "2020-05-04"

	root := flag.String("r", ".", "Set root directory of markdown tree")
	ver := flag.Bool("v", false, "Print version and exit")
	flag.Parse()
	if *ver {
		fmt.Printf("mdtree %s\n", version)
		return
	}
	if err := contents.SetRoot(*root); err != nil {
		log.Fatal(err)
	}
	log.Printf("mdtree serves at: %s", *root)
	if err := server(port); err != nil {
		log.Fatal(err)
	}
}

func server(port int) error {
	gin.SetMode(gin.ReleaseMode)

	sf := staticfile()
	serveCSSFile := func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "text/css")
		c.Writer.Write(sf(c.Param("filename")))
	}
	serveJSFile := func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/javascript")
		c.Writer.Write(sf(c.Param("filename")))
	}

	router := gin.New()
	router.GET("/css/*filename", serveCSSFile)
	router.GET("/js/*filename", serveJSFile)
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(sf("/favicon.ico"))
	})

	router.GET("/", servePage)
	router.GET("page/*page", servePage)

	contents.ReadTree("/")

	log.Printf("mdtree service starts on: %d", port)
	router.Run(fmt.Sprintf(":%d", port))

	return nil
}

func servePage(c *gin.Context) {
	page := c.Param("page")
	if page == "" {
		page = "/"
	}
	log.Printf("reading page: %s", page)

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(`<!doctype html><html><head>
  <meta charset="utf-8">
  <link rel="stylesheet" href="/css/style.css">
</head><body>
  <h1>Markdown Tree</h1>
<div class="container"><div id="tree"><ul>`))

	c.Writer.Write(contents.ReadTree(page))
	c.Writer.Write([]byte(`</ul></div><div id="page">`))
	c.Writer.Write(blackfriday.Run(contents.ReadNode(page)))
	c.Writer.Write([]byte(fmt.Sprintf(`</div></div>
<div id="footer">%s</div>
</body><script src="/js/render.js"></script></html>`,
		time.Now().Format("2006-01-02"))))
}

func staticfile() func(string) []byte {
	sfs := make(map[string][]byte)

	sf := func(filename string) []byte {
		if val, ok := sfs[filename]; ok {
			return val
		}
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		log.Println(filename)
		r, err := statikFS.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		contents, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		sfs[filename] = contents
		return contents
	}
	return sf
}
