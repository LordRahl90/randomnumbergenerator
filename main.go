package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

const folder = "./content"

func main() {

	runtime.GOMAXPROCS(4)
	start := time.Now()
	limit := 100

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}

	fileCount := 0
	fileChannel := make(chan Message, 50)

	go fetchFromChannel(fileChannel)

	for fileCount < limit {
		//create a new file
		newFile := folder + "/file_" + strconv.Itoa(fileCount)
		go writeToChannel(fileChannel, newFile)
		fileCount++
	}

	elapsed := time.Since(start)
	fmt.Printf("Elapsed Time is: %s\n", elapsed.String())

	fmt.Scanln()
}

func fetchFromChannel(fileCh chan Message) {
	for {
		message := <-fileCh
		go writeToFile(message)
	}
}

func writeToChannel(fileCh chan Message, fileName string) {
	for i := 1; i <= 1000000; i++ {
		source := rand.NewSource(time.Now().UnixNano())
		number := rand.New(source).Int()
		message := Message{
			Filename: fileName,
			Number:   number,
		}

		fileCh <- message
	}
}

//spins a new go routine that writes to file
func writeToFile(message Message) {
	// fmt.Printf("Recieved Filename:%s, Number: %d\n", message.Filename, message.Number)
	// return
	go func() {
		fileName := message.Filename
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(strconv.Itoa(message.Number) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}()
}

//Message Struct
type Message struct {
	Filename string
	Number   int
}
