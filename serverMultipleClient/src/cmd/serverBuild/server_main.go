package main

import (
	"flag"
	"server"
)

func main() {
	//parse the command line arguments
	serverV := flag.String("serverAddr", "", "The address of the serverBuild to connect to."+
		"Default is localhost")

	port := flag.String("port", "9999", "Port to listen for incoming connections.")

	flag.Parse()

	server.StartServer(*serverV, *port)
}
