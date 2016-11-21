package main

import (
	"bufio"
	"compress/gzip"
	"hashing"
	"io"
	"log"
	"myheap"
	"os"
	"strconv"
	"strings"
)

func lines(filename string, f func(s string) error) error {
	var r io.Reader
	log.Printf("Opening File: %s", filename)
	{
		fr, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fr.Close()
		r = fr
	}

	if strings.HasSuffix(filename, ".gz") {
		zr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer zr.Close()
		r = zr
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := f(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func readnums(file string) []int64 {
	var nums []int64
	err := lines(file, func(s string) error {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		nums = append(nums, n)
		return nil
	})
	if err != nil {
		log.Fatalf("error reading lines: %s", err)
	}
	log.Printf("found numbers: %d", len(nums))
	return nums
}

func do2Sum(file string) {
	nums := readnums(file)

	buckets, tableSize := hashing.CountBuckets(nums)

	log.Printf("Buckets: %d", buckets)
	c := float32(len(nums)) / float32(buckets)
	log.Printf("Collisions: %f", c)
	a := float32(buckets) / float32(tableSize)
	log.Printf("Alpha: %f", a)

	distinct := hashing.SpecialSums(nums, -10000, 10000)
	log.Printf("Distinct: %d", distinct)
}

func doMedian(file string) {

}

func main() {
	{
		nums := readnums("src/project6/2sum.txt")
		distinct := hashing.SpecialSums(nums, -10000, 10000)
		log.Printf("Distinct: %d", distinct)
	}
	{
		nums := readnums("src/project6/Median.txt")

		var total int64
		m := &myheap.Median{}
		for _, n := range nums {
			m.Add(n)
			//log.Printf("Added %d Median: %d", n, m.Value())
			total += m.Value()
		}
		log.Printf("Total: %d", total%10000)
	}
}
