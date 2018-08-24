package tests

import (
	"testing"
	"server"
	"time"
	"peer"
)

func TestServer(t *testing.T)  {
	serverTest1(t)
}

func TestServer2(t *testing.T)  {
	serverTest2(t)
}

/*
Function which tests for the values
after the serverBuild has started
 */
func serverTest1(t *testing.T)  {
	//Start the serverBuild
	go server.StartServer("127.0.0.1", "9999")

	time.Sleep(time.Second * 2)

	go peer.Start("127.0.0.1", "9999")

	time.Sleep(time.Second * 5)
}

func serverTest2(t *testing.T)  {
	//Start the serverBuild
	go server.StartServer("127.0.0.1", "9999")

	time.Sleep(time.Second * 2)

	go peer.Start("127.0.0.1", "9999")

	time.Sleep(time.Second * 5)
}