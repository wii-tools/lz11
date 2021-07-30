package main

import (
	"encoding/hex"
	"github.com/wii-tools/lz11"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Refer to the lzx(1) manual.")
		os.Exit(-1)
	}

	cmd := os.Args[1]
	filename := os.Args[2]

	switch cmd {
	case "-evb":
		break
	case "-d":
		break
	default:
		log.Println("This does not please lzx(1). Read its manual.")
		os.Exit(-1)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	data, err := lz11.Decompress(file)
	if err != nil {
		panic(err)
	}

	println(hex.EncodeToString(data))
}
