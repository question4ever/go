package main

import (
	"book_challenge/models"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"
	"github.com/gorilla/mux"
)

var p models.Points
var nav string

func Home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static\\index.html")
}

func UpdateIndex() string {

	books := models.GetBooks()

	var buffer bytes.Buffer
	buffer.WriteString("<table><tr><th>Title</th><th>Author</th><th>Page Count</th><th>Pages Read</th></tr>")

	for i, b := range books {
		s := fmt.Sprintf("<tr><td>%s</td><a><td>%s</td><td>%d</td><td>%d</td><td><a href=\"http://localhost:8000/library/read/%d\">Read</a></td><td><a href=\"http://localhost:8000/library/delete/%d\">Delete</a></td></tr>", b.Title, b.Author, b.PageCount, b.PagesRead, i, i)
		buffer.WriteString(s)
	}
	buffer.WriteString("</table><p><a href=\"http://localhost:8000/library/add\">Add Book</a></p>")
	return buffer.String()
}

func Library(w http.ResponseWriter, r *http.Request) {
	out := UpdateIndex()
	out = nav + out
	fmt.Fprint(w, out)
}

func Points(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nav+"<p>Points earned: %d</p><p><a href=\"http://localhost:8000/points/spend\">Spend</a></p>", p.Amount)
}

func SpendPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		str := fmt.Sprintf("<p id=\"points\">Points earned: %d</p><form action=\"/points/spend\" method=\"POST\">", p.Amount)
		str = nav + str
		if p.Amount >= 100 {
			str = str + fmt.Sprint("<label>Slurpee costs 100 points</label><input type=\"checkbox\" name=\"prize\" value=\"Slurpee\">")
		}
		if p.Amount >= 200 {
			str = str + fmt.Sprint("<label>Move Night costs 200 points</label><input type=\"checkbox\" name=\"prize\" value=\"Movie Night\">")
		}
		if p.Amount >= 1000 {
			str = str + fmt.Sprint("<label>Small Toy costs 1000 points</label><input type=\"checkbox\" name=\"prize\" value=\"Small Toy\">")
		}
		if p.Amount >= 5000 {
			str = str + fmt.Sprint("<label>Switch Game costs 5000 points</label><input type=\"checkbox\" name=\"prize\" value=\"Switch Game\">")
		}
		str = str + fmt.Sprint("<button type=\"submit\">Submit</button></form>")
		fmt.Fprint(w, str)
	} else if r.Method == "POST" {
		var str string
		err := r.ParseForm()
		if err != nil {
			fmt.Print(err)
		}
		str = nav
		prizesSelected := r.Form["prize"]
		for _, prize := range prizesSelected {
			p.SpendPoints(prize)
			str = str + fmt.Sprintf("<p>%s</p>", prize)
		}
		fmt.Fprint(w, str)
	}
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
			fmt.Fprintf(w, nav+"<p>%s</p><p>%s</p><p>%d</p><img src=\"%s\">", b.Title, b.Author, b.PageCount, b.Thumbnail)
		} else {
			fmt.Fprint(w)
		}
	}
}

func ReadBook(w http.ResponseWriter, r *http.Request) {
	index := mux.Vars(r)["id"]
	books := models.GetBooks()
	i, err := strconv.Atoi(index)
	if err != nil {
		log.Fatal(err)
	}
	if r.Method == "GET" {
		if i < len(books) {
			b := books[i]
			fmt.Fprintf(w,
				nav+"<p id=\"title\">Title: %s</p><p id=\"author\">Author: %s</p><p>Page Count: %d</p><p>Pages Read: %d</p><img id=\"thumbnail\" src=\"%s\"><form action=\"/library/read/%d\" method=\"POST\"><input type=\"number\" name=\"pagesRead\" id=\"pagesRead\"><button>Submit</button></form>", b.Title, b.Author, b.PageCount, b.PagesRead, b.Thumbnail, i)
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Print(err)
		}
		pr, err := strconv.Atoi(r.FormValue("pagesRead"))
		if err != nil {
			fmt.Print(err)
		}

		if (books[i].PagesRead + pr) >= books[i].PageCount {
			books[i].PagesRead = 0
			p.EarnPoints(10)
		} else {
			books[i].PagesRead = books[i].PagesRead + pr
		}
		p.EarnPoints(pr)

		models.SaveBook(i, books[i])
		fmt.Fprintf(w, nav+"<p>Nice job! You read %d pages of %s and earned %d total points! Keep up the super reading Tom!</p>", pr, books[i].Title, p.Amount)
	}
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
			fmt.Fprintf(w, nav+"<p>%s addded to library!</p>", b.Title)
		}
	}
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	index := mux.Vars(r)["id"]
	books := models.GetBooks()
	i, err := strconv.Atoi(index)
	if err != nil {
		log.Fatal(err)
	}
	if r.Method == "GET" {
		if i < len(books) {
			b := books[i]
			fmt.Fprintf(w,
				nav+"<p id=\"title\">Title: %s</p><p id=\"author\">Author: %s</p><p>Page Count: %d</p><p>Pages Read: %d</p><img id=\"thumbnail\" src=\"%s\"><form action=\"/library/delete/%d\" method=\"POST\"><button>Delete</button></form>", b.Title, b.Author, b.PageCount, b.PagesRead, b.Thumbnail, i)
		}
	} else if r.Method == "POST" {
		fmt.Fprintf(w, nav+"<p>%s successfully deleted</p>", books[i].Title)
		models.DeleteBook(i)
	}
}

func main() {
	//router.PathPrefix("/bc/").Handler(http.StripPrefix("/bc/", http.FileServer(http.Dir(dir))))
	p = models.Points{0, 0, "", 1, *new([]string)}
	nav = "<div class=\"topnav\"><a class=\"active\" href=\"http://localhost:8000/\">Home</a><a href=\"http://localhost:8000/library/\">Library</a><a href=\"http://localhost:8000/points\">Points</a></div>"
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/library/", Library) //Needs to parse books.sav and popluate a table/list of books
	//r.HandleFunc("/library/edit", EditBook) //Allow the library to be edited
	r.HandleFunc("/library/read", ReadBook) //Update pages read on a book
	r.HandleFunc("/library/add", AddBook)
	r.Path("/library/edit/{id}").HandlerFunc(EditBook)
	r.Path("/library/read/{id}").HandlerFunc(ReadBook)
	r.HandleFunc("/points", Points)
	r.HandleFunc("/points/spend", SpendPoints)
	r.Path("/library/delete/{id}").HandlerFunc(DeleteBook)

	port := os.Getenv("PORT")
	if port == "" {
  		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
