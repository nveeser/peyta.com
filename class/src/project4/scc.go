package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"graph"
	"log"
	"os"
	"strconv"
	"strings"
)

type onLine func(s string)

func lines(filename string, ol onLine) ([]string, error) {
	log.Printf("Opening File: %s", name)
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	var lines []string
	scanner := bufio.NewScanner(zr)
	for scanner.Scan() {
		ol(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
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
	var edges []*graph.Edge

	data, err := lines("src/project4/SCC.txt.gz", func(s string) {
		e, err := parse(s)
		if err != nil {
			log.Fatal(err)
		}
		edges = append(edges, e)

	})

	log.Printf("found edges: %d", len(edges))

	groups := graph.Kosaraju(edges)

	log.Printf("found groups: %d", len(groups))

	log.Printf("Result: %v", graph.LargestGroups(groups, 5))
}
