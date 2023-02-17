package lib

import "fmt"

func HandleErr(err error) {
	if err != nil {
		fmt.Printf("\n\n")
		panic(err)
	}
}
