package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
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
		defer file.Close()
		_, err = file.WriteString("this is fake booking content\n")
		if err != nil {
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

func main() {
	flag.Parse()

	http.HandleFunc("/booking", booking)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", html)
	})

	http.ListenAndServe(*fAddr, nil)
}
