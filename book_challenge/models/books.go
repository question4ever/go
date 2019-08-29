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
}

/*
CreateBook creates a book ojbect and returns a pointer of the object
*/
func CreateBook(title string, author string, pageCount int, thumbnail string) *Book {
	b := Book{title, author, pageCount, thumbnail}
	s := fmt.Sprintf("%s,%s,%d,%s\n", b.Title, b.Author, b.PageCount, b.Thumbnail)
	f, err := os.OpenFile("books.sav", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(s); err != nil {
		log.Fatal(err)
	}
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
		b := strings.Split(line, ",")
		if err == io.EOF {
			break
		}
		pc, err := strconv.Atoi(b[2])
		if err != nil {
			panic(err)
		}
		book := Book{b[0], b[1], pc, b[3]}
		books = append(books, book)
	}
	return books
}
