package peer

import (
	"utils"
		"encoding/gob"
			"net"
		"unsafe"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"os"
)

//Function which will be spawned in a different
//go routine and is responsible for maintaining/updating
//the backupExists peer with files of the primary peer
func UpdateBackupPeerStore(fileName string)  {
	//check if the file is updated with the
	//backupExists peer
	if !index[fileName] {
		//establish connection with the backupExists peer
		conn := establishConnectionBackup(peerNode.backupPeer)

		//create network encoders and decoders
		enc = gob.NewEncoder(conn)
		dec = gob.NewDecoder(conn)

		//fileV := client.MakeFileStruct(fileName, dirPath, peerIdentfier)
		//create a struct for the file
		f, err := os.Stat(filepath.Join(dirPath,peerIdentfier+fileName))
		utils.ValidateError(err)

		fileV := utils.File{
			Name:fileName,
			Size:f.Size(),
		}

		ok := sendProcessRequest(enc, dec, utils.STORE, fileV)

		//Read data from the file and write
		//to the network connection
		//read the file contents
		if ok{
			fileData, err := ioutil.ReadFile(filepath.Join(dirPath, f.Name()))
			utils.ValidateError(err)
			p := utils.CreatePacket(utils.DATA_END, string(fileData),
				unsafe.Sizeof(fileData)+unsafe.Sizeof(utils.DATA_END))
			err = enc.Encode(p)
			utils.ValidateError(err)

			fmt.Println(fileName + "  has been updated with the backup peer")
			conn.Close()
		}
	}
}

//Function which establishes connection with the
//backupExists peer specified in the peer structure
//params:
//	@backupPeer: string
//		Variable which holds the network address
//		of the backupExists peer
//Returns: conn net.Conn
//	Returns a pointer to the instance of the
//	TCP connection establishes with the backupExists peer
func establishConnectionBackup(networkAddr string) (conn net.Conn) {

	//dial the connection to the
	//backupExists based on the input parameter
	var err error
	fmt.Println("Establishing connection with backup: ", networkAddr)
	conn, err = net.Dial("tcp", networkAddr)
	utils.ValidateError(err)

	return
}

//Function which sends a request to the backupExists peer
//based on the type of the request specified in the
//parameter
//params:
//	@enc: gob.Encoder
//		Variables which holds the encoder defined
//		on the network connection defined for the
//		backupExists peer
//	@store: int
//		Variable which holds the integer variable
//		defining the type of request being sent
//	@fileV: utils.File
//		Variable which represents the utils.File
//		instance defined for the file being sent to
//		the backupExists peer
//Returns: ok bool
//			Returns true if the file prerequisites have
//			been updated with backupExists else false
func sendProcessRequest(enc *gob.Encoder, dec *gob.Decoder, store int, fileV utils.File) (ok bool) {
	//create the packet to send to the server
	totalSize := unsafe.Sizeof(store) + unsafe.Sizeof(string(fileV.Name))
	packet := utils.CreatePacket(store, string(fileV.Name), totalSize)
	packet.PfileInfo = fileV

	//send the packet
	gob.Register(utils.Packet{})
	gob.Register(utils.File{})
	err := enc.Encode(packet)
	utils.ValidateError(err)

	//Receive the confirmation packet
	//from the master and decode the peer
	//details
	var response utils.Response
	err = dec.Decode(&response)
	utils.ValidateError(err)

	if response.Ptype == utils.RESPONSE {
		ok = true
	}

	return
}
