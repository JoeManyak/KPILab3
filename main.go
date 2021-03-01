package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type cell struct {
	thisX     int
	thisY     int
	parentX   int
	parentY   int
	totalCost float64
}

func toSquare(a int) int {
	return a * a
}

func (c *cell) calculateTotalCost(aimX, aimY int) {
	costToParent := math.Sqrt(float64(toSquare(c.parentX-c.thisX) + toSquare(c.parentY-c.thisY)))
	costToAim := math.Sqrt(float64(toSquare(aimX-c.thisX) + toSquare(aimY-c.thisY)))
	c.totalCost = costToParent + costToAim
}

func main() {
	inpMap := readPathFromFile("input.txt")
	showPath(inpMap)
	pathFinder(inpMap, 1, 1, 2, 2)
}

// x1,y1 — start of path
// x2,y2 — end of path

func pathFinder(inpMap [][]string, x1, y1, x2, y2 int) {
	if inpMap[x1][y1] == "X" || inpMap[x2][y2] == "X" {
		panic("Not Valid start")
	}
}

func readPathFromFile(filename string) [][]string {
	f, err := os.Open(filename)
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
	return inpMap
}

func showPath(inpMap [][]string) {
	for i := range inpMap {
		for j := range inpMap[i] {
			fmt.Printf("%2s", inpMap[i][j])
		}
		fmt.Println()
	}
}
