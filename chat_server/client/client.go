package main

import (
	"../goconfig"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

func main() {
	host, port := readServerConfig()
	addr := host + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	fmt.Println("connecting ", conn.RemoteAddr().String(), "...")
	for {
		timeNow := time.Now().String()
		_, err := conn.Write([]byte(timeNow))
		checkError(err)
		buf := make([]byte, 256)
		conn.Read(buf)
		fmt.Println("recv : ", string(buf))
		time.Sleep(time.Second * 2)
	}
	conn.Close()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func readServerConfig() (string, string) {
	exefile, _ := exec.LookPath(os.Args[0])
	fmt.Println(filepath.Dir(exefile))

	filepath := path.Join(filepath.Dir(exefile), "./client_config.ini")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err)

	host, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.host")
	checkError(err)

	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)

	return host, port
}
