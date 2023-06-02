package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gliderlabs/ssh"
)

type Tunnel struct {
	w      io.Writer
	donech chan struct{}
}

var tunnels = map[int]chan Tunnel{}

func main() {

	go func() {

		http.HandleFunc("/", handleRequest)
		log.Fatal(http.ListenAndServe(":3000", nil))

	}()

	ssh.Handle(func(s ssh.Session) {
		id := rand.Intn(math.MaxInt)
		tunnels[id] = make(chan Tunnel)
		fmt.Println("Tunnel Id ->", id)
		tunnel := <-tunnels[id]
		fmt.Println("Tunnel is ready")
		_, err := io.Copy(tunnel.w, s)
		if err != nil {
			log.Fatal(err)
		}
		close(tunnel.donech)
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	tunnel, ok := tunnels[id]
	if !ok {
		log.Fatal("Error: Cannot Find ")
	}
	donech := make(chan struct{})
	tunnel <- Tunnel{
		w:      w,
		donech: donech,
	}
	<-donech
}
