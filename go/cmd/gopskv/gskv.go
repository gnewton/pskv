package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"sync"
)

type Store struct {
	cmd *exec.Cmd

	in  io.WriteCloser
	out io.ReadCloser
	err io.ReadCloser

	outScanner *bufio.Scanner

	mu sync.Mutex
}

func Open() (*Store, error) {
	s := new(Store)
	s.cmd = exec.Command("gs", "-sDEVICE=nullpage", "-q", "-dNOPAUSE", "-dBATCH", "kv.ps")

	var err error
	if s.in, err = s.cmd.StdinPipe(); err != nil {
		return nil, err
	}

	if s.out, err = s.cmd.StdoutPipe(); err != nil {
		return nil, err
	}

	s.outScanner = bufio.NewScanner(s.out)

	if s.err, err = s.cmd.StderrPipe(); err != nil {
		return nil, err
	}

	if err = s.cmd.Start(); err != nil {
		return nil, err
	}

	return s, nil
}

var PutCommand = []byte("p\n")
var GetCommand = []byte("g\n")
var QuitCommand = []byte("Q\n")
var CountCommand = []byte("c\n")
var DeleteAllCommand = []byte("D\n")

func (kv *Store) Put(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, err := kv.in.Write(PutCommand); err != nil {
		kv.close()
		return err
	}

	if _, err := kv.in.Write([]byte(" " + key + "\n")); err != nil {
		kv.close()
		return err
	}

	if _, err := kv.in.Write([]byte(" " + value + "\n")); err != nil {
		kv.close()
		return err
	}
	return nil
}

func (kv *Store) Get(key string) (bool, string, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, err := kv.in.Write(GetCommand); err != nil {
		kv.close()
		return false, "", err
	}

	if _, err := kv.in.Write([]byte(" " + key + "\n")); err != nil {
		kv.close()
		return false, "", err
	}

	if kv.outScanner.Scan() {
		found := kv.outScanner.Text()
		switch found {
		case "f":
			return false, "", nil
		case "t":
			if kv.outScanner.Scan() {
				return true, kv.outScanner.Text(), nil
			} else {
				kv.close()
				//error: FIX
				return false, "", nil
			}
		}
	} else {
		kv.close()
		//error: FIX
		return false, "", nil
	}
	kv.close()
	//error: FIX
	return false, "", nil
}

func (kv *Store) Count() (uint64, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, err := kv.in.Write(CountCommand); err != nil {
		kv.close()
		return 0, err
	}

	if kv.outScanner.Scan() {
		countStr := kv.outScanner.Text()
		return strconv.ParseUint(countStr, 10, 0)
	} else {
		kv.close()
		return 0, errors.New("Error: count: kv error")
	}
}

func (kv *Store) DeleteAll() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, err := kv.in.Write(DeleteAllCommand); err != nil {
		kv.close()
		return err
	}
	return nil
}

func (kv *Store) Delete(key string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return nil
}

func (kv *Store) ListKeys(start, end uint64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return nil
}

func (kv *Store) close() error {
	if err := kv.in.Close(); err != nil {
		log.Println(err)
		kv.out.Close()
		return err
	}

	if err := kv.out.Close(); err != nil {
		return err
	}
	return nil
}

func (kv *Store) Close() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	var err error
	if _, err := kv.in.Write(QuitCommand); err != nil {
		log.Println(err)
	}

	err2 := kv.close()
	if err2 == nil {
		return err
	}
	return err2
}
