package main

import (
	"bufio"
	"fmt"
	"graph"
	"log"
	"os"
	"strconv"
	"strings"
)

func lines(name string) ([]string, error) {
	log.Printf("Opening File: %s", name)
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var r []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r = append(r, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return r, nil
}

func parse(line string) (*graph.Edge, error) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("line does not seem formated correctly [%q] => %+v", line, parts)
	}

	left, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	right, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &graph.Edge{graph.ID(left), graph.ID(right)}, nil
}

func main() {
	data, err := lines("src/project4/SCC.txt")
	if err != nil {
		log.Fatal(err)
	}

	var edges []*graph.Edge
	for _, l := range data {
		e, err := parse(l)
		if err != nil {
			log.Fatal(err)
		}
		edges = append(edges, e)
	}

	log.Printf("found edges: %d", len(edges))

	groups := graph.Kosaraju(edges)

	log.Printf("found groups: %d", len(groups))

	log.Printf("Result: %v", graph.LargestGroups(groups, 5))
}
