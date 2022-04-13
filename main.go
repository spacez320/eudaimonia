package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
)

const (
	TEMPLATES_DIR = "template"
	POSTS_DIR     = "posts"
)

// Locates post files.
// Returns a list of post file names.
func listPosts(posts_dir string) (posts []string) {
	files, err := ioutil.ReadDir(posts_dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		posts = append(posts, file.Name())
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
	posts := listPosts(POSTS_DIR)

	context.HTML(http.StatusOK, "index.tmpl.html", gin.H{
		"posts": posts,
	})
}

// A post page.
func getPost(context *gin.Context) {
	post := template.HTML(
		string(markdownToHTML(readPost(context.Param("post")))))

	context.HTML(http.StatusOK, "post.tmpl.html", gin.H{
		"post": post,
	})
}

func main() {
	// Set-up Gin.
	router := gin.Default()
	router.Use(gin.Logger())
	router.Delims("{{", "}}")

	// Set-up templates.
	router.LoadHTMLGlob("./templates/*.tmpl.html")

	// Register endpoints.
	router.GET("/", getIndex)
	router.GET("/:post", getPost)

	router.Run()
}
