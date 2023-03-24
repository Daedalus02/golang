/*
author: Mattia Scantamburlo

NOTE:
1)    Ho deciso di far stampare le vendite nel corpo della funzione main così da rendere eventualmente posssibile
	  coordinare la stampa combinata dei valori di andamento delle valute assieme a quelli al momento della vendita.
2)    Ho tenuto i buffer dei channel abbastanza ampi per evitare perdite di segnali.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const eur_usd_MAX float32 = 1.20

const gbp_usd_MIN float32 = 1.35

const jpy_usd_MIN float32 = 0.0085

func main() {
	var wg = sync.WaitGroup{}

	eur_usd := make(chan float32, 60)
	gbp_usd := make(chan float32, 60)
	jpy_usd := make(chan float32, 60)
	end := make(chan bool, 10)
	out := make(chan string, 100)

	wg.Add(3)
	go simulateMarketData(&wg, eur_usd, end, 1, 1.5)
	go simulateMarketData(&wg, gbp_usd, end, 1, 1.5)
	go simulateMarketData(&wg, jpy_usd, end, 0.006, 0.009)
	go selectPair(&wg, eur_usd, gbp_usd, jpy_usd, end, out, 1000)

	for i := 0; i < 60; i++ {
		time.Sleep(1 * time.Second)
		//fmt.Printf("eur_usd: %f32 \njpy_usd: %f32 \ngbp_usd: %f32 \n", <-eur_usd, <-jpy_usd, <-gbp_usd)
		//fmt.Println()
		select {
		case x := <-out:
			fmt.Printf("\n%s\n", x)
		default:
		}
	}

	end <- true
	end <- true
	end <- true
	end <- true
	wg.Wait()
	close(end)
}

func simulateMarketData(wg *sync.WaitGroup, ch chan<- float32, end <-chan bool, min float32, max float32) {
	ender := false
	fact := (max - min)

	for {
		select {
		case x := <-end:
			ender = x
		default:
		}
		if ender {
			break
		}
		val := rand.Float32()*(fact) + min
		time.Sleep(1000 * time.Millisecond)
		ch <- val
	}
	close(ch)
	wg.Done()
}

func selectPair(wg *sync.WaitGroup, eu <-chan float32, gp <-chan float32, jp <-chan float32, end <-chan bool, out chan<- string, amount float32) {
	ender := false
	index := 0
	var order [64]int = [64]int{0}
	var fl [64]float32 = [64]float32{0.0}
	for {
		select {
		case y := <-end:
			ender = y
		case x := <-eu:
			if x > eur_usd_MAX {
				if order[index] == 0 {
					order[index+3] = 1
					fl[index+3] = x
				}
			}
		case x := <-gp:
			if x < gbp_usd_MIN {
				if order[index] == 0 {
					order[index+3] = 2
					fl[index+3] = x
				}
			}
		case x := <-jp:
			if x < jpy_usd_MIN {
				if order[index] == 0 {
					order[index+4] = 3
					fl[index+4] = x
				}
			}
		default:
		}

		if order[index] == 1 {
			out <- fmt.Sprintf("Vendita EUR/USD, valore: %f32 \n ", fl[index])
		} else if order[index] == 2 {
			out <- fmt.Sprintf("Acquisto GPB/USD, valore: %f32 \n ", fl[index])
		} else if order[index] == 3 {
			out <- fmt.Sprintf("Acquisto JPY/USD, valore: %f32 \n ", fl[index])
		}

		if ender {
			break
		}

		index++
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}
