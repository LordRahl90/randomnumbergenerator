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

var fileChannel chan string

func main() {

	runtime.GOMAXPROCS(4)
	start := time.Now()
	limit := 100000

	fileCount := 0

	msgChannel := make(chan Message, 10)
	fileChannel = make(chan string, 50)

	innerFolderPre := "inner_"

	for l := 1; l <= 10; l++ {
		go fetchFromChannel(msgChannel)
	}

	innerFolder := innerFolderPre + strconv.Itoa(fileCount)

	for fileCount < limit {

		if fileCount%100 == 0 {
			innerFolder = innerFolderPre + strconv.Itoa(fileCount)
		}
		if _, err := os.Stat(innerFolder); os.IsNotExist(err) {
			os.MkdirAll(folder+"/"+innerFolder, os.ModePerm)
		}
		//create a new file
		newFile := folder + "/" + innerFolder + "/file_" + strconv.Itoa(fileCount)
		fileChannel <- newFile
		go writeToChannel(msgChannel, fileChannel)

		fileCount++
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed Time is: %s\n", elapsed.String())

	fmt.Scanln()
}

func fetchFromChannel(fileCh chan Message) {
	for {
		message := <-fileCh
		writeToFile(message)
	}
}

func writeToChannel(messageChan chan Message, fileChannel chan string) {
	filename := <-fileChannel

	for i := 1; i <= 10000; i++ {
		var buffer bytes.Buffer
		for j := 1; j <= 100; j++ {
			number := rand.Intn(1001)

			buffer.WriteString(strconv.Itoa(number) + ",")
		}

		message := Message{
			Filename: filename,
			Number:   buffer.String() + "\n",
		}

		messageChan <- message
	}
}

//spins a new go routine that writes to file
func writeToFile(message Message) {
	fileName := message.Filename
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(message.Number)
	if err != nil {
		log.Fatal(err)
	}
}

//Message Struct
type Message struct {
	Filename string
	Number   string
}
