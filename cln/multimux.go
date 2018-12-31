package main

import (
	"io"
	"log"
	"net"

	"github.com/inconshreveable/muxado"
)

type MultiMuxClient struct {
	connections []muxado.Session
	sconns      chan net.Conn
}

func NewMultiMuxClient1() *MultiMuxClient {
	mmClient := &MultiMuxClient{}
	mmClient.sconns = make(chan net.Conn, 16)

	return mmClient
}

func NewMultiMuxClient(c io.ReadWriteCloser) *MultiMuxClient {
	mmClient := &MultiMuxClient{}
	mmClient.sconns = make(chan net.Conn, 16)

	sess := muxado.Client(c, nil)
	go mmClient.handleSession(sess)

	return mmClient
}

func (m *MultiMuxClient) AddConnection(c io.ReadWriteCloser) {
	sess := muxado.Client(c, nil)
	go m.handleSession(sess)
}

func (m *MultiMuxClient) handleSession(sess muxado.Session) {
	for {
		sconn, err := sess.Accept()
		if err != nil {
			log.Println("Can't accept, connection is dead", err)
			break
		}
		m.sconns <- sconn
	}
}
func (m *MultiMuxClient) Accept() (net.Conn, error) {
	// for _, sess := range m.connections {
	// 	go m.handleSession(sess)
	// }
	sconn := <-m.sconns

	return sconn, nil
}