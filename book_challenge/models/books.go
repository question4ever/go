package models

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
Book object
	- title
	- author
	- page count
	- thumbnail
*/
type Book struct {
	Title     string
	Author    string
	PageCount int
	Thumbnail string
	PagesRead int
}

/*
CreateBook creates a book ojbect and returns a pointer of the object
*/
func CreateBook(title string, author string, pageCount int, thumbnail string) *Book {
	b := Book{title, author, pageCount, thumbnail, 0}
	s := fmt.Sprintf("%s*%s*%d*%d*%s\n", b.Title, b.Author, b.PageCount, b.PagesRead, b.Thumbnail)
	f, err := os.OpenFile("books.sav", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(s); err != nil {
		log.Fatal(err)
	}
	f.Close()
	return &b
}

/*GetBooks returns a list of books from books.sav*/
func GetBooks() []Book {
	f, err := os.OpenFile("books.sav", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var books []Book
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		b := strings.Split(line, "*")
		if err == io.EOF {
			break
		}
		pc, err := strconv.Atoi(b[2])
		if err != nil {
			panic(err)
		}
		pr, err := strconv.Atoi(b[3])
		if err != nil {
			panic(err)
		}
		book := Book{b[0], b[1], pc, b[4], pr}
		books = append(books, book)
		fmt.Println(b[0])
	}
	f.Close()
	return books
}

func DeleteBook(index int) {
	books := GetBooks()
	e := os.Remove("books.sav")
	if e != nil {
		log.Fatal(e)
	}
	f, err := os.OpenFile("books.sav", os.O_CREATE|os.O_WRONLY, 0644)
	var s string
	if err != nil {
		log.Fatal(err)
	}
	for j, book := range books {
		if index != j {
			s = fmt.Sprintf("%s*%s*%d*%d*%s", book.Title, book.Author, book.PageCount, book.PagesRead, book.Thumbnail)
			f.WriteString(s)
		}
	}
	f.Close()
}

func SaveBook(index int, b Book) {
	books := GetBooks()
	books[index] = b
	fmt.Println(b.PagesRead)
	f, err := os.OpenFile("books.sav", os.O_WRONLY, 0644)
	var s string
	if err != nil {
		log.Fatal(err)
	}
	for _, book := range books {
		s = fmt.Sprintf("%s*%s*%d*%d*%s", book.Title, book.Author, book.PageCount, book.PagesRead, book.Thumbnail)
		f.WriteString(s)
	}
	f.Close()
}
