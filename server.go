package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/xtaci/kcp-go"
)

var addr string
var replayport int

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	//if r.URL.Path == "query_cur_region" {
	//	query_cur_region(w, r)
	//}
	//if r.URL.Path == "query_region_list" {
	//	query_region_list(w, r)
	//}

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
func handleEcho(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(buf)
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func client(test string) {
	//key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	//block, _ := kcp.NewAESBlockCrypt(key)

	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions(test, nil, 0, 0); err == nil {
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			log.Println("sent:", data)
			if _, err := sess.Write([]byte(data)); err == nil {
				// read back the data
				if _, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", string(buf))
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}

func udp(test string) {
	//key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	//block, _ := kcp.NewAESBlockCrypt(key)
	fmt.Println(test)
	if listener, err := kcp.ListenWithOptions(test, nil, 0, 0); err == nil {
		// spin-up the client
		//go client(test)
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(s)
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}

}

func udpc(address string) {

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
	http := "0.0.0.0:22223"
	//serverdd = "127.0.0.1:9997"
	//pss = "000"
	//test()
	//fmt.Println(os.Args)
	flag.StringVar(&addr, "l", "0.0.0.0:22222", "Listen")
	flag.Parse()
	go httpreq(http)
	udp(addr)
	fmt.Scanf("%s", &addr)
	//fmt.Print("123")
}
