/*
author: Mattia Scantamburlo
NOTE:...
*/

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

func main() {
	var clients [10]Cliente
	var mapCar map[string]int = make(map[string]int) //to keep track of the number of car borrowd for each type
	var keys []string = []string{"Berlina", "SUV", "Station Wagon"}
	var wg sync.WaitGroup

	mapCar[keys[0]] = 0
	mapCar[keys[1]] = 0
	mapCar[keys[2]] = 0

	//naming clients
	for i := 0; i < len(clients); i++ {
		clients[i].nome = strconv.Itoa(i) //give the client a number as a name
	}

	//borrowing car phase
	for i := 0; i < len(clients); i++ {
		wg.Add(1)
		go noleggia(&wg, clients[i], mapCar, keys)
	}
	wg.Wait()

	//printing phase
	fmt.Println()
	stampa(mapCar)
}

// allow a client to borrow a car
func noleggia(wg *sync.WaitGroup, cl Cliente, map1 map[string]int, keys []string) {
	r := rand.Intn(len(keys))
	cl.tipo = keys[r]
	map1[keys[r]]++
	fmt.Printf("Il cliente %s ha noleggiato il veicolo %s \n", cl.nome, cl.tipo)
	wg.Done()
}

// print the number od car borrowed divided by type
func stampa(map1 map[string]int) {
	fmt.Println("Il numero di veicoli noleggiati per tipo:")
	for k, v := range map1 {
		fmt.Printf("%s : %d \n", k, v)
	}
}

// associating a Cliente element a Veicolo element using composition
type Cliente struct {
	Veicolo
	nome string
}

type Veicolo struct {
	tipo string
}
