/*
Author: Mattia Scantamburlo.

Goal: Write a program in GO that simulates a car rental agency that has to manage
bookings of 10 customers. Each customer hires a vehicle from those available:
Sedan, SUV or Station Wagon.

NOTE: Since more than once, even with the go routines operating separately, the execution
ended up with the error: "Golang fatal error: concurrent map read and map write".
So i decided to make the map carMap a mutexMap to avoid having the concurrent accesses.
*/

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

func main() {
	// This variable contains all the customers considered.
	var customers [10]Cliente
	// This variable is used to keep track of the number of car borrowed for each type of car.
	var mapCar mutexMap
	mapCar.carMap = make(map[string]int)
	// This variable contains all the car names used as keys in the map "mapCar".
	var keys []string = []string{"Berlina", "SUV", "Station Wagon"}
	var wg sync.WaitGroup

	// Assigning all the values of the map "mapCar" to 0 before making the customers start borrowing.
	mapCar.carMap[keys[0]] = 0
	mapCar.carMap[keys[1]] = 0
	mapCar.carMap[keys[2]] = 0

	// Naming customers.
	for i := 0; i < len(customers); i++ {
		customers[i].nome = "Cliente_" + strconv.Itoa(i) // Give the customer a number as a name (0,1,2,3,4,5,6,7,8,9).
	}

	// Borrowing car phase.
	for i := 0; i < len(customers); i++ {
		wg.Add(1)
		go noleggia(&wg, customers[i], &mapCar, keys)
	}
	wg.Wait()

	// Printing the number of car borrowed by type at the end of the process.
	defer fmt.Println()
	defer stampa(&mapCar)
}

/*
	This function allow a customer to borrow a car.

It takes as arguments the WaitGroup, the Cliente instance who borrow the car,
the map with the Car borrowed count, and the name of the possible cars for the Cliente
instance to borrow.
*/
func noleggia(wg *sync.WaitGroup, customer Cliente, map1 *mutexMap, keys []string) {
	r := rand.Intn(len(keys))
	customer.tipo = keys[r]
	map1.Lock()
	map1.carMap[keys[r]]++
	map1.Unlock()
	fmt.Printf("The customer %s has borrowed a car of type %s. \n", customer.nome, customer.tipo)
	wg.Done()
}

/*
	This function print the number od car borrowed divided by type.

It takes as argument the map with the Car borrowed count.
*/
func stampa(map1 *mutexMap) {
	fmt.Println("Number of car borrowed ordered by car type:")
	map1.Lock()
	for k, v := range map1.carMap {
		fmt.Printf("-%s : %d \n", k, v)
	}
	map1.Unlock()
}

/*
	Representing the type Cliente with the borrowed car and a name.

NOTE: Associating a Cliente element a Veicolo element using composition.
*/
type Cliente struct {
	Veicolo        // Borrowed car.
	nome    string // Cliente name.
}

/*
Representing the type Veicolo with the name of the car.
*/
type Veicolo struct {
	tipo string // Veicolo name.
}

/*
	This type is used to keep track of the number of car borrowed for each type of car.

It uses a map of integer as values indexed by name of cars. Also, it uses a mutex to ensure
the go routines does not access the map simultaneously.
*/
type mutexMap struct {
	carMap map[string]int
	mutex  sync.Mutex
}

/*
This method is just a dummy for the internal mutex Lock method.
*/
func (mutex_map *mutexMap) Lock() {
	mutex_map.mutex.Lock()
}

/*
This method is just a dummy for the internal mutex Unlock method.
*/
func (mutex_map *mutexMap) Unlock() {
	mutex_map.mutex.Unlock()
}
