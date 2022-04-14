package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

const (
	POSTS_DIR     = "posts"
	STYLE_FILE    = "style.css"
	TEMPLATES_DIR = "template"
)

// Locates post files.
// Returns a list of post file names.
func listPosts() (posts []string) {
	files, err := ioutil.ReadDir(POSTS_DIR)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		posts = append(posts, strings.TrimSuffix(file.Name(), ".md"))
	}

	return
}

// Reads markdown post data.
// Returns the text within the post file.
func readPost(post string) (post_data []byte) {
	post_data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", POSTS_DIR, post))
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Converts markdown to HTML.
// Returns HTML output.
func markdownToHTML(markdown []byte) (html []byte) {
	html = blackfriday.MarkdownCommon(markdown)

	return
}

// The index page.
func getIndex(context *gin.Context) {
	posts := listPosts()
	log.Println(posts)

	context.HTML(http.StatusOK, "index.tmpl.html", gin.H{"posts": posts})
}

// A post page.
func getPost(context *gin.Context) {
	post_file := fmt.Sprintf(
		"%s.md", strings.TrimPrefix(context.Param("post"), "/posts/"))
	post := template.HTML(string(markdownToHTML(readPost(post_file))))

	context.HTML(http.StatusOK, "post.tmpl.html", gin.H{"post": post})
}

func main() {
	// Set-up Gin.
	router := gin.Default()
	router.Use(gin.Logger())
	router.Delims("{{", "}}")

	// Set-up templates.
	router.LoadHTMLGlob("./templates/*.tmpl.html")
	router.StaticFile("/style.css", STYLE_FILE)

	// Register endpoints.
	router.GET("/", getIndex)
	router.GET("/posts/:post", getPost)

	router.Run()
}
