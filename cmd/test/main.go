package main

import (
	"fmt"
	"log"

	"github.com/Muchprow/higuruma/internal/modrinth"
)

func main() {
	results, err := modrinth.Search("jei", 5, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found:", results.Total)
	for _, p := range results.Hits {
		fmt.Printf("- %s (%s) — %d downloads\n", p.Title, p.Slug, p.Downloads)
	}
}
