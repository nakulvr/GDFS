package server

import (
	"time"
	"net"
			"utils"
	"encoding/gob"
	"unsafe"
	)

//Function which is spawned in a go routine
//when the main server is initialized and is 
//responsible to send heart beat messages to
//the peers, if any and update its peer maintainance
//structure, if required
func heartBeatHandler()  {
	doEvery(200*time.Millisecond, sendHeartBeat)
}

//Function which invokes the function which will
//be sending the heart beat signals to the peers. Additionally,
//this function will be invoking the function specified
//in the parameter for every 30 Miliseconds
//params:
//	@timeV: time.Duration
//		Specifies the time duration interval
//	@hBeatFunc: func()
//		Function which gets invoked every step
//		in the interval
func doEvery(timeV time.Duration, hBeatFunc func())  {
	for range time.Tick(timeV) {
		hBeatFunc()
	}
}

//Function which sends the heart beat
//to the peers updated in the master node
//structure
func sendHeartBeat()  {
	//check if there exists any primary
	//peers
	primary := masterNode.peers
	if len(primary) != 0{
		//fmt.Println("Reached here")
		for k := range primary{
			establishConnection(k, "", "primary")
			//if there exists a backup node for the
			//primary
			if b, ok := masterNode.backupPeers[k];ok{
				establishConnection(b, k, "backup")
			}
		}
	}
}

//Function which establishes a UDP connection
//with the peer and sends a heartbeat message
func establishConnection(networkAddr string, placeholer string, nodeType string)  {
	//dial a TCP connection to peer
	//if the connection is successful, the peer exists,
	//fmt.Println("establishing connection with "+nodeType)
	conn, err := net.Dial("tcp", networkAddr)

	//get the pointer to the primary peer and backup
	//peer map
	primary := masterNode.peers
	backup := masterNode.backupPeers

	//based on the error, perform the
	//structure modifications
	if err != nil{
		mutex.Lock()
		{
			//if the node is of type primary
			//then, add the delete the node
			//from the primary peer structure
			//and replace backup peer address as
			//primary
			if nodeType == "primary"{
				delete(primary, networkAddr)
				if _, ok := backup[networkAddr]; ok{
					backupAddr := backup[networkAddr]
					delete(backup, networkAddr)
					go updatePeer(backupAddr)
					primary[backupAddr] = 0
				}
			} else {
				//if  the node type is of backup
				//delete the backup peer from the
				//structure and reset the primary peer
				//backup status
				delete(backup, placeholer)
				go updatePeer(placeholer)
			}

		}
		mutex.Unlock()
	} else {
		conn.Close()
	}
}

//Function which establishes with the peer
//address in the parameter and sends a update
//signal stating that, it has no backup peer
//params:
//	@addr: string
//		Specifies the network address of the peer
//		to connect
func updatePeer(addr string)  {
	//dial a tcp connection
	conn, err := net.Dial("tcp", addr)
	utils.ValidateError(err)

	//create network encoder and decoder
	//to communicate
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	//create and send the packet of type
	//update
	pkt := utils.CreatePacket(utils.UPDATE, ",false", unsafe.Sizeof(utils.UPDATE))
	err = encoder.Encode(pkt)
	utils.ValidateError(err)

	//expect a packet of type response from
	//the peer
	err = decoder.Decode(&pkt)
	utils.ValidateError(err)
	if pkt.Ptype == utils.RESPONSE{
		conn.Close()
	}
}