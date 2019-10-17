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

	stop_size_watcher bool
	stop_time_watcher bool

	wg sync.WaitGroup
}

var Logger LoggerDetails

func startTimerWatcher() {
	fmt.Println("Time watcher started for log file ", Logger.src.Name())

	defer func(){
		fmt.Println("Time watcher is over")
		Logger.mtx.Unlock()
		WG.Done()
	}()

	for ; true; {
		Logger.mtx.Lock()
		//fmt.Println("Time diff is : ", time.Now().Sub(Logger.createdTime))

		if Logger.stop_time_watcher == true {
			break;
		}

		fileInfo, err := Logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}
		if time.Now().Sub(Logger.createdTime) >= (Logger.max_file_time_in_min * time.Minute) {
			fmt.Println("Timer expired : file size is ", fileInfo.Size())
			Logger.stop_size_watcher = true
			closeFile()
			Logger.wg.Add(1)
			WG.Add(1)
			go ZipLogFile(fileInfo.Name())
			go createFile()
			break
		}
		Logger.mtx.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func startSizeWatcher(){
	fmt.Println("Size watcher started for log file ", Logger.src.Name())
	defer func(){
		fmt.Println("Size watcher is over")
		Logger.mtx.Unlock()
		WG.Done()
	}()

	for ; true ; {
		Logger.mtx.Lock()

		if Logger.stop_size_watcher == true {
			break;
		}

		fileInfo, err := Logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}

		if fileInfo.Size() > Logger.max_file_size {
			fmt.Println("Max File size reached, Current size is ", fileInfo.Size())
			Logger.stop_time_watcher = true
			closeFile()
			Logger.wg.Add(1)
			WG.Add(1)
			go ZipLogFile(fileInfo.Name())
			go createFile()
			break
		}

		Logger.mtx.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
}

func createFile( )(){

	Logger.wg.Wait()
	Logger.mtx.Lock()

	if Logger.logfileCount >= Logger.max_file_cout {
		Logger.mtx.Unlock()
		return
	}

	Logger.stop_size_watcher = false
	Logger.stop_time_watcher = false

	fileName := fmt.Sprintf("%s%d%d%d%d%d%d", "mylog.file.",
		time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	fmt.Println("creating file with name ", fileName)

	var err error
	Logger.src,  err = os.Create(fileName)
	if err != nil{
		fmt.Println("Error while creating file with name : ", fileName)
	}

	Logger.createdTime = time.Now()
	Logger.mtx.Unlock()

	WG.Add(1)
	go startTimerWatcher()

	WG.Add(1)
	go startSizeWatcher()

	return
}

func closeFile(){
	Logger.src.Close()
	Logger.logfileCount++
}

func InitLogFileDetails( file_time, file_size, file_count int){

	fmt.Println("Initializing logger")

	fmt.Println("max file size is :", file_size)
	fmt.Println("max file time is :", file_time)
	fmt.Println("max file count is :", file_count)

	Logger.logfileCount = 0
	Logger.src = nil
	Logger.max_file_cout = file_count
	Logger.max_file_size = int64(file_size) * 1024 *1024
	Logger.max_file_time_in_min = time.Duration(file_time)

	go createFile()
}

func WriteLog( msg string)( count int){
	Logger.mtx.Lock()
	if Logger.src != nil {
		Logger.src.Write(bytes.NewBufferString(msg).Bytes())
	}
	count = Logger.logfileCount
	Logger.mtx.Unlock()

	return
}


