package main

import (
	"bufio"
	"graph"
	"log"
	"os"
	"strings"
)

type onLine func(s string) error

func lines(filename string, ol onLine) error {
	log.Printf("Opening File: %s", filename)
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	// zr, err := gzip.NewReader(r)
	// if err != nil {
	// 	return err
	// }
	// defer zr.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := ol(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func main() {
	var edges []*graph.SizeEdge
	err := lines("src/project5/paths.txt", func(s string) error {
		e, err := graph.ParseNodePaths(s, "\t")
		if err == nil {
			edges = append(edges, e...)
		}
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found edges: %d", len(edges))

	log.Printf("First: %+v", edges[0])

	got := graph.FindDistances(edges, 1)

	ids := []graph.ID{7, 37, 59, 82, 99, 115, 133, 165, 188, 197}

	var r []string
	for _, id := range ids {
		log.Printf("D[%d]: %d", id, got[id])
		r = append(r, got[id].String())
	}
	log.Printf("Result: %s", strings.Join(r, ","))
}
