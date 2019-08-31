package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"book_challenge/models"
	"net/http"

	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static\\index.html")
}

func UpdateIndex() string {

	books := models.GetBooks()

	var buffer bytes.Buffer
	buffer.WriteString("<table><tr><th>Title</th><th>Author</th><th>Page Count</th></tr>")

	for i, b := range books {
		s := fmt.Sprintf("<tr><td>%s</td><a><td>%s</td><td>%d</td><td><a href=\"http://localhost:8000/library/edit/%d\">edit</a></td></tr>", b.Title, b.Author, b.PageCount, i)
		buffer.WriteString(s)
	}
	buffer.WriteString("</table>")

	return buffer.String()
}

func Library(w http.ResponseWriter, r *http.Request) {
	out := UpdateIndex()
	fmt.Fprint(w, out)
}

func EditBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		index := mux.Vars(r)["id"]
		books := models.GetBooks()
		i, err := strconv.Atoi(index)
		if err != nil {
			log.Fatal(err)
		}
		if i < len(books) {
			b := books[i]
			fmt.Print(b.Title)
			fmt.Fprintf(w, "<p>%s</p><p>%s</p><p>%d</p><img src=\"%s\">", b.Title, b.Author, b.PageCount, b.Thumbnail)
		}
	}
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
			fmt.Fprintf(w, "<h2>%s addded to library!</h2>", b.Title)
		}
	}
}

func main() {
	//router.PathPrefix("/bc/").Handler(http.StripPrefix("/bc/", http.FileServer(http.Dir(dir))))
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/library/", Library) //Needs to parse books.sav and popluate a table/list of books
	//r.HandleFunc("/library/edit", EditBook) //Allow the library to be edited
	r.HandleFunc("/library/read", ReadBook) //Update pages read on a book
	r.HandleFunc("/library/add", AddBook)
	r.Path("/library/edit/{id}").HandlerFunc(EditBook)

	log.Fatal(http.ListenAndServe(":8000", r))
}
