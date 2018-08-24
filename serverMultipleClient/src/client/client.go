package client

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"utils"
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

//Structures
type client struct {
	address    string
	port       int
	masterNode string
	//backupPeer    string
	//masterNode    string
}

type peerInfo struct {
	primaryPeer string
	backupPeer  string
}

// global variables
var (
	clientNode  client
	encode      *gob.Encoder
	decode      *gob.Decoder
	dirPath     string
	filePeerMap map[string]peerInfo
)

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

//Entry point to the client process. This function
//takes the master server IP address and port as the
//input and acts as an interface to the client process
//by initializing the client process.
//Params:
//	@remoteAddr: string
//		Takes the master server IP address in the
//		string format
//	@remotePort: string
//		Takes the master server port in the string
//		format
//Returns: nil
func Start(remoteAddr string, remotePort string) {

	//initialize the client instance
	initializeClient(remoteAddr, remotePort)

	//instantiate the command line interface
	//to the user
	initializeCLI()

}

// init the client
func initializeClient(remoteAddr string, remotePort string) {

	//form the network address for the node
	address := remoteAddr + ":" + remotePort
	clientNode = client{masterNode: address}

	//set the directory path from which the
	//client can read the files
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the Directory Path\n(Path where files that needs to be transferred are stored)")
	d, _, _ := r.ReadLine()
	dirPath = string(d)

	//dirPath = "C:\\Users\\mohan\\Desktop\\Courses\\Projects\\MDFS\\serverMultipleClient\\clientFiles"
	////dirPath = "Z:\\MS_NEU\\Courses\\CS\\Project\\MDFS\\serverMultipleClient\\clientFiles"
}

//Function which initializes the Command-line
//interface to the user, making the client features
//available to the user in terms of commmands.
func initializeCLI() {

	filePeerMap = make(map[string]peerInfo)
	//cli declarations
	clientIpAddr, err := utils.ExternalIP()
	utils.ValidateError(err)

	cliMessage := "client@" + clientIpAddr + ">>"
	reader := bufio.NewReader(os.Stdin)

	//in an infinite loop
	for {
		fmt.Print(cliMessage)

		//read the input command.
		command, err := reader.ReadString('\n')
		utils.ValidateError(err)

		//process and validate the input command
		processAndValidate(command)
	}
}

//Function which processes the input command
//and validates it against the valid options.
func processAndValidate(command string) {
	//Step-1: Remove unexpected suffixes
	command = strings.TrimSuffix(command, "\n")

	//Step-2: Split the string into tokens
	tokens := strings.Split(command, " ")

	switch tokens[0] {
	case "send":
		conn := establishConnection()

		fileV := MakeFileStruct(tokens[1], dirPath, "")

		defer conn.Close()
		sendFile(conn, fileV)
		break
	case "receive":
		primary := filePeerMap[tokens[1]].primaryPeer
		backup := filePeerMap[tokens[1]].backupPeer


		if primary != "" {
			fmt.Printf("Primary %s, Secondary %s\n", primary, backup)
			data := fetchDataFromPeer(tokens[1], filePeerMap[tokens[1]].primaryPeer, filePeerMap[tokens[1]].backupPeer)
			o := writeToFile(data, tokens[1], dirPath)
			if o{
				//fmt.Println("File Successfully written to the disk")
				fmt.Printf("File Successfully received, data:\n%s\n", data)
			}
		} else {
			fmt.Println("Error retrieving file. System Failure. No Peers standing!!!")
		}
		break
	case "quit": os.Exit(0)
	}
}

func fetchDataFromPeer(fileName string, primaryPeer string, backupPeer string) string {
	conn, err := net.DialTimeout("tcp", primaryPeer, time.Duration(time.Duration(200*time.Millisecond)))
	//utils.ValidateError(err)
	if err != nil {
		fmt.Printf("Primary Peer down Fetching from Backup\n")
		conn, err = net.DialTimeout("tcp", backupPeer, time.Duration(time.Duration(200*time.Millisecond)))
		utils.ValidateError(err)
	}
	totalSize := unsafe.Sizeof(utils.FETCH) + unsafe.Sizeof(fileName)
	packet := utils.CreatePacket(utils.FETCH, fileName, totalSize)
	gob.Register(utils.Packet{})
	//gob.Register(utils.File{})
	encode = gob.NewEncoder(conn)
	decode = gob.NewDecoder(conn)

	err = encode.Encode(packet)
	utils.ValidateError(err)

	var resp utils.Packet
	err = decode.Decode(&resp)
	utils.ValidateError(err)
	conn.Close()
	return resp.Pcontent
}

//Function which writes to the file retrieved from the
//peers
//params:
//	@data: string
//		Consists of data converted from bytes into
//		string format
//	@fileName: string
//		fileName of the file to be written.
//	@dir: string
//		directory path where the files is to be
//		created
//Returns: ok bool
//		If the file is successfully written to
//		the disk, returns true, else false
func writeToFile(data string, fileName string, dir string) (ok bool) {
	//open the file if exists or create one with
	//that name
	fileLoc := filepath.Join(dirPath, fileName)
	err := ioutil.WriteFile(fileLoc, []byte(data), 0755)
	utils.ValidateError(err)
	ok = true
	return
}

//Function which establishes a connection
//to the master server on demand
func establishConnection() (conn net.Conn) {
	//dial a TCP connection to the master node/server
	conn, err := net.Dial("tcp", clientNode.masterNode)
	utils.ValidateError(err)

	//initialize the client instance with its
	//system local address and port
	if clientNode.address == "" && clientNode.port == 0 {
		//get the local address from the connection
		networkAddr := conn.LocalAddr().String()
		addr := strings.Split(networkAddr, ":")

		//initialize the processed values
		clientNode.address = addr[0]
		clientNode.port, err = strconv.Atoi(addr[1])
		utils.ValidateError(err)
	}

	return
}

//Function which sends the file to the
//server established on the conn instance passed
//as a parameter
//Params:
//	@conn: net.Conn
//		Instance which holds the TCP connection
//		to the master server
//	@fileV: file
//		An instance of the file structure
//		which holds the file name, the address
//		of the primary and backup peers where the
//		file is present at, in a string format.
//Returns: nil
func sendFile(conn net.Conn, fileV utils.File) {
	var err error

	//create the packet to send to the server
	totalSize := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(fileV.Name))
	packet := utils.CreatePacket(utils.STORE, string(fileV.Name), totalSize)
	packet.PfileInfo = fileV

	//println("Reached here")
	//send the packet
	gob.Register(utils.Packet{})
	gob.Register(utils.File{})
	encode = gob.NewEncoder(conn)
	err = encode.Encode(packet)
	utils.ValidateError(err)

	//Receive the confirmation packet
	//from the master and decode the peer
	//details
	var response utils.ClientResponse
	decode = gob.NewDecoder(conn)
	err = decode.Decode(&response)
	//println("Reached here too")
	utils.ValidateError(err)

	if response.Ptype == utils.RESPONSE {
		if response.PrimaryNetAddr != "" {
			fmt.Printf("Primary: %s, Secondary: %s\n", response.PrimaryNetAddr, response.BackupNetAddr)
			fileV.PrimaryPeer = response.PrimaryNetAddr
			fileV.BackupPeer = response.BackupNetAddr
			filePeerMap[fileV.Name] = peerInfo{response.PrimaryNetAddr, response.BackupNetAddr}

			//once the primary and backup peer
			//credentials have been established,
			//contact the primary and send the file.
			sendData(fileV)
			fmt.Println("Data Successfully sent to: "+response.PrimaryNetAddr)
		} else {
			fmt.Println("Error storing the file. System Failure. No Peers standing!!!")
		}
	}


}

//Function which establishes connection
//with the primary peer address received
//from the master node and sends the file
func sendData(fileV utils.File) {

	//read the file contents
	fileData, err := ioutil.ReadFile(filepath.Join(dirPath, fileV.Name))
	utils.ValidateError(err)

	//establish connection with the
	//primary peer
	conn, err := net.Dial("tcp", fileV.PrimaryPeer)
	utils.ValidateError(err)

	//create a network encoder and deocder
	encode = gob.NewEncoder(conn)
	decode = gob.NewDecoder(conn)

	//send a STORE request to the primary peer
	totalSize := unsafe.Sizeof(utils.STORE) + unsafe.Sizeof(string(fileData))
	packet := utils.CreatePacket(utils.STORE, string(fileData), totalSize)

	//register the interface with the gob
	gob.Register(utils.Packet{})
	packet.PfileInfo = fileV
	utils.ValidateError(err)

	//once the packet is ready,
	//send a STORE request to the primary
	err = encode.Encode(packet)

	//Expect a RESPONSE from the primary
	//confirming the client that the necessary
	//setups are done and file can now be sent.
	var resp utils.Packet
	err = decode.Decode(&resp)
	utils.ValidateError(err)
	//validate the packet received. If
	//it is of the type response, send
	//the data on the same established connection
	if resp.Ptype == utils.RESPONSE {
		totalSize := unsafe.Sizeof(utils.DATA_END) + unsafe.Sizeof(string(fileData))
		packet := utils.CreatePacket(utils.DATA_END, string(fileData), totalSize)
		err = encode.Encode(packet)
		utils.ValidateError(err)
	}
}

//Function which takes the token as input
//and returns a file struct as output.
//Params:
//	@fileName: string
//		A variable which holds the value of
//		file name processed from the command
//Returns: fileV utils.File
func MakeFileStruct(fileName string, d string, suffix string) (fileV utils.File) {
	//create a struct for the file
	f, err := os.Stat(filepath.Join(d, suffix+fileName))

	utils.ValidateError(err)

	fileV = utils.File{
		Name: f.Name(),
		Size: f.Size(),
	}

	return
}
