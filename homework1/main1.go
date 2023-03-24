/*
author: Mattia Scantamburlo
NOTE:...
*/

package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {

	var str string = "aaaaaaaaaaaaabbbbbbbbcccccddddccccccfff"
	var chr rune = 'b'
	var runeConv []rune = []rune(str)
	ch := make(chan int, len(str))

	wg.Add(1)
	go sum(&wg, ch, len(str))

	for i := 0; i < len(str); i++ {
		wg.Add(1)
		go counter(&wg, ch, runeConv[i], chr)
	}

	wg.Wait()
	close(ch)
}

// sum printer: it only read a number of value equal to the length of the string
func sum(wg *sync.WaitGroup, ch <-chan int, max int) {
	sum := 0
	for i := 0; i < max; i++ {
		sum += <-ch
	}
	fmt.Println("There are ", sum, " iteration of the selected character")
	wg.Done()
}

// counter: it always send a number representing wether or not current
// character is equal to the searched one
func counter(wg *sync.WaitGroup, wr chan<- int, comp rune, chr rune) {
	if comp == chr {
		wr <- 1
	} else {
		wr <- 0
	}
	wg.Done()
}
