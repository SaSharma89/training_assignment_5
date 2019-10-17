package main

import (
	"fmt"
	"sync"
	"flag"
)

var WG sync.WaitGroup


func startLogger( file_size, file_time * int){
	fmt.Println("starting logger")
	defer func () {
		fmt.Println("Logger is done")
		WG.Done()
	}()

	fmt.Println("max file size is :", *file_size)
	fmt.Println("max file time is :", *file_time)

	InitLogFileDetails( *file_time, *file_size)
	count := 0

	for ; count < 1; {
		count = WriteLog(" Testing with random msg...........")
		//time.Sleep( 1 * time.Second)
	}
}

func main() {
	fmt.Println("Creating log file with rotation")

	file_size := flag.Int("log_file_size", 1, "This is max file size in MB")
	file_time := flag.Int("log_file_time", 1,"This is max file time in min")

	flag.Parse()

	WG.Add(1)

	go startLogger(file_time, file_size)

	WG.Wait()

	fmt.Println("Main is over")
}
