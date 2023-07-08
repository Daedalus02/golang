/*
Author: Mattia Scantamburlo.

Goal: Write a program in Go that simulates a couple trading activity in a fictitious market.
The program must simulate using the competition three couple pairs: EUR/USD, GBP/USD and JPY/USD,
and simulate buying and selling transactions in parallel.

NOTE:
1)   I decided to have the sales printed in the main function body to make it possible to
	 coordinate the combined printing of the trend values of the couples with those at the time of sale.
2)   I kept the channel buffers wide enough to avoid signal loss.
3)   I considered the moment for buying the couples only when the value go over or under the given limit
	(example jpy/usd: 0.0090, 0.0084 (buy), 0.0071 (do not buy)).
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// This constant is the maximum value the EUR/USD market value can reach before selling.
const eur_usd_MAX_Sale float32 = 1.20

// This constant is the minimum value the GPB/USD market value can reach before selling.
const gbp_usd_MIN_Purchase float32 = 1.35

// This constant is the minimum value the JPY/USD market value can reach before selling.
const jpy_usd_MIN_Purchase float32 = 0.0085

// This constant is the maximum value the EUR/USD couple can reach.
const eur_usd_MAX float32 = 1.5

// This constant is the maximum value the GPB/USD couple can reach.
const gbp_usd_MAX float32 = 1.5

// This constant is the maximum value the JPY/USD couple can reach.
const jpy_usd_MAX float32 = 0.009

// This constant is the minimum value the EUR/USD couple can reach.
const eur_usd_MIN float32 = 1.0

// This constant is the minimum value the GPB/USD couple can reach.
const gpb_usd_MIN float32 = 1.0

// This constant is the minimum value the JPY/USD couple can reach.
const jpy_usd_MIN float32 = 0.006

func main() {
	var wg = sync.WaitGroup{}

	// This variable is the channel for receiving the values of the simulated market EUR/USD.
	eur_usd := make(chan float32, 60)
	// This variable is the channel for receiving the values of the simulated market GPB/USD.
	gbp_usd := make(chan float32, 60)
	// This variable is the channel for receiving the values of the simulated market JPY/USD.
	jpy_usd := make(chan float32, 60)
	// This variable is the number of seconds the simulation is run for.
	time := 60

	// Adding one to the WaitGroup for each simulated market.
	wg.Add(4)
	// Starting all the simulated markert routines.
	go simulateMarketData(&wg, eur_usd, eur_usd_MIN, eur_usd_MAX, time)
	go simulateMarketData(&wg, gbp_usd, gpb_usd_MIN, gbp_usd_MAX, time)
	go simulateMarketData(&wg, jpy_usd, jpy_usd_MIN, jpy_usd_MAX, time)
	// Starting the Sale/Purchase couples routine.
	go selectPair(&wg, eur_usd, gbp_usd, jpy_usd, time)

	// Waiting for all the routines added to the WaitGroup to end.
	wg.Wait()
	fmt.Println("\nThe simulated market has finished.")
}

/*
	This function simulate the market couple with all the given parameters for a specified time.

It takes the following arguments:
-The WaitGroup pointer (wg).
-The couple value sending channel (ch).
-The end receiving signal channel.
-The maximim value of the simulated couple.
-The mininum value of the simulated couple.
-The amount of time the simulated market should exists for.
*/
func simulateMarketData(wg *sync.WaitGroup, ch chan<- float32, min float32,
	max float32, simulation_time int) {
	// This factor is used to set the possible values (minimum and maximum) of the market data.
	factor := (max - min)

	for i := 0; i < simulation_time; i++ {
		// This variable simulate the value of the market.
		// (by summing a random value in interval (0, max-min) to min).
		val := rand.Float32()*(factor) + min
		// Sending the value back through the given receive channel "ch".
		ch <- val
		// Sleeping one second after generating next value.
		time.Sleep(1 * time.Second)
	}
	// Closing the sending channel "ch".
	close(ch)
	// Telling the other routines using the WaitGroup that this process has ended its execution.
	wg.Done()
}

/*
	This function simulate the market based decisions to either buying or selling the couple.

It takes the following arguments:
-The WaitGroup pointer (wg).
-The EUR/USD couple value receiving channel (ch).
-The GPB/USD couple value receiving channel (ch).
-The JPY/USD couple value receiving channel (ch).
*/
func selectPair(wg *sync.WaitGroup, eu <-chan float32, gp <-chan float32, jp <-chan float32,
	simulation_time int) {
	// Checking to see if the simulation time is a positive number.
	if simulation_time < 0 {
		panic("The time simulation must be a positive number!")
	}
	// This variable is the previous value of the couple EUR/USD.
	var last_eur_usd float32 = 0
	// This variable is the previous value of the couple GPB/USD.
	var last_gpb_usd float32 = 1.5
	// This variable is the previous value of the couple JPY/USD.
	var last_jpy_usd float32 = 0.009
	// This variable is the current value of the couple EUR/USD.
	var current_eur_usd float32 = 0.0
	// This variable is the current value of the couple GPB/USD.
	var current_gpb_usd float32 = 0.0
	// This variable is the current value of the couple JPY/USD.
	var current_jpy_usd float32 = 0.0

	// This variable increase by one every second, used to index the operation and value arrays.
	index := 0
	// This variable store a value indicating which one of the possible operations was done.
	// (0 --> not done, 1 --> Sale EUR/USD, 2 --> Purchase GPB/USD, 3 --> Purchase JPY/USD ).
	var operation [64]int = [64]int{0}
	// This variable store the value of the couple when purchasing.
	var value [64]float32 = [64]float32{0.0}

	fmt.Println("+------------------------------+")
	fmt.Println("| OPERATION | COUPLES | VALUES |")
	fmt.Println("|------------------------------|")
	for i := 0; i < simulation_time; i++ {
		// Reading the value of the couple EUR/USD.
		current_eur_usd = <-eu
		// Checking the presence of sale condition.
		if (current_eur_usd > eur_usd_MAX_Sale) && (float32(last_eur_usd) < eur_usd_MAX_Sale) {
			if operation[index] == 0 {
				// The addition of 3 to the index means that the operation is scheduled after 3 seconds.
				operation[index+3] = 1
				value[index+3] = current_eur_usd
			}
		}

		// Reading the value of the couple GPB/USD.
		current_gpb_usd = <-gp
		// Checking the presence of purchase condition.
		if (current_gpb_usd < gbp_usd_MIN_Purchase) && (float32(last_gpb_usd) > gbp_usd_MIN_Purchase) {
			if operation[index] == 0 {
				// The addition of 3 to the index means that the operation is scheduled after 3 seconds.
				operation[index+3] = 2
				value[index+3] = current_gpb_usd
			}
		}
		// Reading the value of the couple JPY/USD.
		current_jpy_usd = <-jp
		// Checking the presence of purchase condition.
		if (current_jpy_usd < jpy_usd_MIN_Purchase) && (float32(last_jpy_usd) > jpy_usd_MIN_Purchase) {
			if operation[index] == 0 {
				// The addition of 4 to the index means that the operation is scheduled after 3 seconds.
				operation[index+4] = 3
				value[index+4] = current_jpy_usd
			}
		}

		// Writing the (formatted) values related to the operation done to the printing routine.
		if operation[index] == 1 {
			// Considering the sale EUR/USD.
			fmt.Printf("| SALE      | EUR/USD | %.4f |\n", value[index])
			fmt.Println("|------------------------------|")
		} else if operation[index] == 2 {
			// Considering the purchase GPB/USD.
			fmt.Printf("| PURCHASE  | GPB/USD | %.4f |\n", value[index])
			fmt.Println("|------------------------------|")
		} else if operation[index] == 3 {
			//Considering the purchase JPY/USD.
			fmt.Printf("| PURCHASE  | JPY/USD | %.4f |\n", value[index])
			fmt.Println("|------------------------------|")
		}

		// Setting the last values of the couples:
		last_eur_usd = current_eur_usd
		last_gpb_usd = current_gpb_usd
		last_jpy_usd = current_jpy_usd

		// Incrementing the value of value and operation index.
		index++
		// Sleeping one second after trying to do another operation.
		time.Sleep(1 * time.Second)
	}

	// Telling the other routines using the WaitGroup that this process has ended its execution.
	wg.Done()
}
