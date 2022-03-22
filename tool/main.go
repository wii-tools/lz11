package main

import (
	"github.com/wii-tools/lz11"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("Refer to the lzx(1) manual.")
		os.Exit(-1)
	}

	cmd := os.Args[1]
	filename := os.Args[2]
	output := os.Args[3]

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	switch cmd {
	case "-evb":
		data, err := lz11.Compress(file)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}

		break
	case "-d":
		data, err := lz11.Decompress(file)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}
		break
	default:
		log.Println("This does not please lzx(1). Read its manual.")
		os.Exit(-1)
	}

}
