package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/libdns/libdns"
	"github.com/libdns/loopia"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	user := os.Getenv("LOOPIA_USER")
	password := os.Getenv("LOOPIA_PASSWORD")
	var override, zone string
	if len(os.Args) > 1 {
		zone = os.Args[1]
	}

	if len(os.Args) > 2 {
		override = os.Args[2]
	}

	if zone == "" {
		fmt.Fprintf(os.Stderr, "ZONE not set\n")
		os.Exit(1)
	}

	if user == "" {
		exitOnError(fmt.Errorf("user is not set"))
	}

	if password == "" {
		exitOnError(fmt.Errorf("password is not set"))
	}

	fmt.Printf("zone: %s, user: %s\n", zone, user)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	go show(ctx, &wg, zone, user, password, override)

	// Wait for SIGINT.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		cancel()
	}()

	wg.Wait()
	fmt.Println("Done!")
}

func show(ctx context.Context, wg *sync.WaitGroup, zone, user, password, override string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	p := &loopia.Provider{
		Username:       user,
		Password:       password,
		OverrideDomain: override,
	}
	fmt.Println("getting records")
	resAll, err := p.GetRecords(ctx, zone)
	exitOnError(err)
	printRecords("All records", resAll)
	wg.Done()
}

func printRecords(title string, records []libdns.Record) {
	fmt.Println(title)
	for i, r := range records {
		fmt.Printf("  [%d] %+v\n", i, r)
	}
}
