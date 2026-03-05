package tun

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"regexp"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

func CreateTunInterface(name string, cidr string) *water.Interface {
	// Create TUN interface
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = name

	tun, err := water.New(config)
	if err != nil {
		log.Fatalf("Error creating TUN interface: %v", err)
	}
	log.Printf("[tun] TUN interface %s created successfully.", tun.Name())

	// Configure the TUN interface with the specified CIDR
	re := regexp.MustCompile(`(\d+\.\d+\.\d+)\.\d+/(\d+)`)
	ipAddress := re.ReplaceAllString(cidr, "$1.1/$2")

	// Use netlink to set the IP address and bring up the interface
	link, err := netlink.LinkByName(name)
	if err != nil {
		log.Fatalf("Error finding TUN interface: %v", err)
	}
	addr, _ := netlink.ParseAddr(ipAddress)
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Fatalf("Error adding IP address to TUN interface: %v", err)
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		log.Fatalf("Error bringing up TUN interface: %v", err)
	}

	return tun
}

func StartTunServer(tun *water.Interface) {
	go tunToConn(tun)

	// Listen for incoming TCP connections on port 443
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[tun] Server listening on port 443.")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		log.Println("[tun] Client connected.")
		go connToTun(conn, tun)
	}
}

var clientConn net.Conn

func connToTun(conn net.Conn, tun *water.Interface) {
	clientConn = conn
	for {
		lengthBuf := make([]byte, 2)
		_, err := io.ReadFull(conn, lengthBuf)
		if err != nil {
			return
		}

		length := binary.BigEndian.Uint16(lengthBuf)
		packet := make([]byte, length)

		_, err = io.ReadFull(conn, packet)
		if err != nil {
			return
		}

		tun.Write(packet)
	}
}

func tunToConn(tun *water.Interface) {
	buf := make([]byte, 2000)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			continue
		}

		if clientConn != nil {
			lengthBuf := make([]byte, 2)
			binary.BigEndian.PutUint16(lengthBuf, uint16(n))
			clientConn.Write(lengthBuf)
			clientConn.Write(buf[:n])
		}
	}
}
