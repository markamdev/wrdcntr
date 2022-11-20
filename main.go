package main

import (
	"fmt"

	"github.com/markamdev/wrdcntr/counter"
)

func main() {
	fmt.Println("word occurence counter")

	// simple test code
	newCounter := counter.CreateCounter()
	fmt.Println("Stats from empty counter:")
	statsPrinter(newCounter.GetStats())

	newCounter.AddSentence("This is a new sentence.")
	newCounter.AddSentence("This is another new sentence.")
	newCounter.AddSentence("This is some test string")

	fmt.Println("Stats from counter with 3 sentences:")
	statsPrinter(newCounter.GetStats())
}

func statsPrinter(stats []counter.StatElem) {
	fmt.Printf("%-4s\t%-20s\t%s\t%s\n", "id", "word", "count", "sentences")
	for idx, elem := range stats {
		fmt.Printf("%04d\t%-20s\t%d\t%v\n", idx, elem.Word, elem.Count, elem.Sentences)
	}
}
