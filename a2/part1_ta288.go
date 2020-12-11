// 1

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
			fmt.Printf("\nPatient %d has been treated", <-p)
		//if waiting room is empty
		default:
			//dentist sleeps and waits for patient
			fmt.Printf("\nDentist is asleep.")
			p := <-dent

			//when p arrives from dent
			time.Sleep(2 * time.Second)
			fmt.Printf("\nPatient %v has been treated", <-p)
		}
	}
}

func patient(wait chan<- chan int, dent chan<- chan int, id int) {

	//patient represented as this int channel
	p := make(chan int)

	fmt.Printf("\nPatient %v wants to be treated.", id)

	select {
	//check if dentist is asleep
	case dent <- p:
		fmt.Printf("\nPatient %v is having treatment.", id)

		//patient falls asleep
		p <- id
	//if dentist is busy, join waiting room
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

	n := 5
	m := 7

	dent := make(chan chan int)
	wait := make(chan chan int, n)

	go dentist(wait,dent)
	time.Sleep(3 * time.Second)

	for i:=0 ; i<m ; i++ {
		go patient(wait, dent, i)
		time.Sleep(1 * time.Second)
	}

	time.Sleep(3000 * time.Millisecond)

}


//old code, keep just in case,
//formatiing on the refactor looks better.
		//if the channel wait is empty.
		// if len(wait) == 0 {

		// 	//dentist sleeps and waits for patient
		// 	fmt.Printf("\nDentist is asleep.")
		// 	p := <-dent

		// 	//when p arrives from dent
		// 	time.Sleep(2 * time.Second)
		// 	fmt.Printf("\nPatient %v has been treated", <-p)
		// }

			
	
