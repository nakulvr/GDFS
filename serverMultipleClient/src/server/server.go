package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
//global variable declaration
type master struct {
	address     string
	port        int
	networkAddr string
	peers       map[string]int
	backupPeers map[string]string
}

var (
	masterNode master
	enc *gob.Encoder
	dec *gob.Decoder
	mutex = &sync.Mutex{}
	previousPeer string
	)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/*
Function which initializes the master struct with initial
values which are parsed from the command line.

Returns: nil
*/
func initializeMaster(server string, port string) {

	//form the network address for the node
	address := server + ":" + port

	//initialize the global variable
	//representing master node
	p, err := strconv.Atoi(port)
	if err != nil {
		fmt.Printf("Conversion Error: %s", err.Error())
	}

	masterNode = master{address: server, port: p, networkAddr: address,
		peers: make(map[string]int), backupPeers: make(map[string]string)}

}

/*
Function which listens on the dedicated port specified
for the master node and accepts the clients requests only
to pass it to a go routine to handle the requests

Returns: nil
*/
func acceptAndProcess(node master) {

	//heart beat signal handler
	go heartBeatHandler()

	//listen on the designates network address
	adapter, err := net.Listen("tcp", node.networkAddr)
	if err != nil {
		fmt.Printf("Error while listening to the on port: %d", node.port)
		return
	}

	//until a SIGNAL interrupt is passed or an exception is
	//raised, keep on accepting peerBuild connections and add it
	//to the peer map.
	for {

		//debug information
		fmt.Printf("\nListening on Port: %d\n", node.port)

		//accept incoming connections
		conn, err := adapter.Accept()
		if err != nil {
			println(err.Error())
			continue
		}

		// Will write to network.
		enc = gob.NewEncoder(conn)
		// Will read from network.
		dec = gob.NewDecoder(conn)

		//start a go routine to handle
		//the incoming connections
		go handleConnection(conn)
	}
}

/*
Function which returns the address of master struct
for testing functions
*/
func StructAddr() *master {
	return &masterNode
}

/*
Function which starts the serverBuild and
passes the initialized values for listening
on the designated port.
*/
func StartServer(server string, port string) {

	//initialize the structure to define the
	//master node
	initializeMaster(server, port)

	//once the master node is initialized,
	//listen on the dedicated port and accept
	//connections
	acceptAndProcess(masterNode)

}
