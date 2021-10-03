package main

import (
	"crypto/rand"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	fWorkDir = flag.String("work-dir", "", "working directory to save bookings")
	fAddr    = flag.String("addr", ":8090", "listen addr")
)

//go:embed index.html
var html string

func booking(w http.ResponseWriter, req *http.Request) {
	err := func() error {
		file, err := ioutil.TempFile(*fWorkDir, "*.bk")
		if err != nil {
			return err
		}
		// err = file.Close()
		// if err != nil {
		// 	return err
		// }
		// file, err = os.OpenFile(file.Name(), os.O_RDWR|os.O_SYNC, 0600)
		// if err != nil {
		// 	return err
		// }
		defer file.Close()
		b := make([]byte, 1111)
		rand.Read(b)
		_, err = file.Write(b)
		if err != nil {
			return err
		}
		if err := file.Sync(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "fail to book: %s\n", err.Error())
		return
	}
	fmt.Fprintf(w, "[%s] booking succeeded\n", time.Now().Format(time.StampMilli))
}

// func save(file string) error

func main() {
	flag.Parse()

	http.HandleFunc("/booking", booking)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", html)
	})

	srv := &http.Server{
		Addr: *fAddr,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 5 * time.Second,
	}
	log.Println(srv.ListenAndServe())
}
