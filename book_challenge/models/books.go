package models

import (
	"fmt"
	"log"
	"os"
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
	s := fmt.Sprintf("Title: %s Author: %s Page Count: %d Thumbnail: %s\n", b.Title, b.Author, b.PageCount, b.Thumbnail)
	f, err := os.OpenFile("books.sav", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.WriteString(s); err != nil {
		log.Fatal(err)
	}
	return &b
}
