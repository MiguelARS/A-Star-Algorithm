package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Node struct {
	x int
	y int
	h int
	w int
	g int
	p [][]int
}

type Point struct {
	x int
	y int
}

var rowNumber int
var columnNumber int

var openNodes []Node
var closedPoints []Point

var maze [][]int = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func sortNodes() {
	sort.Slice(openNodes[:], func(i, j int) bool {
		if openNodes[i].g != openNodes[j].g {
			return openNodes[i].g < openNodes[j].g
		} else {
			return openNodes[i].h < openNodes[j].h
		}
	})
}

func calculateHeuristic(x, y, xf, yf int) int {
	tx := x - xf
	ty := y - yf
	if tx < 0 {
		tx *= -1
	}
	if ty < 0 {
		ty *= -1
	}
	t := tx + ty
	return t
}

func findTheWay(xf, yf, index int, doChan chan struct{}, doChannelFound chan struct{}) {
	currentNode := openNodes[index]
	w := currentNode.w + 1
	x := currentNode.x
	y := currentNode.y
	p := currentNode.p
	// p = append(p, []int{x, y})
	maze[x][y] = 2

	closedPoints = append(closedPoints, Point{x, y})

	if currentNode.x == xf && currentNode.y == yf {
		doChannelFound <- struct{}{}
		return
	}

	neighborDirections := getNeighborDirections(currentNode)

	for _, point := range neighborDirections {
		h := calculateHeuristic(point.x, point.y, xf, yf)
		g := h + w
		addOrUpdateNode(point, h, w, g, p)
		maze[point.x][point.y] = 2
	}
	doChan <- struct{}{}
}

func addOpenNode(xf int, yf int, point Point, w int, p [][]int, wg *sync.WaitGroup) {
	h := calculateHeuristic(point.x, point.y, xf, yf)
	g := h + w
	addOrUpdateNode(point, h, w, g, p)
	maze[point.x][point.y] = 2
	wg.Done()
}

func addOrUpdateNode(point Point, h int, w int, g int, p [][]int) {
	newP := deepCopySlice(p)
	newP = append(newP, []int{point.x, point.y})

	for _, node := range openNodes {
		if node.x == point.x && node.y == point.y {
			if node.g < g {
				node.h = h
				node.w = w
				node.g = g
				node.p = newP
				return
			}
		}
	}

	var newNode = Node{point.x, point.y, h, w, g, newP}
	openNodes = append(openNodes, newNode)
}

func getNeighborDirections(node Node) []Point {
	neighborDirections := make([]Point, 0, 4)

	AddNeighborIfValid(&neighborDirections, node.x, node.y+1)
	AddNeighborIfValid(&neighborDirections, node.x, node.y-1)
	AddNeighborIfValid(&neighborDirections, node.x-1, node.y)
	AddNeighborIfValid(&neighborDirections, node.x+1, node.y)

	return neighborDirections
}

func AddNeighborIfValid(neighborDirections *[]Point, x int, y int) {
	if validValue(maze[x][y]) {
		newPoint := Point{x, y}
		if !checkIfPointIsClosed(newPoint) {
			*neighborDirections = append(*neighborDirections, newPoint)
		} else {
			fmt.Println("This didnt work")
		}
	}
}

func checkIfPointIsClosed(newPoint Point) bool {
	for _, point := range closedPoints {
		if point.x == newPoint.x && point.y == newPoint.y {
			return true
		}
	}
	return false
}

func validValue(value int) bool {
	return value == 0 || value == 4
}

func printMaze(currentMap [][]int) {
	for _, h := range currentMap {
		for _, value := range h {
			fmt.Print(value)
		}
		fmt.Println()
	}
}

func deepCopySlice(currentSlice [][]int) [][]int {
	newSlice := make([][]int, len(currentSlice))
	for i := range currentSlice {
		newSlice[i] = make([]int, len(currentSlice[i]))
		copy(newSlice[i], currentSlice[i])
	}
	return newSlice
}

func getRandomFinishPoint(x, y int) (xf, yf int) {
	xf, yf = getRandomPointInMaze()

	if calculateHeuristic(x, y, xf, yf) > 6 {
		return xf, yf
	} else {
		return getRandomFinishPoint(x, y)
	}
}

func getRandomPointInMaze() (x, y int) {
	x = rand.Intn(rowNumber-2) + 1
	y = rand.Intn(columnNumber-2) + 1
	return x, y
}

func addObstacles(x, y, xf, yf, numberOfObstacles int) {
	var currentNumberOfObstacles = 0
	for {
		if currentNumberOfObstacles >= numberOfObstacles {
			return
		}

		obstacleX, obstacleY := getRandomPointInMaze()
		if (obstacleX != x && obstacleY != y) && (obstacleX != xf && obstacleY != yf) && (maze[obstacleX][obstacleY] == 0) {
			maze[obstacleX][obstacleY] = 3
			currentNumberOfObstacles = currentNumberOfObstacles + 1
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	rowNumber = len(maze)
	columnNumber = len(maze[0])

	x, y := getRandomPointInMaze()
	maze[x][y] = 2
	xf, yf := getRandomFinishPoint(x, y)
	maze[xf][yf] = 4
	addObstacles(x, y, xf, yf, 30)
	tracedPathMap := deepCopySlice(maze)

	fmt.Println("Laberinto")
	fmt.Println("x: ", x, " y: ", y)
	fmt.Println("fx: ", xf, " fy: ", yf)
	printMaze(maze)

	node := Node{x, y, 0, 0, 0, [][]int{}}
	openNodes = append(openNodes, node)

	doChan := make(chan struct{})
	doFound := make(chan struct{})
	found := false
	for {
		start := len(openNodes)
		if start > 5 {
			start = 5
		}
		for i := 0; i < start; i++ {
			go findTheWay(xf, yf, i, doChan, doFound)
		}

		finish := 0
		for finish < start {
			select {
			case <-doFound:
				found = true
				finish++
				break
			case <-doChan:
				finish++
			}
		}
		if found {
			break
		}
		openNodes = openNodes[finish:]
		if len(openNodes) == 0 {
			break
		}

		sortNodes()
	}

	fmt.Println("Posiciones revisadas")
	printMaze(maze)

	fmt.Println("Camino encontrado")
	if found {
		fmt.Println("Dibujando el camino")
		currentNode := openNodes[0]
		for _, pathNode := range currentNode.p {
			tracedPathMap[pathNode[0]][pathNode[1]] = 7
		}
	} else {
		fmt.Println("No se encontro el camino")
	}
	printMaze(tracedPathMap)
	elapsed := time.Since(start)
	fmt.Printf("Tomo %s\n", elapsed)
}