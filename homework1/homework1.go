/*
Author: Mattia Scantamburlo.

Goal: Write a program in Go that counts the number of times a certain
character "x" appears in a string. The program must use competition,
starting a goroutine for each character in the string and checking if the character
corresponds to the character searched.

NOTE: The below code does not distinguish between the upper case and lower case
character.
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {

	var wg = sync.WaitGroup{}
	// This variable is the string containing all possible character to confront with.
	var str string = "aaaaaaaaaaaaabbbbbbbbcccccddddccccccfff"
	// Character to confront each element of the string with.
	var chr rune = 'c'
	// This variable is used to get the answer of the qeuestions from the user.
	var answer string

	// User interaction.
	// Asking wheter the user want to use the default string or not.
	fmt.Printf("Do you want to use the default string: %s (y/n)\n", str)
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	// Making the string lower case.
	answer = strings.ToLower(answer)
	// Removing the new line character.
	answer = strings.TrimRight(answer, "\r\n")
	if err != nil {
		log.Fatal(err)
	}
	for (answer != "y") && (answer != "n") {
		fmt.Println("Please enter a valid answer (y/n)")
		answer, err = reader.ReadString('\n')
		// Removing the new line character.
		answer = strings.TrimRight(answer, "\r\n")
		// Making the string lower case.
		answer = strings.ToLower(answer)
		// Checking the presence of errors.
		if err != nil {
			log.Fatal(err)
		}
	}

	// If not then the user is asked to enter the string.
	if answer == "n" {
		fmt.Print("Enter the string to search in:\nString: ")
		str, err = reader.ReadString('\n')
		// Removing the new line character.
		str = strings.TrimRight(str, "\r\n")
		// Making the string lower case.
		str = strings.ToLower(str)
		// Checking the presence of errors.
		if err != nil {
			log.Fatal(err)
		}
		if len(str) == 0 {
			str = " "
		}
	}

	// Asking wheter the user want to use the default character or not.
	fmt.Printf("Do you want to use the default character for the research (y/n): %c \n", chr)
	answer, err = reader.ReadString('\n')
	// Making the string lower case.
	answer = strings.ToLower(answer)
	// Removing the new line character.
	answer = strings.TrimRight(answer, "\r\n")
	// Checking the presence of errors.
	if err != nil {
		log.Fatal(err)
	}
	for (answer != "y") && (answer != "n") {
		fmt.Println("Please enter a valid answer (y/n)")
		answer, err = reader.ReadString('\n')
		// Removing the new line character.
		answer = strings.TrimRight(answer, "\r\n")
		// Making the string lower case.
		answer = strings.ToLower(answer)
		// Checking the presence of errors.
		if err != nil {
			log.Fatal(err)
		}
	}

	// If not then the user is asked to enter the character.
	if answer == "n" {
		// Input from terminal (only first element, if present will be read)
		var line string
		// User interaction.
		fmt.Print("Enter the character to search in the given string ",
			"(if too long considered first element):\nCharacter: ")
		line, err = reader.ReadString('\n')
		// Removing the new line character.
		line = strings.TrimRight(line, "\r\n")
		// Making the string lower case.
		line = strings.ToLower(line)
		// Checking the presence of errors.
		if err != nil {
			log.Fatal(err)
		}
		if len(line) == 0 {
			chr = ' '
		} else {
			chr = []rune(line)[0]
		}
	}

	// Converting the string into an array of runes.
	var runeConv []rune = []rune(str)
	ch := make(chan int, len(str))

	fmt.Println("\n-->Calculating the number of iteration of the character: \"", string(chr),
		"\" in the given string.")

	wg.Add(1)
	// Starting the sum routine
	go sum(&wg, ch, len(str))

	for i := 0; i < len(str); i++ {
		wg.Add(1)
		// Starting the counters routine to give the sum routine the correct values for each character.
		go counter(&wg, ch, runeConv[i], chr)
	}

	wg.Wait()
	close(ch)
}

/*
Sum printer: this function sums the value associated by the function counter to each rune
in the string.
*/
func sum(wg *sync.WaitGroup, ch <-chan int, max int) {
	sum := 0
	for i := 0; i < max; i++ {
		sum += <-ch
	}
	fmt.Print("\n-->Result: there are/is ", sum, " iteration/s of the selected character in the considered string.\n\n")
	wg.Done()
}

/*
Counter: this function send a number in the given integer channel, representing the value of the
comparison between the two given rune. Return 1 if equals, 0 if not.
*/
func counter(wg *sync.WaitGroup, wr chan<- int, comp rune, chr rune) {
	if comp == chr {
		wr <- 1
	} else {
		wr <- 0
	}
	wg.Done()
}
