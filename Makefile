#Project root
ROOT=serverMultipleClient
client:
	go run $(ROOT)/src/cmd/clientBuild/client_main.go
master:
	go run $(ROOT)/src/cmd/serverBuild/server_main.go
peer_1:peer
peer_2:peer
peer_3:peer
peer_4:peer
peer:
	go run $(ROOT)/src/cmd/peerBuild/peer_main.go
