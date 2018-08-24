package main

import (
	"peer"
	"flag"
)

func main() {
	//parse the command line arguments
	remoteAddr := flag.String("addr", "127.0.0.1", "The address of the serverBuild to connect to."+
		"Default is localhost")

	remotePort := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	peer.Start(*remoteAddr, *remotePort)
}
