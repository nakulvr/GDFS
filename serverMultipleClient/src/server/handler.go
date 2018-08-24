package server

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"utils"
)


func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}

/*
Function which handles the incoming
peerBuild requests to the serverBuild.
It performs any necessary action and/or invokes
other functions to complete the tasks

Returns: nil
*/
func handleConnection(conn net.Conn) {

	//Receive and Decode the packet on the
	//network.
	var recv utils.Packet
	err := dec.Decode(&recv)
	if err != nil {
		print("Error while decoding packet: ", err.Error())
	}

	//parse the packet
	if recv.Ptype == utils.PEER {
		validatePeer(conn, recv)
	}

	if recv.Ptype == utils.STORE {
		//var response utils.ClientResponse
		// TODO: selecting the peers according to hash

		//var primary string
		//for ip, _ := range masterNode.peers {
		//	primary = ip
		//	//fmt.Printf("IP: %s, Port: %d", ip, port)
		//}
		var response utils.ClientResponse
		if len(masterNode.peers) > 0{
			primary := MapRandomKeyGet(masterNode.peers)
			backupPeer, _ := masterNode.backupPeers[primary.(string)]
			response = utils.CreateClientResponse(utils.RESPONSE, primary.(string), backupPeer)
		} else {
			response = utils.CreateClientResponse(utils.RESPONSE, "", "")
		}

		enc = gob.NewEncoder(conn)
		err := enc.Encode(response)
		utils.ValidateError(err)
		fmt.Println("Response sent to the client")
	}

	//close the connection
	conn.Close()
}

/*
Function which handles the incoming request
from a peer
*/
func validatePeer(conn net.Conn, recv utils.Packet) {
	//debug
	n := len(masterNode.peers)
	b_n := len(masterNode.backupPeers)

	//get the address of the tcp-peerBuild
	clientAddr := conn.RemoteAddr().String()

	//add the peerBuild to the peer list
	networkAddr := strings.Split(clientAddr, ":")
	clientPort, err := strconv.Atoi(networkAddr[1])
	if err != nil {
		fmt.Printf("Conversion Error: %s", err.Error())
	}

	//if every peer registered has a backup
	if n == b_n {
		mutex.Lock()
		_, ok := masterNode.peers[clientAddr]
		if !ok {
			masterNode.peers[clientAddr] = clientPort
		} else {
			fmt.Println("Peer already registered. If the peer needs to update " +
				"the details, send an UPDATE message to the master")
		}
		mutex.Unlock()

		//send response packet
		p := utils.Response{Ptype: utils.RESPONSE, Backup: false, NetAddress: ""}
		err = enc.Encode(p)
		if err != nil {
			print("Error while encoding peer packet: ", err.Error())
		}

		//check if the peer has been added
		if n != len(masterNode.peers) {
			fmt.Println("Peer added")
			previousPeer = clientAddr
		} else {
			fmt.Println("Peer registration unsuccessful.")
		}
	} else if n > b_n { //when the recently added peer does
		//not have a backup
		mutex.Lock()
		//masterNode.peers[clientAddr] = clientPort
		masterNode.backupPeers[previousPeer] = clientAddr
		mutex.Unlock()

		//createResponse Packet
		resp := utils.CreateResponse(utils.RESPONSE, true, previousPeer)
		err = enc.Encode(resp)
		if err != nil {
			print("Error while encoding response packet: ", err.Error())
		}

		//check if the peer has been added
		if b_n != len(masterNode.backupPeers) {
			fmt.Println("Backup Peer added")
		} else {
			fmt.Println("Peer registration unsuccessful.")
		}
	}
}
