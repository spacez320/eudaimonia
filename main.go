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
	ASSETS_DIR    = "assets"
	IMAGES_DIR    = "images"
	TEMPLATES_DIR = "template"
)

// Retrieves a list of posts.
// Returns a list of post names.
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

// Converts Markdown to HTML.
// Returns HTML output.
func markdownToHTML(markdown []byte) (html []byte) {
	html = blackfriday.MarkdownCommon(markdown)

	return
}

// Reads Markdown post data.
// Returns the text within the post file.
func readPost(post string) (post_data []byte) {
	post_data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", POSTS_DIR, post))
	if err != nil {
		log.Fatal(err)
	}

	return
}

// The index page.
func getIndex(context *gin.Context) {
	posts := listPosts()
	log.Println(posts)

	for k, v := range context.Request.Header {
		log.Printf("%s: %s\n", k, v)
	}

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
	router.Delims("{{", "}}")
	router.LoadHTMLGlob("./templates/*.tmpl.html")
	router.Static("/assets", ASSETS_DIR)
	router.Static("/images", IMAGES_DIR)
	router.Use(gin.Logger())

	// Register endpoints.
	router.GET("/", getIndex)
	router.GET("/posts/:post", getPost)

	// Start the webserver.
	router.Run()
}
