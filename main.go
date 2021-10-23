// Application which greets you.
package main

import "fmt"
import "github.com/asdine/storm/v3"

// Contains relevant information for LEX listings.
type Listing struct {
	ID                  int `storm:"id,increment"`
	IssueNumber         int
	PageNumber          int
	IndexedMemberNumber int    `storm:"index"`
	IndexedCategory     string `storm:"index"`
	ListingText         string
}

func main() {
	fmt.Println(greet())

	db, err := storm.Open("ogma.db")
	if err != nil {
		fmt.Println("Failed to open db: ", err)
	}

	defer db.Close()

	listing := Listing{
		ID:                  1,
		IssueNumber:         56,
		PageNumber:          1,
		IndexedMemberNumber: 2989,
		IndexedCategory:     "Art & Photography",
		ListingText:         "Fingerpainting exchange.",
	}

	err = db.Save(&listing)
	if err != nil {
		fmt.Println("Failed to save to db: ", err)
	}

}

func greet() string {
	return "Hi!"
}
