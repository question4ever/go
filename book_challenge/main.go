package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"book_challenge/models"
	"net/http"

	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static\\index.html")
}

func UpdateIndex() {

	books := models.GetBooks()

	path := "static\\library\\index.html"
	var buffer bytes.Buffer
	buffer.WriteString("<table><tr><th>Title</th><th>Author</th><th>Page Count</th></tr>")

	for _, b := range books {
		s := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td></tr>", b.Title, b.Author, b.PageCount)
		buffer.WriteString(s)
	}
	buffer.WriteString("</table>")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	info, _ := os.Stat(path)
	mode := info.Mode()

	array := strings.Split(string(file), "\n")
	fmt.Print(array)
	for i, line := range array {
		if line == "" {
			array[i] = buffer.String()
		}
	}

	ioutil.WriteFile(path, []byte(strings.Join(array, "\n")), mode)
}

func Library(w http.ResponseWriter, r *http.Request) {
	UpdateIndex()
	http.ServeFile(w, r, "static\\library\\index.html")
}

func EditBook(w http.ResponseWriter, r *http.Request) {

}

func ReadBook(w http.ResponseWriter, r *http.Request) {
}

/*
AddBook
	- takes input from a form and calls the create function from the book model
*/
func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static\\library\\newBook.html")
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Print(err)
		}
		i, err := strconv.Atoi(r.FormValue("pageCount"))
		if err == nil {
			var b = models.CreateBook(r.FormValue("title"), r.FormValue("author"), i, r.FormValue("thumbnail"))
			fmt.Printf("Title: %s\nAuthor: %s\nPage Count: %d\n", b.Title, b.Author, b.PageCount)
		}
	}
}

func main() {
	//router.PathPrefix("/bc/").Handler(http.StripPrefix("/bc/", http.FileServer(http.Dir(dir))))
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/library/", Library)      //Needs to parse books.sav and popluate a table/list of books
	r.HandleFunc("/library/edit", EditBook) //Allow the library to be edited
	r.HandleFunc("/library/read", ReadBook) //Update pages read on a book
	r.HandleFunc("/library/add", AddBook)

	log.Fatal(http.ListenAndServe(":8000", r))
}
