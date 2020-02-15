package main

import (
	"bufio"
	"inversion"
	"log"
	"os"
	"strconv"
)

func numbers() ([]uint64, error) {
	file, err := os.Open("IntegerArray.txt")
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
	input, err := numbers()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found numbers: %d", len(input))

	output, count := inversion.Invert(input)
	log.Printf("output: %v", output)
	log.Printf("Inverted count: %d", count)
}
