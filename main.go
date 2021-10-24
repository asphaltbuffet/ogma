// Application which greets you.
package main

import (
	"fmt"

	"github.com/asdine/storm/v3"
)

// A Listing contains relevant information for LEX listings.
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

	defer func() {
		err = db.Close()
	}() // use function closure to allow checking error from deferred db.Close
	if err != nil {
		fmt.Println("Failed to close db: ", err)
	}

	listing := Listing{
		ID:                  1,
		IssueNumber:         56, //nolint:gomnd // preliminary dev magic number use
		PageNumber:          1,
		IndexedMemberNumber: 2989, //nolint:gomnd // preliminary dev magic number use
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
