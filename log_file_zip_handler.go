package main

import (
	"os"
	"archive/zip"
	"fmt"
	"io"
)

func ZipLogFile(filename string){

	defer func () {
		fmt.Println("zip file created")
		WG.Done()
	}()

	zipFileName := filename + ".zip"
	zipfile, err := os.OpenFile(zipFileName, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Failed to open zip for writing: ", err)
		return
	}
	defer zipfile.Close()

	zipw := zip.NewWriter(zipfile)
	defer zipw.Close()

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Failed to open ", filename, " : ", err)
		return
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		fmt.Println("Failed to create entry for", filename, "in zip file: ", err)
		return
	}

	if _, err := io.Copy(wr, file); err != nil {
		fmt.Println("Failed to write ", filename, " to zip: ", err)
		return
	}
}