package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var inpMap [][]string
	for scanner.Scan() {
		text := scanner.Text()
		inpMap = append(inpMap, strings.Split(text, ""))
	}
	for i := range inpMap {
		for j := range inpMap[i] {
			fmt.Printf("%2s", inpMap[i][j])
		}
		fmt.Println()
	}
}
