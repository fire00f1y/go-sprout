package main

import (
	"context"
	gosprout "github.com/fire00f1y/go-sprout"
	"github.com/fire00f1y/go-sprout/resource/file"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	fileName := "temp.file"
	f, e := os.Create(fileName)
	if e != nil {
		log.Fatalf("error creating file: %v\n", e)
	}
	res := file.NewResource(fileName)

	// Remove file after the sample
	defer func() {
		if e := os.Remove(fileName); e != nil {
			log.Print(e)
		}
	}()

	// Close file
	defer func() {
		log.Print(f.Close())
	}()

	_, e = f.WriteString("")
	if e != nil {
		log.Fatalf("could not write to file: %v\n", e)
	}

	done := make(chan string)
	gosprout.Watch(context.Background(),
		500*time.Millisecond,
		res,
		func(r io.Reader) {
			s, _ := ioutil.ReadAll(r)
			done <- string(s)
		},
		func(e error) {
			log.Printf("error during update: %v\n", e)
		},
	)

	text := "lol"
	log.Printf("writing %s to the file\n", text)
	_, e = f.WriteString(text)
	if e != nil {
		log.Printf("failed to write to file: %v\n", e)
		return
	}

	d := <-done
	log.Printf("got update from update function: %s\n", d)
}
