package main

import (
	"client"
	"flag"
)

func main() {

	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the Master to connect to."+
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port of the Master daemon.")

	flag.Parse()

	//start the client process
	client.Start(*remoteAddr, *remotePort)
}
