package main

import "log"

func main() {
	result_chan2 := make(chan string)
	result_chan3 := make(chan string)
	result_chan4 := make(chan string)
	go record_file_reader_routine("2.txt", result_chan2)
	go record_file_reader_routine("3.txt", result_chan3)
	go record_file_reader_routine("4.txt", result_chan4)
	stop2 := <- result_chan3
	stop1 := <- result_chan2
	stop3 := <- result_chan4
	if stop1 == "ok" && stop2 == "ok" && stop3 == "ok" {
		log.Printf("Finished")
	} else {
		log.Printf("Error: %s %s %s", stop1, stop2, stop3)
	}
}
