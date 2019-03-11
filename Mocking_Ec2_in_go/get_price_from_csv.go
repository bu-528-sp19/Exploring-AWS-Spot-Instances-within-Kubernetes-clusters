package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"log"
)

func main(){
	csvFile, _ := os.Open("Mock.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.ReadAll()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		fmt.Println(line[0][1])
	}
}
