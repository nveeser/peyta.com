package main

import (
	"bufio"
	"graph"
	"log"
	"os"
	"strconv"
	"strings"
)

func read(name string) ([]graph.Row, error) {
	log.Printf("Opening File: %s", name)
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []graph.Row

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids, err := parse(scanner.Text())
		if err != nil {
			return nil, err
		}
		result = append(result, ids)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func parse(line string) (graph.Row, error) {
	parts := strings.Split(line, "\t")

	var ids []uint64
	for _, s := range parts {
		if s == "" {
			continue
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, n)
	}
	return graph.Row(ids), nil
}

func main() {
	input, err := read("src/project3/Graph.txt")
	if err != nil {
		log.Fatal(err)
	}
	edges := graph.MakeEdges(input)
	mincut := graph.MinCut(len(input), edges)

	log.Printf("Final: %d", len(mincut))
	log.Printf("Graph: %d", mincut)
}
