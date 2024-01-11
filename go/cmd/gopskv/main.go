package main

import (
	"fmt"
	//"io"
	"bufio"
	//"os"
	"log"
	"os/exec"
)

func main() {

	cmd := exec.Command("gs", "-sDEVICE=nullpage", "-q", "-dNOPAUSE", "-dBATCH", "kv.ps")

	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	out, err := cmd.StdoutPipe()
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
		fmt.Println("Writer: start")
		//defer writer.Close()
		// the writer is connected to the reader via the pipe
		// so all data written here is passed on to the commands
		// standard input
		// writer.Write([]byte("p\n"))
		// writer.Write([]byte(" 10\n"))
		// writer.Write([]byte(" apple\n"))
		// writer.Write([]byte("c\n"))
		// writer.Write([]byte("g\n"))
		// writer.Write([]byte(" 10\n"))
		// writer.Write([]byte("c\n"))
		// writer.Write([]byte("Q\n"))
		// writer.Close()

		in.Write([]byte("p\n"))
		in.Write([]byte(" 10\n"))
		in.Write([]byte(" apple\n"))
		in.Write([]byte("p\n"))
		in.Write([]byte(" foobar\n"))
		in.Write([]byte(" butter\n"))
		in.Write([]byte("c\n"))
		in.Write([]byte("g\n"))
		in.Write([]byte(" 10\n"))
		in.Write([]byte("c\n"))
		in.Write([]byte("k\n"))
		in.Write([]byte("0 100\n"))
		in.Write([]byte("Q\n"))
		in.Close()
		fmt.Println("Writer: END")
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Wait: could not run command: ", err)
	}

}
