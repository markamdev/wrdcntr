package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/markamdev/wrdcntr/counter"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "f", "", "path to input file")
	flag.Parse()

	// temp
	fileName = "/home/markamdev/Development/arista-assignment/files/test_1.txt"

	if len(fileName) == 0 {
		fmt.Println("input file is a mandatory param")
		flag.PrintDefaults()
		return
	}

	fmt.Println("word occurence counter: reading from", fileName)

	// open file and check if successfull
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error while opening file:", err)
	}
	defer file.Close()

	// create new word counter instance
	newCounter := counter.CreateCounter()

	// create text file scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(counter.SentenceSplitter)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			newCounter.AddSentence(line)
		}
	}

	statsPrinter(newCounter.GetStats())
}

func statsPrinter(stats []counter.StatElem) {
	fmt.Printf("%-4s\t%-20s\t%s\t%s\n", "id", "word", "count", "sentences")
	for idx, elem := range stats {
		fmt.Printf("%04d\t%-20s\t%d\t%v\n", idx, elem.Word, elem.Count, elem.Sentences)
	}
}
