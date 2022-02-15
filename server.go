package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

var addr string
var replayport int

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	mulu := "./http" + r.URL.Path
	//r.RequestURI

	b, err := ioutil.ReadFile(mulu)
	if checkError(err) {
		fmt.Fprint(w, "")
		return
	}
	fmt.Fprint(w, string(b))
	//fmt.Fprintln(w, "hello world")
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("Error:", err.Error())
		return true
	}
	return false
}
func udp(address string) {

	//fmt.Print("OK")
	udpAddr, _ := net.ResolveUDPAddr("udp", address)
	UDPlistener, err := net.ListenUDP("udp", udpAddr)

	//UDPlistener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0})
	if checkError(err) {
		return
	}
	fmt.Println("Using Local UDP port:", UDPlistener.LocalAddr().(*net.UDPAddr).Port)
	//var replayport int
	//
	//mihoyoAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:22223")
	mihoyoAddr, _ := net.ResolveUDPAddr("udp", "47.103.25.145:22102")
	for {
		data2 := make([]byte, 2048)
		//读取
		n, addr, err := UDPlistener.ReadFromUDP(data2)
		if checkError(err) {
			continue
		}
		//fmt.Printf("lenth:%d\n", n)
		//fmt.Printf("Port:%d,Prot2:%d\n", addr.Port, udpAddr.Port)
		//fmt.Println(addr.IP.String()[0:3])
		//var replay *net.UDPAddr
		//replay, _ := net.ResolveUDPAddr("udp", "0.0.0.0:0")

		//if addr.Port != udpAddr.Port {
		//if addr.Port != replay.Port {
		if addr.Port != mihoyoAddr.Port {
			replayport = addr.Port
			//replay = addr
			//fmt.Println(replay)
			//data2 = BytesCombine([]byte{1, uint8(len(pss))}, []byte(pss), tempport, data2)
			//fmt.Println(data2)
			//data2 = xor(data2)
			//fmt.Println("Send Data")
			_, err = UDPlistener.WriteToUDP(data2[:n], mihoyoAddr)
			//_, err = UDPlistener.WriteToUDP(data2[:n], udpAddr)
			if checkError(err) {
				continue
			}

		} else {
			//fmt.Println("Reveive Data")
			fmt.Println(replayport)
			//data2 = xor(data2)
			//fmt.Println(data2)
			//_, err = UDPlistener.WriteToUDP(data2[:n], mihoyoAddr)
			//_, err = UDPlistener.WriteToUDP(data2[:n], replay)
			_, err = UDPlistener.WriteToUDP(data2[:n], &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: replayport})
			if checkError(err) {
				continue
			}
		}
	}

}
func httpreq(addr string) {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(addr, nil)
}
func main() {
	addr := "0.0.0.0:22102"
	http := "0.0.0.0:22220"
	//serverdd = "127.0.0.1:9997"
	//pss = "000"
	//fmt.Println(os.Args)
	flag.StringVar(&addr, "l", "0.0.0.0:22222", "Listen")
	flag.Parse()
	go httpreq(http)
	udp(addr)
	fmt.Scanf("%s", &addr)
	//fmt.Print("123")
}
