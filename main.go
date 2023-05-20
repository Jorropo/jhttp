package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := flag.String("a", "localhost:8081", "Address to serve at.")
	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintln(o, `Usage: jhttp [path/to/serve]`)
		flag.PrintDefaults()
	}
	flag.Parse()

	var path string
	switch args := flag.Args(); len(args) {
	case 0:
		path = "."
	case 1:
		path = args[0]
	}

	s := http.Server{
		Handler: http.FileServer(http.Dir(path)),
		Addr:    *addr,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Closing")
		s.Close()
		os.Exit(0)
	}()

	fmt.Printf("will serve on http://%s\n", *addr)
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
