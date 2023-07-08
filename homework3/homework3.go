/*
Author: Mattia Scantamburlo.

Goal: Write a program in Go that simulates the production of 5 cakes by 3 pastry chefs.
The production of each cake requires 3 steps that must take place in order: first the cake is cooked,
then garnished and finally decorated.

NOTE:
1) I preferred making the program clearer in the development of routines by associating state variables
to each chef that are printed about every second (although this adds a computational cost).
2) The chefs can hold a cake in their hands after putting it into the next chef spaces, which is why for
example chef 1 is able to cook 4 cakes when the chef 2 has decorated only one. (he/she first fill the 2 spaces
available, then chef 2 start decorating one of these, so chef 1 can cake another cake to put into chef 2
freed space and then one more to hold in his hands).
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Starting timer.
	time_init := time.Now()
	fmt.Println("Starting timer...")

	// Defining the end channel and the comunication channel between the chefs (1-2, 2-3).
	chEnd := make(chan bool, 4)
	ch12 := make(chan int, 2)
	ch23 := make(chan int, 2)

	wg := sync.WaitGroup{}
	wg.Add(3) // Adding elements to WaitGroup.

	go chef1(&wg, ch12, chEnd, 5)
	go chef2(&wg, ch12, ch23, chEnd, 5)
	go chef3(&wg, ch23, chEnd, 5) //Notice: chef3 is also responsibnle for telling other cooks when to stop.

	wg.Wait()    // Waiting unitil the go routines get completed.
	close(chEnd) // Closing the end channel.
	fmt.Printf("\nIt took the 3 cooks: %f32 seconds", float32((time.Now().Sub(time_init)).Seconds()))
}

/*
COOKS THE CAKES (COOK1):
This function simulates the functionality of the chef number 1. The argument are:
-The WaitGroup pointer (wg).
-The channel to send chef number 2 signals (ch12).
-The channel to aknowledge the end of the execution (chEnd).
-The max number of cakes to cook (limit).
*/
func chef1(wg *sync.WaitGroup, ch12 chan int, chEnd <-chan bool, limit int) {
	// Check to see if the limit is a positive number.
	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}

	// Starting timer for chef number 1.
	time1 := time.Now()

	// STATE VARIABLE OF CHEF NUMBER 1:
	// Cake decorated by chef number 1.
	done1 := 0

	for {
		// Ending chef number 2 work when the number of cake garnished reach the limit.
		if done1 == limit {
			break
		}

		//COOKING CAKES PHASE
		if time.Now().Sub(time1) >= 1*time.Second {
			// Signal chef number 3 that a cake has been garnished (or wait until chef number 2 aknowledge it).
			ch12 <- 1
			time1 = time.Now() // "zerify" the time gap.
			done1++
			fmt.Printf("\n Chef number 1 has cooked %d cakes!\n", done1)
		}
		time.Sleep(50 * time.Millisecond) // Delay: to avoid doing to many check operation.
	}

	// Waiting to get an end signal.
	for {
		if <-chEnd {
			break
		}
		time.Sleep(50 * time.Millisecond) // Delay: to avoid checking necessaire resources constantly.
	}
	close(ch12) // Closing channel between chef number 1 and number 2.
	wg.Done()   // Decrease the WaitGroup count of one.
}

/*
GUARNISH CAKES (CHEF2):
This function simulates the functionality of the chef number 2. The argument are:
-The WaitGroup pointer (wg).
-The channel to receive chef number 1 signals (ch12).
-The channel to send chef number 3 signals (ch12).
-The channel to aknowledge the end of the execution (chEnd).
-The max number of cakes to garnish (limit).
*/
func chef2(wg *sync.WaitGroup, ch12 <-chan int, ch23 chan<- int, chEnd <-chan bool, limit int) {

	// Check to see if the limit is a positive number.
	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}

	// Starting timer for chef number 2.
	time2 := time.Now()

	// STATE VARIABLE OF CHEF NUMBER 2:
	// Cake decorated by chef number 2.
	done2 := 0
	// Cake available for chef number 2.
	available2 := 0
	// This variable tells whether the chef number 2 is cooking or not.
	cooking2 := false

	for {
		// Ending chef number 2 work when the number of cake garnished reach the limit.
		if done2 == limit {
			break
		}

		if (available2+done2 < limit) && (available2 < 2) {
			value := <-ch12 // Check to see if some cakes are ready to be garnished.
			if value > 0 {
				available2 += value
			}
			if available2 > 2 {
				panic("Chef number 2 only have 2 places!")
			}
		}

		// CHECKING GRANISHING CONDITION:
		if !cooking2 && available2 != 0 { // Checks if guarnishing cakes conditions are present.
			cooking2 = true
			time2 = time.Now() // "zerify" the time gap.
		}

		//GUARNISHING CAKES PHASE:
		if time.Now().Sub(time2) >= 4*time.Second && cooking2 {
			// Signal chef number 3 that a cake has been garnished (or wait until chef number 3 aknowledge it).
			ch23 <- 1
			done2++
			available2--
			cooking2 = false
			fmt.Printf("\n Chef number 2 has garnished %d cakes! \n", done2)
			//fmt.Printf(" Chef number 2 has %d available cakes!\n", available2)
		}

		time.Sleep(50 * time.Millisecond) // Delay: to avoid checking necessaire resources constantly.
	}

	// Waiting to get an end signal.
	for {
		if <-chEnd {
			break
		}
		time.Sleep(50 * time.Millisecond) // Delay: to avoid checking end channel constantly.
	}

	close(ch23) // Closing channel betwenn chef number 2 and number 3.
	wg.Done()   // Decrease WaitGroup count.
}

/*
DECORATE CAKES (CHEF 3):
This function simulates the functionality of the chef number 3. The argument are:
-The WaitGroup pointer (wg).
-The channel to receive chef number 2 signals (ch23).
-The channel to send the end of the execution signal (chEnd).
-The max number of cakes to decorate (limit).
*/
func chef3(wg *sync.WaitGroup, ch23 <-chan int, chEnd chan<- bool, limit int) {

	// Check to see if the limit is a positive number.
	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}

	// Starting timer for chef number 3.
	time3 := time.Now()

	// STATE VARIABLE OF CHEF NUMBER 3:
	// Cake decorated by chef number 3.
	done3 := 0
	// Cake available for chef number 3.
	available3 := 0
	// This variable tells whether the chef number 3 is cooking or not.
	cooking3 := false

	for {
		// Ending chef number 3 work when the number of cake garnished reach the limit.
		if done3 == limit {
			break
		}

		// CHECKING CAKES AVAILABILITY PHASE
		if (available3+done3 < limit) && (available3 < 2) { // Checks to see if cakes are necessaire.
			value := <-ch23 // Check to see if some cakes are ready to be decorated.
			if value > 0 {
				available3 += value
			}
			if available3 > 2 {
				panic("Chef number 3 only have 2 places!")
			}
		}

		// CHECKING DECORATING CONDITION:
		if !cooking3 && available3 != 0 { // Checks if decorating cakes conditions are present.
			cooking3 = true
			time3 = time.Now() // "zerify" the time gap.
		}

		//DECORATING CAKES PHASE:
		if time.Now().Sub(time3) >= 8*time.Second && cooking3 {
			available3--
			done3++
			cooking3 = false
			fmt.Printf("\n Chef number 3 has decorated %d cakes!\n", done3)
			//fmt.Printf(" Chef number 3 has %d available cakes!\n", available3)
		}

		time.Sleep(50 * time.Millisecond) // Delay: to avoid checking necessaire resources constantly.
	}

	// Signaling first and second chef that the cakes are completely done.
	chEnd <- true
	chEnd <- true
	wg.Done() // Decrease WaitGroup count.
}
