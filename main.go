package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type cell struct {
	thisX       int
	thisY       int
	parentX     int
	parentY     int
	parent      *cell
	parentsCost float64
	fullCost    float64
	roadCost    float64
}

func toSquare(a int) int {
	return a * a
}

func (c *cell) calculateTotalCost(aimX, aimY int) {
	costToAim := math.Sqrt(float64(toSquare(aimX-c.thisX) + toSquare(aimY-c.thisY)))
	c.fullCost = c.roadCost + costToAim
}

func (c *cell) roadCostCalculate() {
	costToNeighbor := math.Sqrt(float64(toSquare(c.parentX-c.thisX) + toSquare(c.parentY-c.thisY)))
	c.roadCost = costToNeighbor + c.parentsCost
}

func main() {
	inpMap := readPathFromFile("input.txt")
	showPath(inpMap)
	pathFinder(inpMap, 1, 1, 3, 6)
}

func initCell(x, y int) cell {
	var c cell
	c.thisX = x
	c.thisY = y
	c.parentX = -1
	c.parentY = -1
	c.parentsCost = 0
	c.fullCost = 0
	return c
}

// x1,y1 — start of path
// x2,y2 — end of path

func pathFinder(inpMap [][]string, startX, startY, aimX, aimY int) ([]cell, []cell, string) {
	c := initCell(startX, startY)
	openList := make([]cell, 0, 0)
	closedList := make([]cell, 0, 0)
	openList = append(openList, c)
	for {
		selected := getCheapestCell(openList)
		if checkIfSuccess(selected.thisX, selected.thisY, aimX, aimY) {
			return openList, closedList, "FAILURE"
		}
		deleteCellFromSlice(openList, findIndexOfCell(selected, openList))
		closedList = append(closedList, selected)
		for _, v := range getNotClosedNeighbors(&selected, closedList, openList, inpMap, aimX, aimY) {
			temp := math.Sqrt(float64(toSquare(v.thisX-c.thisX) + toSquare(v.thisY-c.thisY)))
			if findIndexOfCell(c, openList) == -1 || temp < v.roadCost {
				v.parent = &c
				v.roadCost = temp
				v.calculateTotalCost(aimX, aimY)
				if findIndexOfCell(c, openList) == -1 {
					openList = append(openList, v)
				}
			}
		}
	}
	//return openList,closedList, "HORRAY"
}

func findIndexOfCell(c cell, cArr []cell) int {
	for i := range cArr {
		if c == cArr[i] {
			return i
		}
	}
	return -1
}

func deleteCellFromSlice(c []cell, i int) []cell {
	temp := make([]cell, len(c)-1, len(c)-1)
	copy(temp[:i], c[:i])
	copy(temp[i:], c[i+1:])
	c = temp
	return c
}

func getNotClosedNeighbors(c *cell, closedList, openList []cell, initMap [][]string, aimX, aimY int) []cell {
	cellArr := make([]cell, 0, 0)
	for i := c.thisX - 1; i < c.thisX+2; i++ {
		for j := c.thisY - 1; j < c.thisY+2; j++ {
			if i == 1 && j == 1 {
				continue
			}
			if initMap[i][j] != "X" {
				var b = true
				for item := range closedList {
					if closedList[item].thisX == i && closedList[item].thisY == j {
						b = false
						break
					}
				}
				for item := range openList {
					if openList[item].thisX == i && openList[item].thisY == j {
						b = false
						cellArr = append(cellArr, openList[item])
						break
					}
				}
				if b {
					cellArr = append(cellArr, cell{
						thisX:       i,
						thisY:       j,
						parentX:     c.thisX,
						parentY:     c.thisY,
						parent:      c,
						parentsCost: c.fullCost,
					})
					cellArr[len(cellArr)-1].roadCostCalculate()
					cellArr[len(cellArr)-1].calculateTotalCost(aimX, aimY)
				}
			}
		}
	}
	return cellArr
}

func checkIfSuccess(x1, y1, x2, y2 int) bool {
	return (math.Abs(float64(x1-x2)) <= 1) && (math.Abs(float64(y1-y2)) <= 1)
}

func getCheapestCell(openList []cell) cell {
	if len(openList) == 0 {
		panic("Empty open list!")
	}
	var minCell = openList[0]
	for i := range openList {
		if minCell.fullCost > openList[i].fullCost {
			minCell = openList[i]
		}
	}
	return minCell
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
