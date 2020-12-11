// 2.a

package main

import (
	"fmt"
	"time"
)

func dentist(h_wait chan chan int , l_wait <-chan chan int, dent <-chan chan int) {

	//start a timer
	timer := time.NewTimer(4 * time.Second)

	//loop forever.
	for {
		//check high priority queue first, if its empty
		if len(h_wait) == 0 {

			//then check if low priority queue, if this is empty
			if len(l_wait) == 0 {
				//dentist sleeps and waits for patient
				fmt.Printf("\nDentist is asleep.")
				p := <-dent

				//when p arrives from dent, do some work
				time.Sleep(2 * time.Second)
				//and wake up the patient
				fmt.Printf("\nPatient %v has been treated", <-p)
			} else {
				//if there is no high priority patients, but low priority patients
				//see l_wait patients
				select {
					//patient leaves l_wait room
					case <- l_wait:	
						p := <- dent
						//dentist releases the patient
						time.Sleep(2 * time.Second)
						fmt.Printf("\nPatient %d has been treated from l", <-p)
				}

			}

		} else {
			// if high priority treatment is required
			select {
				//patient leaves h_wait room
				case <- h_wait:	
					p := <- dent
					//dentist releases the patient
					time.Sleep(2 * time.Second)
					fmt.Printf("\nPatient %d has been treated from h", <-p)
			}

		}

		select {
		//if timer ends
		case <-timer.C:
			//and low p patients are waiting
			if len(l_wait) > 0 {
				//move the patients
				p := <- l_wait
				h_wait <- p
				fmt.Printf("\nPatient has been moved from low to high")
				//restart the timer
				timer = time.NewTimer(4 * time.Second)
			}
		default:
		}
	}
}

func patient(wait chan<- chan int, dent chan<- chan int, id int) {

	p := make(chan int)

	fmt.Printf("\nPatient %v wants to be treated.", id)

	
	select {
	//check if the dentist is busy
	case dent <- p:
		fmt.Printf("\nPatient %v is having treatment.", id)

		//patient falls asleep
		p <- id
	default:
		fmt.Printf("\nPatient %v is waiting.", id)
		//patient goes to waiting room
		wait <- p

		//patient sees the dentist
		dent <- p 

		//patient falls asleep
		fmt.Printf("\nPatient %v is having treatment.", id)
		p <- id
	}
	
}

func main() {

	dent := make(chan chan int)
	
	h_wait := make(chan chan int, 100)
	l_wait := make(chan chan int, 5)

	h := 10
	l := 3

	go dentist(h_wait, l_wait, dent)
	time.Sleep(3 * time.Second)

	for i:= l ; i < h ; i++ {
		go patient(h_wait, dent, i)
	}

	for i := 0 ; i < l ; i++ {
		fmt.Printf("\nPatient %v is a low p patient.", i)
		go patient(l_wait, dent, i)
	}

	time.Sleep(50 * time.Second)

}

// 2.b
// Let's assume that go doesn't have fairness.
// If there was a fixed number of patients, they would all be served
// just not in a fair order.

// If there was an infinite amount of patient, or a constant flow of patients
// then some patients may never get served.

// The best approach would be to implement a some kind of queue, using ageing to make sure
// low priority patients also get served.


			
	
