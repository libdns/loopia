package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	zone := os.Getenv("ZONE")
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
	go show(ctx, &wg, zone, user, password)

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

func show(ctx context.Context, wg *sync.WaitGroup, zone, user, password string) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	p := &loopia.Provider{
		Username: user,
		Password: password,
	}
	fmt.Println("getting records")
	resAll, err := p.GetRecords(ctx, zone)
	exitOnError(err)
	printRecords("All records", resAll)
	toDelete := []libdns.Record{}
	for _, r := range resAll {
		rr := r.RR()
		if rr.Type == "TXT" && strings.Contains(rr.Name, "_acme-challenge") {
			toDelete = append(toDelete, r)
		}
	}

	if len(toDelete) > 0 {
		fmt.Println("deleting records")
		res, err := p.DeleteRecords(ctx, zone, toDelete)
		if err != nil {
			fmt.Printf("  error deleting %s\n", err)
		} else {
			printRecords("Deleted records", res)
		}
	}

	wg.Done()
}

func printRecords(title string, records []libdns.Record) {
	fmt.Println(title)
	for i, r := range records {
		fmt.Printf("  [%d] %+v\n", i, r)
	}
}
