package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"welsh_powell_and_a_star/src/astar"
	"welsh_powell_and_a_star/src/welshpowell"
)

func parseGraphInput(graphJSONBytes []byte) (astar.GraphInput, error) {
	graphInput := astar.GraphInput{}

	err := json.Unmarshal(graphJSONBytes, &graphInput)
	if err != nil {
		return graphInput, errors.Wrap(err, "couldn't parse graph input")
	}

	return graphInput, nil
}

func main() {
	graphJSONBytes, err := ioutil.ReadFile("../graph.json")
	if err != nil {
		log.Panic(err)
	}

	graphInput, err := parseGraphInput(graphJSONBytes)
	if err != nil {
		log.Panic(err)
	}

	path, err := astar.AStar(&graphInput)
	if err != nil {
		log.Panic(err)
	}

	pathJSONbytes, _ := json.MarshalIndent(path, "", "  ")
	fmt.Println("A* -> ", string(pathJSONbytes))

	coloredVertices := welshpowell.WelshPowell(graphInput.Graph)

    colorMax := 0
    for _, color := range coloredVertices {
        if colorMax <= color{
            colorMax = color
        }

    }
    fmt.Println("Amount of colors used ->",colorMax + 1, "\n")

	coloredVerticesJSONbytes, _ := json.MarshalIndent(coloredVertices, "", "  ")
	fmt.Println("Welsh Powell Colors ->", string(coloredVerticesJSONbytes))
}
