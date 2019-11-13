package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

// LoggerDetails is main strucure
type LoggerDetails struct {
	src          *os.File
	mtx          sync.Mutex
	logfileCount int
	createdTime  time.Time

	maxFileSize      int64
	maxFileTimeInMin time.Duration
	maxFileCout      int

	stopSizeWatcher bool
	stopTimeWatcher bool

	wg sync.WaitGroup
}

//Logger is handler
var Logger LoggerDetails

func startTimerWatcher() {
	fmt.Println("Time watcher started for log file ", Logger.src.Name())

	defer func() {
		fmt.Println("Time watcher is over")
		Logger.mtx.Unlock()
		WG.Done()
	}()

	for true {
		Logger.mtx.Lock()
		//fmt.Println("Time diff is : ", time.Now().Sub(Logger.createdTime))

		if Logger.stopTimeWatcher == true {
			break
		}

		fileInfo, err := Logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}
		if time.Now().Sub(Logger.createdTime) >= (Logger.maxFileTimeInMin * time.Minute) {
			fmt.Println("Timer expired : file size is ", fileInfo.Size())
			Logger.stopSizeWatcher = true
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

func startSizeWatcher() {
	fmt.Println("Size watcher started for log file ", Logger.src.Name())
	defer func() {
		fmt.Println("Size watcher is over")
		Logger.mtx.Unlock()
		WG.Done()
	}()

	for true {
		Logger.mtx.Lock()

		if Logger.stopSizeWatcher == true {
			break
		}

		fileInfo, err := Logger.src.Stat()
		if err != nil {
			fmt.Println("Error while getting stats of file")
			return
		}

		if fileInfo.Size() > Logger.maxFileSize {
			fmt.Println("Max File size reached, Current size is ", fileInfo.Size())
			Logger.stopTimeWatcher = true
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

func createFile() {

	Logger.wg.Wait()
	Logger.mtx.Lock()

	if Logger.logfileCount >= Logger.maxFileCout {
		Logger.mtx.Unlock()
		return
	}

	Logger.stopSizeWatcher = false
	Logger.stopTimeWatcher = false

	fileName := fmt.Sprintf("%s%d%d%d%d%d%d", "mylog.file.",
		time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())

	fmt.Println("creating file with name ", fileName)

	var err error
	Logger.src, err = os.Create(fileName)
	if err != nil {
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

func closeFile() {
	Logger.src.Close()
	Logger.logfileCount++
}

// InitLogFileDetails called from main
func InitLogFileDetails(fileTime, fileSize, fileCount int) {

	fmt.Println("Initializing logger")

	fmt.Println("max file size is :", fileSize)
	fmt.Println("max file time is :", fileTime)
	fmt.Println("max file count is :", fileCount)

	Logger.logfileCount = 0
	Logger.src = nil
	Logger.maxFileCout = fileCount
	Logger.maxFileSize = int64(fileSize) * 1024 * 1024
	Logger.maxFileTimeInMin = time.Duration(fileTime)

	go createFile()
}

// WriteLog to log a msg
func WriteLog(msg string) (count int) {
	Logger.mtx.Lock()
	if Logger.src != nil {
		Logger.src.Write(bytes.NewBufferString(msg).Bytes())
	}
	count = Logger.logfileCount
	Logger.mtx.Unlock()

	return
}
