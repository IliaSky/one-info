package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/nerdgguy/go-imap"
	"os"
)

func usage() {
	fmt.Print("Only works on Gmail accounts\n")
	fmt.Print("Make sure IMAP is enabled in your Gmail account.")
	fmt.Printf("usage: %s email password\n", os.Args[0])
	os.Exit(0)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}
	user := args[0]
	pass := args[1]

	conn, err := tls.Dial("tcp", "imap.gmail.com:993", nil)
	if err != nil {
		fmt.Print("Connection to imap.gmail.com:993 failed")
		os.Exit(1)
	}

	im := imap.New(conn, conn)
	//im.Unsolicited = make(chan interface{}, 100)

	hello, err := im.Start()
	if err != nil {
		fmt.Print("Hello failed")
		os.Exit(1)
	}
	fmt.Printf("Server hello: %s\n", hello)

	fmt.Printf("logging in...\n")
	resp, caps, err := im.Auth(user, pass)
	if err != nil {
		fmt.Print("Auth failed\n")
		os.Exit(1)
	}
	fmt.Printf("Server capabilities: %s\n", caps)
	fmt.Printf("Server resp: %s\n", resp)

	examine, err := im.Examine("INBOX")
	if err != nil {
		fmt.Print("Examine failed\n")
		os.Exit(1)
	}
	fmt.Printf("Mailbox status: %+v\n", examine)

	fresp, err := im.Fetch("1:*", []string{"RFC822"})
	if err != nil {
		fmt.Print("Fetch failed\n")
		os.Exit(1)
	}
	for _, fr := range fresp {
		fmt.Printf("Fetch resp: %v\n", string(fr.Rfc822))
	}
}
