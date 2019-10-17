package main

import (
	"os"
	"fmt"
	"bytes"
	"time"
	"sync"
)

type LoggerDetails struct{
	src *os.File
	mtx sync.Mutex
	logfileCount int
	createdTime time.Time

	max_file_size int64
	max_file_time_in_min time.Duration
	max_file_cout int
}

var logger LoggerDetails

func startTimerWatcher() {
	fmt.Println("Time watcher started for log file ", logger.src.Name())

	defer func(){
		fmt.Println("Time watcher is over")
		logger.mtx.Unlock()
		WG.Done()
	}()

	for ;true; {
		logger.mtx.Lock()
		fmt.Println("Time diff is : ", time.Now().Sub(logger.createdTime))
		fileInfo, err := logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}
		if time.Now().Sub(logger.createdTime) >= (logger.max_file_time_in_min * time.Minute) {
			fmt.Println("Timer expired : file size is ", fileInfo.Size())

			closeFile()
			WG.Add(1)
			go ZipLogFile(fileInfo.Name())
			_ = createFile()
			break
		}
		logger.mtx.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func startSizeWatcher(){
	fmt.Println("Size watcher started for log file ", logger.src.Name())
	defer func(){
		fmt.Println("Size watcher is over")
		logger.mtx.Unlock()
		WG.Done()
	}()

	for ;true; {
		logger.mtx.Lock()

		fileInfo, err := logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}

		if fileInfo.Size() > logger.max_file_size {
			fmt.Println("MAx File size reached, Current size is ", fileInfo.Size())
			closeFile()
			WG.Add(1)
			go ZipLogFile(fileInfo.Name())
			_ = createFile()
			break
		}

		logger.mtx.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func createFile( ) ( err error){

	if logger.logfileCount >= logger.max_file_cout {
		return
	}

	fileName := fmt.Sprintf("%s%d%d%d%d%d%d", "mylog.file.",
		time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	fmt.Println("creating file with name ", fileName)

	logger.src, err = os.Create(fileName)
	if err != nil{
		fmt.Println("Error while creating file with name : ", fileName)
	}

	logger.createdTime = time.Now()
	WG.Add(1)
	go startTimerWatcher()

	WG.Add(1)
	go startSizeWatcher()

	return
}

func closeFile(){
	logger.src.Close()
	logger.logfileCount++
}

func InitLogFileDetails( file_time, file_size int){

	fmt.Println("Initializing logger")

	logger.logfileCount = 0
	logger.src = nil
	logger.max_file_cout = 1
	logger.max_file_size = int64(file_size) * 1024 *1024
	logger.max_file_time_in_min = time.Duration(file_time)

	_ = createFile()
}

func WriteLog( msg string)( count int){
	logger.mtx.Lock()
	if logger.src != nil {
		logger.src.Write(bytes.NewBufferString(msg).Bytes())
	}
	count = logger.logfileCount
	logger.mtx.Unlock()

	return
}


