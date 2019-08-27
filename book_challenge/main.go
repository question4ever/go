package main

import (
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

/*
AddBook
	- takes input from a form and calls the create function from the book model
*/
func AddBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "static\\newBook.html")
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
	r.HandleFunc("/library/add", AddBook)   //Needs google api integration to get thumbnail

	log.Fatal(http.ListenAndServe(":8000", r))
}
