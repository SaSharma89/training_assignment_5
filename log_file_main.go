package main

import (
	"flag"
	"fmt"
	"sync"
)

// WG for sync
var WG sync.WaitGroup

func startLogger(fileTime, fileSize, fileCount *int) {
	fmt.Println("starting logger")
	defer func() {
		fmt.Println("Logger is done")
		WG.Done()
	}()

	InitLogFileDetails(*fileTime, *fileSize, *fileCount)
	count := 0

	for count < *fileCount {
		count = WriteLog(" Testing with random msg...........")
		//time.Sleep( 1 * time.Second)
	}
}

func main() {
	fmt.Println("Creating log file with rotation")

	fileSize := flag.Int("log_fileSize", 1, "This is max file size in MB")
	fileTime := flag.Int("log_fileTime", 1, "This is max file time in min")
	fileCount := flag.Int("log_fileCount", 2, "This is max file count")

	flag.Parse()

	WG.Add(1)

	go startLogger(fileTime, fileSize, fileCount)

	WG.Wait()

	fmt.Println("Main is over")
}
