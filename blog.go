package main

import (
	"encoding/json"
	"github.com/mfindlater/blog/lib"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BlogController struct {
	Config  blog.Config
	Context blog.BlogContext
	Data    map[string]interface{}
}

func NewBlogController() *BlogController {
	b := new(BlogController)
	b.Context = blog.BlogContext{}
	b.Context.Connect()
	b.Data = make(map[string]interface{})
	return b
}

func (b *BlogController) LoadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	var config blog.Config
	err = json.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
		return err
	}

	b.Config = config
	b.Data["config"] = config

	return nil
}

func (b *BlogController) mainHandler(w http.ResponseWriter, r *http.Request) {

	b.Data["posts"] = b.Context.GetPosts()
	var tmpl = template.Must(template.New("index").ParseFiles("template/posts.tmpl", "template/layout.tmpl", "template/header.tmpl", "template/footer.tmpl"))
	tmpl.ExecuteTemplate(w, "layout", &b.Data)
}

func (b *BlogController) newPostHandler(w http.ResponseWriter, r *http.Request) {

	var tmpl = template.Must(template.New("index").ParseFiles("template/new.tmpl", "template/layout.tmpl", "template/header.tmpl", "template/footer.tmpl"))
	tmpl.ExecuteTemplate(w, "layout", &b.Data)
}

func (b *BlogController) saveHandler(w http.ResponseWriter, r *http.Request) {
	post := blog.Post{}

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	post.Title = r.FormValue("title")
	post.Body = []byte(r.FormValue("body"))
	status, err := strconv.ParseInt(r.FormValue("status"), 0, 32)

	if err != nil {
		log.Fatal(err)
	}

	post.Status = int(status)
	post.Posted = time.Now().Format(time.Stamp)
	post.Updated = time.Now().Format(time.Stamp)

	context := blog.BlogContext{}
	context.Connect()
	context.SavePost(post)

	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {

	blogController := NewBlogController()
	err := blogController.LoadConfig("config.json")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", blogController.mainHandler)
	http.HandleFunc("/new", blogController.newPostHandler)
	http.HandleFunc("/save", blogController.saveHandler)
	http.ListenAndServe(":80", nil)
}
