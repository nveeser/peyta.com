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
	"time"
)

type onLine func(s string) error

func lines(filename string, ol onLine) error {
	log.Printf("Opening File: %s", filename)
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer zr.Close()

	scanner := bufio.NewScanner(zr)
	for scanner.Scan() {
		if err := ol(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
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
	err := lines("src/project4/SCC.txt.gz", func(s string) error {
		e, err := parse(s)
		if err == nil {
			edges = append(edges, e)
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found edges: %d", len(edges))

	start := time.Now()
	groups := graph.Kosaraju(edges)
	elapsed := time.Since(start)

	log.Printf("found groups: %d (Time: %s)", len(groups), elapsed)

	log.Printf("Result: %v", graph.LargestGroups(groups, 5))
}
