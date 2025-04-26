package main

import (
	"fmt"
	"strings"
)

func main() {
	filter := " fdssf, fdsfd, dgfdgfdg     "
	froms := strings.FieldsFunc(filter, func(r rune) bool { return r == ',' || r == ' ' })
	for _, v := range froms {
		fmt.Println(v)
	}
}
