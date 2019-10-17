package main

import (
	"fmt"
	"sync"
	"flag"
)

var WG sync.WaitGroup


func startLogger( file_time, file_size, file_count * int){
	fmt.Println("starting logger")
	defer func () {
		fmt.Println("Logger is done")
		WG.Done()
	}()

	InitLogFileDetails( *file_time, *file_size, *file_count)
	count := 0

	for ; count < *file_count; {
		count = WriteLog(" Testing with random msg...........")
		//time.Sleep( 1 * time.Second)
	}
}

func main() {
	fmt.Println("Creating log file with rotation")

	file_size := flag.Int("log_file_size", 1, "This is max file size in MB")
	file_time := flag.Int("log_file_time", 1,"This is max file time in min")
	file_count := flag.Int("log_file_count", 2,"This is max file count")

	flag.Parse()

	WG.Add(1)

	go startLogger(file_time, file_size, file_count)

	WG.Wait()

	fmt.Println("Main is over")
}
