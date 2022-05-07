package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Node struct {
	x int
	y int
	h int // Costo estimado desde un nodo n hasta el nodo final u objetivo.
	g int // Costo exacto de la ruta desde el nodo inicial a cualquier nodo n.
	f int // Costo mínimo del nodo vecino.
	p [][]int // Arreglos de nodos recorridos hasta este nodo.
}

type Point struct {
	x int
	y int
}

var rowNumber int
var columnNumber int

var openNodes []Node

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
	// se ordena asc de acuerdo al f.
	sort.Slice(openNodes[:], func(i, j int) bool {
		if openNodes[i].f != openNodes[j].f {
			return openNodes[i].f < openNodes[j].f
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
	t := tx + ty // Costo estimado desde un nodo A hasta el nodo B.
	return t
}

func findTheWay(xf, yf int) bool {
	currentNode := openNodes[0]
	g := currentNode.g + 1 // el costo para moverse es uno.
	x := currentNode.x
	y := currentNode.y
	p := currentNode.p
	p = append(p, []int{x, y})
	maze[x][y] = 2

	if currentNode.x == xf && currentNode.y == yf {
		return true
	}

	openNodes = openNodes[1:]
	neighborDirections := getNeighborDirections(currentNode)

	for _, point := range neighborDirections {
		h := calculateHeuristic(point.x, point.y, xf, yf)
		f := h + g
		addOrUpdateNode(point, h, g, f, p)
		maze[point.x][point.y] = 2
	}

	sortNodes() // se ordena asc de acuerdo al f.
	if len(openNodes) == 0 {
		return true // no se encontró solución.
	}
	return false
}

func addOrUpdateNode(point Point, h int, g int, f int, p [][]int) {
	for _, node := range openNodes {
		// Se actualiza los valores de los nodos abiertos respecto al nuevo punto
		if node.x == point.x && node.y == point.y {
			if f < node.f {
				node.h = h
				node.g = g
				node.f = f
				node.p = p
				return
			}
		}
	}

	var newNode = Node{point.x, point.y, h, g, f, p}
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
		*neighborDirections = append(*neighborDirections, Point{x, y})
	}
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

func deepCopyMap(currentMap [][]int) [][]int {
	newMap := make([][]int, len(currentMap))
	for i := range currentMap {
		newMap[i] = make([]int, len(currentMap[i]))
		copy(newMap[i], currentMap[i])
	}
	return newMap
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
	startTime := time.Now()
	rowNumber = len(maze)
	columnNumber = len(maze[0])

	x, y := getRandomPointInMaze()
	maze[x][y] = 2
	xf, yf := getRandomFinishPoint(x, y)
	maze[xf][yf] = 4
	addObstacles(x, y, xf, yf, 30)
	tracedPathMap := deepCopyMap(maze)

	fmt.Println("Laberinto")
	fmt.Println("x: ", x, " y: ", y)
	fmt.Println("fx: ", xf, " fy: ", yf)
	printMaze(maze)

	node := Node{x, y, 0, 0, 0, [][]int{}}
	openNodes = append(openNodes, node)

	var found bool
	for {
		if findTheWay(xf, yf) {
			found = true
			break
		}
	}

	fmt.Println("Camino encontrado")
	if found {
		currentNode := openNodes[0]
		for _, pathNode := range currentNode.p {
			tracedPathMap[pathNode[0]][pathNode[1]] = 7
		}
	}
	printMaze(tracedPathMap)
	elapsed := time.Since(startTime)
	fmt.Printf("Tomó %s\n", elapsed)
}