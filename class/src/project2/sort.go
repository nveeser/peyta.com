package main

import (
	"bufio"
	"log"
	"os"
	"qsort"
	"strconv"
)

func numbers(name string) ([]uint64, error) {
	log.Printf("Opening File: %s", name)
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var numbers []uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, u)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return numbers, nil
}

func main() {

	//input, err := numbers("src/project2/QuickSort.txt")
	input, err := numbers("src/project2/1000.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found numbers: %d", len(input))

	count := qsort.Sort(input, qsort.FirstElement)
	log.Printf("Sorted First: %d", count)

	count = qsort.Sort(input, qsort.LastElement)
	log.Printf("Sorted Last: %d", count)

	count = qsort.Sort(input, qsort.Median3)
	log.Printf("Sorted Median3: %d", count)
}
