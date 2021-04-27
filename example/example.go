package main

import (
	"fmt"
	"math/rand"
	"os"

	invert "github.com/Noahnut/invertIndex"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testFileGenerator(fileNumber int) {
	fileBaseName := "test"
	for i := 0; i < fileNumber; i++ {
		f, err := os.Create(fileBaseName + fmt.Sprint(i))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		f.WriteString("ID:" + "111" + "\n")
		f.WriteString("Name:" + "222" + "\n")
		f.WriteString("title:" + "333" + "\n")
	}

}

func main() {
	testFileGenerator(5)
	in := invert.InvertIndex{}
	in.NewInvertIndex("ID", "Name", "title")
}
