package main

import (
	"net/http"
	"log"
	"fmt"
	"os"
	"encoding/json"
)

const defaultAddr=":8080"

//Version information injected at build time
var (
	Service string
	GoVersion string
	GitCommit string
	GOOS string
	GOARCH string
)

//Response defines the json response
type Response struct{
	Service string `json:"Service"`
	GoVersion string `json:"GoVersion"`
	GitCommit string `json:"Gitcommit"`
	GOOS string `json:"GOOS"`
	GOARCH string `json:"GOARCH"`
}

func versionHandler(w http.ResponseWriter, r *http.Request){
	log.Printf("received request: %s %s", r.Method, r.URL.Path)
	response := &Response{
		Service: Service,
		GoVersion: GoVersion,
		GitCommit: GitCommit,
		GOOS: GOOS,
		GOARCH: GOARCH,
	}
	if err:= json.NewEncoder(w).Encode(response); err!=nil{
		http.Error(w, fmt.Sprintf("failed to read entries: %+v", err), http.StatusInternalServerError)
		return
	}
	//log.Print("entries returned version information")
}

func main(){
	addr:= defaultAddr

	if p:=os.Getenv("PORT"); p!=""{
		addr =":"+ p
	}

	log.Printf("server starting to listen on %s", addr)
	http.HandleFunc("/version", versionHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server listen error: %+v", err)
	}
}