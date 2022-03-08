/**
 * Copyright 2017 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// [START all]
package main

import (
	"flag"
	"fmt"
	"net"
	"log"
	"net/http"
)

var httpAddress string
var udpAddress string

func init() {
	flag.StringVar(&httpAddress, "http-addr", ":8080", "Address of the HTTP server for /metrics and /health ")
	flag.StringVar(&udpAddress, "udp-addr", ":25001", "Address of the UDP server")
}



func serveUdp(pc net.PacketConn, addr net.Addr, buf []byte) {
	log.Printf("Received %d bytes from %s: %s", len(buf), addr, string(buf))

	// try to write a response to the caller
	pc.WriteTo([]byte(fmt.Sprintf("%d bytes received\n", len(buf))), addr)
}

func main() {
	// http server (for /metrics and /healthz)
	go func() {
		mux := http.NewServeMux()

		// Enable observability, add /metrics endpoint to extract and examine stats.
		enableObservabilityAndExporters(mux)

		// Enable health check /healthz endpoint
		enableHealthCheck(mux)

		/* TODO: print out stats about the udp server ? */
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "OK")

		})

		// start the web server on port and accept requests
		log.Printf("HTTP server listening on %s", httpAddress)
		err := http.ListenAndServe(httpAddress, mux)
		log.Fatal(err)
	}()

	// listen to incoming udp packets
	udpAddr, err := net.ResolveUDPAddr("udp", udpAddress)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UDP server listening on %s", udpAddress)
	defer conn.Close()

	// infinite loop
	for {
		// note here max 1024 byte messages
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %s", err.Error())
			continue
		}
		// serve in its own goroutine, note the buffer copy
		go serveUdp(conn, addr, buf[:n])
	}
}
// [END all]
