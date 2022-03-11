package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
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

var waitFlag = flag.Int("wait", 20, "time to wait in seconds")

func main() {
	flag.Parse()
	user := os.Getenv("LOOPIA_USER")
	password := os.Getenv("LOOPIA_PASSWORD")
	zone := os.Getenv("ZONE")

	if flag.NArg() > 0 {
		zone = flag.Arg(0)
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

	z := zone
	components := strings.Split(zone, ".")
	if components[0] == "*" {
		z = strings.Join(components[1:], ".")
	}

	p := &loopia.Provider{
		Username: user,
		Password: password,
	}
	ctx := context.TODO()
	fmt.Println("appending")
	res, err := p.AppendRecords(ctx, z,
		[]libdns.Record{
			{Name: "_acme-challenge.test", Type: "TXT", Value: "Zgu7tw287LB-LpXyTHYLeROag9-4CLHnM77zvTEvH6o"},
		})
	exitOnError(err)
	printRecords("after append", res)

	fmt.Printf("Will sleep for %d seconds...\n", *waitFlag)
	time.Sleep(time.Second * time.Duration(*waitFlag))

	resAll, err := p.GetRecords(ctx, z)
	exitOnError(err)
	printRecords("after all", resAll)

	// delete first
	res, err = p.DeleteRecords(ctx, z, res)
	exitOnError(err)
	printRecords("after delete", res)

	// check final result
	resAll, err = p.GetRecords(ctx, z)
	exitOnError(err)
	printRecords("after all", resAll)

	fmt.Println("Done!")
}

func printRecords(title string, records []libdns.Record) {
	fmt.Println(title)
	for i, r := range records {
		fmt.Printf("  [%d] %+v\n", i, r)
	}
}
