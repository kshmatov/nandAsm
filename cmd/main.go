package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kshmatov/nandAssembler/internal/parser"
)

func main() {
	out := flag.String("o", "", "output file, defualt stdout")
	in := flag.String("i", "", "source file")
	bin := flag.Bool("b", false, "store in binary format, in in binary string")
	flag.Parse()
	if *in == "" {
		fmt.Printf("no source file is given\n")
		flag.Usage()
		return
	}

	cnt, err := os.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	src := strings.Split(string(cnt), "\n")
	res, err := parser.Parse(src)
	if err != nil {
		fmt.Println(err)
		return
	}
	var df io.Writer
	if *out == "" {
		df = os.Stdout
	} else {
		f, err := os.OpenFile(*out, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		df = f
	}
	if *bin {
		_, err := df.Write(res.Binary())
		if err != nil {
			fmt.Println(err)
			return
		}
		if *out == "" {
			fmt.Println()
		}
	} else {
		for _, s := range res.String() {
			_, err := df.Write([]byte(s + "\n"))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	fmt.Println("Done")
}
