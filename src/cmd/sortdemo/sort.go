package main

import (
	"fmt"
	"sort"
)

//create a simple sort
func main() {
	//creates a slicen
	df := []int{3, 4, 2, 75, 8, 1}
	sort.Ints(df)
	for i, v := range df {
		fmt.Println(i, v)
	}
}
