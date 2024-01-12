package main

import (
	"bufio"
	"fmt"
	//"io"
	"strconv"
	//"os"
	"log"
	"os/exec"
)

func main() {

	cmd := exec.Command("/home/gnewton/local/bin/gs", "-sDEVICE=nullpage", "-q", "-dNOPAUSE", "-dBATCH", "../../../kv.ps")

	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Start: could not run command: ", err)
	}

	go func() {
		fmt.Println("Reader: START")
		defer out.Close()

		sc := bufio.NewScanner(out)

		for sc.Scan() {
			fmt.Println(sc.Text())
			if err := sc.Err(); err != nil {
				log.Fatalf("scan file error: %v", err)
				break
			}
		}

		fmt.Println("Reader: END")
	}()

	go func() {
		defer stderr.Close()
		fmt.Println("err: START")

		sc := bufio.NewScanner(stderr)

		for sc.Scan() {
			fmt.Println(sc.Text())
			if err := sc.Err(); err != nil {
				log.Fatalf("scan file error: %v", err)
				break
			}
		}

		fmt.Println("stderr Reader: END")
	}()
	var by int64

	go func() {
		fmt.Println("Writer: start")

		by = 0

		for i := 0; i < 1000000; i++ {
			s := strconv.Itoa(i)
			//key := "p " + "K" + s + "\n"
			key := "p " + s + "\n"
			in.Write([]byte(key))
			by += int64(len([]byte(key)))

			//value := s + "\n"
			value := "V\n"
			in.Write([]byte(value))
			by += int64(len([]byte(value)))
		}

		in.Write([]byte("p 10\n"))
		in.Write([]byte(" apple\n"))
		in.Write([]byte("p foobar\n"))
		in.Write([]byte("butter\n"))
		in.Write([]byte("c\n"))
		in.Write([]byte("g 10\n"))
		in.Write([]byte("c\n"))
		in.Write([]byte("b\n"))
		//in.Write([]byte("k 0 100\n"))
		//in.Write([]byte("D\n"))
		in.Write([]byte("Q\n"))
		in.Close()
		fmt.Println("Writer: END")
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: could not run command: ", err)
	}
	fmt.Println("Bytes=", by/1024/1024, " MB")
}
