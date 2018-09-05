package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

const folder = "./files"

func init() {
	rand.Seed(time.Now().Unix())
}

var limit = 100000
var fileCount = 0

func main() {
	runtime.GOMAXPROCS(4)
	fileChannel := make(chan string, 500)
	var innerFolder string
	innerFolderPre := "inner_"

	go readFromChannel(fileChannel)

	for limit > fileCount {
		if fileCount%1000 == 0 {
			innerFolder = innerFolderPre + strconv.Itoa(fileCount)
		}
		if _, err := os.Stat(innerFolder); os.IsNotExist(err) {
			os.MkdirAll(folder+"/"+innerFolder, os.ModePerm)
		}
		//create a new file
		newFile := folder + "/" + innerFolder + "/file_" + strconv.Itoa(fileCount)
		fileChannel <- newFile
		fileCount++
	}

	fmt.Scanln()
}

func readFromChannel(fileChannel chan string) {
	for {
		filename := <-fileChannel
		// we open the file, create 10 worker threads to work on it

		generateNumber(filename)
	}
}

func generateNumber(fileName string) {
	for i := 1; i <= 10000; i++ {
		var buffer bytes.Buffer
		for j := 1; j <= 100; j++ {
			number := rand.Intn(1001)

			buffer.WriteString(strconv.Itoa(number) + ",")
		}

		go writeToFile(fileName, buffer.String()+"\n")
	}
}

func writeToFile(fileName string, message string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		log.Fatal(err)
	}
}
