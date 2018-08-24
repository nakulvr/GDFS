package utils

import (
		"log"
)

//global variable declaration
//Usage:
//Peer:
//	Ptype: Peer/backup
//	Pcontent: ""
//	Psize: size of Ptype
//Master:
//	Ptype: Response
//	Pcontent:""
//	Psize: sizeof Ptype
//Client:
//	Ptype: Fetch/Store
//	Pcontent: FileName/Data chunk
//	Psize:""/chunk size

type Packet struct {
	Ptype     int
	PfileInfo File
	Pcontent  string
	Psize     uintptr
}

type File struct{
	Name string
	Size int64
	PrimaryPeer string
	BackupPeer string
}

/*
Usage:
	Master:
		Ptype: Reponse
		Backup: True/False
		NetAddress: Backup peer addr
*/
type Response struct {
	Ptype      int
	Backup     bool
	NetAddress string
}

/*
Usage:
	Master:
		Ptype: Response
		PrimaryNetAddr: Primary Peer Network
						 address
		BackupNetAddr: Backup Peer Network
						address
*/
type ClientResponse struct {
	Ptype          int
	PrimaryNetAddr string
	BackupNetAddr  string
}

//constants which identify the
//the packet type
const (
	PEER     = 1
	FETCH    = 2
	STORE    = 3
	BACKUP   = 4
	RESPONSE = 5
	UPDATE = 6
	DATA = 7
	DATA_END = 8
)

/*
Function creates the packet based on
the input and returns the instace of the
packet
Params:
	p_type: Packet Type(PEER/
						FETCH/
						STORE/
						BACKUP/
						RESPONSE
	content: Content the packet will carry
	p_size: size of the packet

Returns: Instance of packet struct
*/
func CreatePacket(p_type int, content string, p_size uintptr) Packet {
	//create a packet with defined parameters
	packet_t := Packet{Ptype: p_type, Pcontent: content, Psize: p_size}

	//return the created packet
	return packet_t
}

/*
Function creates the Response packet based
on the input and returns the instance of the
response packet struct
Params:
	p_type: Packet Type(PEER/
						FETCH/
						STORE/
						BACKUP/
						RESPONSE
	content: Content the packet will carry
	p_size: size of the packet

Returns: Instance of packet struct
*/
func CreateResponse(p_type int, backup bool, parent string) Response {
	//create a packet with defined parameters
	//and return the created packet
	return Response{Ptype: p_type, Backup: backup, NetAddress: parent}
}

/*

 */
func CreateClientResponse(p_type int, primary string, backup string) ClientResponse {
	return ClientResponse{Ptype: p_type, PrimaryNetAddr: primary, BackupNetAddr: backup}
}

/*
Responsible for writing the err response
to the log.

Params:
	err: Error
 */
func ValidateError(err error)  {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err != nil{
		log.Fatal(err)
	}
}
