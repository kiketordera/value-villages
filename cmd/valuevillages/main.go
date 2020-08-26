package main

import (
	"sync"

	"github.com/kiketordera/value-villages/internal/app/village"
)

var (
	debug = false
)

func main() {
	// Initialize the app
	v := village.New()

	// New waitgroup for sync
	var wg sync.WaitGroup

	// We inicialize 1 app, so we wait 1 process
	wg.Add(1)
	//I nit the app (iun the background due to GO)
	go v.Start(&wg)

	// Wait until all apps stop
	wg.Wait()
	v.DataBase.Close()
}
