// 3
package main

import (
	"fmt"
	"time"
)

func dentist(wait <-chan chan int, dent <-chan chan int) {

	//loop forever.
	for {
		select {
		//patient leaves waiting room
		case <- wait:	
			p := <- dent
			//dentist releases the patient
			time.Sleep(2 * time.Second)
			fmt.Printf("\nPatient %d has been treated", <-p)
		//if the channel wait is empty.
		default:
			//dentist sleeps and waits for patient
			fmt.Printf("\nDentist is asleep.")
			//blocks until a patient arrives
			p := <- dent
			//when p arrives from dent, do some work
			time.Sleep(2 * time.Second)
			//and wake up the patient
			fmt.Printf("\nPatient %v has been treated", <-p)
		}
	
	}
}

func assistant(h_wait chan chan int, l_wait <-chan chan int, wait chan<- chan int) {

	timer := time.NewTimer(5 * time.Second)

	//loop forever
	for {
		select {
			//if timer ends
			case <-timer.C:
				//and low p patients are waiting
				if len(l_wait) > 0 {
					//move the patients low
					p := <- l_wait
					//move them to high
					h_wait <- p
					fmt.Printf("\nAssistant has moved patient from low to high")
					//restart the timer
					timer = time.NewTimer(5 * time.Second)
				}
			// if patients in high wait
			case p := <- h_wait:
				//move them to waiting room
				wait <- p
			//otherwise just fall asleep
			case p := <- l_wait:
				if len(h_wait) == 0 {
					wait <- p
				}
			default:
			}

	}

}

func patient(hl_wait chan<- chan int, dent chan<- chan int, id int) {

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
		hl_wait <- p

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

	wait := make(chan chan int, 100)

	h := 10
	l := 3

	go dentist(wait, dent)
	time.Sleep(3 * time.Second)

	go assistant(h_wait, l_wait, wait)
	// time.Sleep(3 * time.Second)

	for i:= l ; i < h ; i++ {
		go patient(h_wait, dent, i)
	}

	for i := 0 ; i < l ; i++ {
		fmt.Printf("\nPatient %v is a low p patient.", i)
		go patient(l_wait, dent, i)
	}

	time.Sleep(50 * time.Second)

}

//3.b
// if the patient fails for whatever reason, there will be no communication
// between patients and the dentist. Part 2 does not have any type of patient.
// meaning that the it will be stuck in a state where patients are waiting,
// and the dentist will recieve nothing. The assistant in part 3 resolves this issue.


// OLD CODE
// if len(wait) == 0 {

// 	//dentist sleeps and waits for patient
// 	fmt.Printf("\nDentist is asleep.")
// 	p := <-dent

// 	//when p arrives from dent
// 	time.Sleep(2 * time.Second)
// 	fmt.Printf("\nPatient %v has been treated", <-p)


			
	
