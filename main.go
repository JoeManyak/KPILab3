package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const e = 5

type cell struct {
	thisX       int
	thisY       int
	parentX     int
	parentY     int
	parent      int
	parentsCost float64
	fullCost    float64
	roadCost    float64
}

type queue struct {
	el []string
}

func (q *queue) addOne(str string) {
	temp := make([]string, len(q.el)+1, len(q.el)+1)
	temp[0] = str
	copy(temp[1:], q.el[0:])
	q.el = temp
}

func (q *queue) takeNext() string {
	str := q.el[len(q.el)-1]
	q.el = q.el[:len(q.el)-1]
	return str
}

func main() {
	inpMap := readPathFromFile("input.txt")
	showPath(inpMap)
	b := pathFinder(inpMap, 1, 1, 8, 1)
	if b {
		fmt.Println("Road built successfully!")
	} else {
		fmt.Println("No way to aim!")
	}
}

func queueGenerate() queue {
	var q queue
	for i := 49; i < 58; i++ {
		q.addOne(string(i))
	}
	for i := 97; i < 123; i++ {
		q.addOne(string(i))
	}
	for i := 65; i < 91; i++ {
		q.addOne(string(i))
	}
	return q
}

func toSquare(a int) int {
	return a * a
}

func (c *cell) calculateTotalCost(aimX, aimY int) {
	costToAim := math.Sqrt(float64(toSquare(aimX-c.thisX) + toSquare(aimY-c.thisY)))
	c.fullCost = c.roadCost + costToAim*e
}

func (c *cell) roadCostCalculate() {
	costToNeighbor := math.Sqrt(float64(toSquare(c.parentX-c.thisX) + toSquare(c.parentY-c.thisY)))
	c.roadCost = costToNeighbor + c.parentsCost
}

func returner(c cell, closedList []cell, inpMap [][]string, q queue) {
	inpMap[c.thisX][c.thisY] = q.takeNext()
	if c.parent != -1 {
		returner(closedList[c.parent], closedList, inpMap, q)
		return
	}
	showPath(inpMap)
	return
}

func initCell(x, y int) cell {
	var c cell
	c.thisX = x
	c.thisY = y
	c.parentX = -1
	c.parentY = -1
	c.parent = -1
	c.parentsCost = 0
	c.fullCost = 0
	return c
}

// x1,y1 — start of path
// x2,y2 — end of path

func pathFinder(inpMap [][]string, startX, startY, aimX, aimY int) bool {
	c := initCell(startX, startY)
	openList := make([]cell, 0, 0)
	closedList := make([]cell, 0, 0)
	openList = append(openList, c)
	for {
		selected, b := getCheapestCell(openList)
		if b {
			return false
		}
		if checkIfSuccess(selected.thisX, selected.thisY, aimX, aimY) {
			inpMap[aimX][aimY] = "0"
			q := queueGenerate()
			returner(selected, closedList, inpMap, q)
			return true
		}
		openList = deleteCellFromSlice(openList, findIndexOfCell(selected, openList))
		closedList = append(closedList, selected)
		for _, v := range getNotClosedNeighbors(closedList[len(closedList)-1],
			closedList, openList, inpMap, aimX, aimY) {
			temp := math.Sqrt(float64(toSquare(v.thisX-c.thisX) + toSquare(v.thisY-c.thisY)))
			if findIndexOfCell(v, openList) == -1 || temp < v.roadCost {
				v.parent = len(closedList) - 1
				v.roadCost = temp
				v.calculateTotalCost(aimX, aimY)
				if findIndexOfCell(v, openList) == -1 {
					openList = append(openList, v)
				}
			}
		}
	}
}

func findIndexOfCell(c cell, cArr []cell) int {
	for i := range cArr {
		if c.thisX == cArr[i].thisX && c.thisY == cArr[i].thisY {
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

func getNotClosedNeighbors(c cell, closedList, openList []cell, initMap [][]string, aimX, aimY int) []cell {
	cellArr := make([]cell, 0, 0)
	for i := c.thisX - 1; i < c.thisX+2 && i < len(initMap); i++ {
		for j := c.thisY - 1; j < c.thisY+2 && j < len(initMap[i]); j++ {
			if i == 1 && j == 1 {
				continue
			}
			if initMap[i][j] != "X" { //Англійська X
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
						parent:      len(closedList) - 1,
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

func getCheapestCell(openList []cell) (cell, bool) {
	if len(openList) == 0 {
		return cell{}, true
	}
	var minCell = openList[0]
	for i := range openList {
		if minCell.fullCost > openList[i].fullCost {
			minCell = openList[i]
		}
	}
	return minCell, false
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
	fmt.Println(">>")
	for i := range inpMap {
		for j := range inpMap[i] {
			fmt.Printf("%2s", inpMap[i][j])
		}
		fmt.Println()
	}
}
