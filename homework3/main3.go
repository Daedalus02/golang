/*
author: Mattia Scantamburlo
NOTE:
1)	Ho preferito rendere il programma più chiaro nello sviluppo delle routines associando delle variabili di stato
	a ciscun cuoco che vengono stampate circa ogni secondo (nonostante questo aggiunga un costo computazionale).
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

// this is the function that manage the process of cooking the decided number of cakes
func main() {
	time_init := time.Now()
	fmt.Println("Starting timer...")

	//end channel
	chEnd := make(chan bool, 4)

	//channel 12 (cook1 <--> cook2)
	mut12 := mutCh{
		ch:      make(chan int, 2),
		counter: 2,
		id:      "12",
	}

	//channel 23 (cook2 <--> cook3)
	mut23 := mutCh{
		ch:      make(chan int, 2),
		counter: 2,
		id:      "23",
	}
	wg := sync.WaitGroup{}
	wg.Add(3) //adding element to wait group
	go cook1(&wg, &mut12, chEnd, 5)
	go cook2(&wg, &mut12, &mut23, chEnd, 5)
	go cook3(&wg, &mut23, chEnd, 5) //notice: cook3 is also responsibnle for telling other cooks when to stop
	wg.Wait()                       //waiting unitil the go routines get completed
	close(chEnd)
	fmt.Printf("It took the 3 cooks: %f32 seconds", float32((time.Now().Sub(time_init)).Seconds()))
}

// COOKS THE CAKES (COOK1)
func cook1(wg *sync.WaitGroup, mut12 *mutCh, chEnd <-chan bool, limit int) {
	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}
	//state variables of cook1
	time1 := time.Now()
	done1 := 0
	for {
		if done1 == limit {
			break
		}
		//COOKING CAKES PHASE
		if time.Now().Sub(time1) >= 1*time.Second {
			if err := mut12.send(1); err == nil { //signal cook2 that a cake is ready if possible
				time1 = time.Now() //"zerify" the time gap
				done1++
				fmt.Printf("\nDone1: %d \n", done1)
			} else {
				//fmt.Println(err) //possible printing of error
			}
		}
		time.Sleep(50 * time.Millisecond) //delay: to avoid doing to many check operation
	}

	//waiting to get an end signal
	for {
		if <-chEnd {
			break
		}
		time.Sleep(50 * time.Millisecond) //delay: to avoid checking necessaire resources constantly
	}
	mut12.close() //closing channel: wait until mutex is unlocked
	wg.Done()     //decrease wg count
}

// GUARNISH CAKES (COOK2)
func cook2(wg *sync.WaitGroup, mut12 *mutCh, mut23 *mutCh, chEnd <-chan bool, limit int) {
	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}

	//state variable of cook2
	time2 := time.Now()
	done2 := 0
	available2 := 0
	cooking2 := false
	for {
		if done2 == limit {
			break
		}

		if available2 < 2 {
			if err, value := mut12.receive(); err == nil { //if possible, chek to see if some cakes are ready
				if value > 0 {
					available2 += value
				}
				if available2 > 2 {
					panic("Cook2 only have 2 places!")
				}
			} else {
				//fmt.Println(err) //possible printing of error
			}
		}
		if !cooking2 && available2 != 0 { //checks if guarnishing cakes conditions are present
			cooking2 = true
			time2 = time.Now()
		}

		//GUARNICHING CAKES PHASE
		if time.Now().Sub(time2) >= 4*time.Second && cooking2 {
			if err := mut23.send(1); err == nil {
				done2++
				available2--
				cooking2 = false
				fmt.Printf("\nDone2: %d\n", done2)
			} else {
				//fmt.Println(err) //possible printing of error
			}
		}

		//fmt.Printf("\nAvailable2: %d\n", available2)
		time.Sleep(50 * time.Millisecond) //delay: to avoid checking necessaire resources constantly
	}
	for {
		if <-chEnd {
			break
		}
		time.Sleep(50 * time.Millisecond) //delay: to avoid checking end channel constantly
	}
	mut23.close() //closing channel: wait until mutex is unlocked
	wg.Done()     //decrease wg count
}

func cook3(wg *sync.WaitGroup, mut23 *mutCh, chEnd chan<- bool, limit int) {

	if limit < 0 {
		panic("Cake orders must be positive numbers!")
	}

	//condition variable
	done3 := 0
	available3 := 0
	time3 := time.Now()
	cooking3 := false

	for {
		if done3 == limit {
			break
		}
		if available3 < 2 {
			if err, value := mut23.receive(); err == nil { // if possible, checks to see if some cakes are ready
				if value > 0 {
					available3 += value
				}
				if available3 > 2 {
					panic("Cook3 only have 3 places!")
				}
			} else {
				//fmt.Println(err) // possible printing of error
			}
		}

		if !cooking3 && available3 != 0 { //checks if decorating cakes conditions are present
			cooking3 = true
			time3 = time.Now()
		}

		//DECORATING CAKES PHASE
		if time.Now().Sub(time3) >= 8*time.Second && cooking3 {
			available3--
			done3++
			cooking3 = false
			fmt.Printf("\nDone3: %d\n", done3)
		}
		//fmt.Printf("\nAvailable3: %d\n", available3)
		time.Sleep(50 * time.Millisecond) //delay: to avoid checking necessaire resources constantly
	}

	//signaling that the cakes are completely done
	chEnd <- true
	chEnd <- true
	wg.Done() //decrease wg count
}

// mutual exclusive channel
type mutCh struct {
	ch      chan int
	counter int        // number of free spaces in the channel (MAX: 2, MIN: 0)
	m       sync.Mutex // to synchronize send and read operations on the same channel
	id      string     //to identifie the channel (NB: this is not used unless the print statement is uncommented)
}

// send values when possible through the channel
func (mut *mutCh) send(value int) error {
	mut.m.Lock()
	//fmt.Println("Sending ...")
	//fmt.Printf("counter sender from %s: %d\n", mut.id, mut.counter)
	if mut.counter == 0 {
		mut.m.Unlock()
		//fmt.Println("finished Sending not possible...")
		//fmt.Println()
		return fmt.Errorf("Not possible to send data")
	} else {
		mut.ch <- value
		mut.counter--
		//fmt.Printf("counter sender from %s after writing : %d\n", mut.id, mut.counter)
	}
	mut.m.Unlock()
	//fmt.Println("finished Sending done smth...")
	//fmt.Println()
	return nil
}

// read value from the channel when possible
func (mut *mutCh) receive() (error, int) {
	var value int
	mut.m.Lock()

	//fmt.Println("Receiving ...")
	//fmt.Printf("counter receiver from %s: %d\n", mut.id, mut.counter)

	if mut.counter == 2 {
		mut.m.Unlock()
		//fmt.Println("finished receiving not possible...")
		//fmt.Println()
		return fmt.Errorf("Not possible to receive data"), -1
	} else {
		select {
		case value = <-mut.ch:
			mut.counter++
			//fmt.Printf("counter receiver from %s after reading : %d\n", mut.id, mut.counter)
		default:
			mut.m.Unlock()
			//fmt.Println("finished receiving in select...")
			//fmt.Println()
			return fmt.Errorf("No message"), -2
		}
	}
	//fmt.Println("finished receiving done smth...")
	//fmt.Println()
	mut.m.Unlock()
	return nil, value
}

func (mut *mutCh) close() error {
	mut.m.Lock()
	close(mut.ch)
	//fmt.Printf("\nclosed n: %s\n", mut.id)
	mut.m.Unlock()
	return nil
}
