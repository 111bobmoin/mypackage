package main

import (
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/riotpot/internal/globals"
	lr "github.com/riotpot/internal/logger"
	"github.com/riotpot/internal/services"
)

var Plugin string
var vmIP = "192.168.000.000"//根据自己要连接的虚拟机ip决定

const (
	name    = "HTTP"
	network = globals.TCP
	port    = 80
)

func init() {
	Plugin = "Httpd"
}

func Httpd() services.Service {
	mx := services.NewPluginService(name, port, network)

	return &Http{
		mx,
	}
}

type Http struct {
	// Anonymous fields from the mixin
	services.Service
}

func (h *Http) Run() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(h.valid))

	srv := &http.Server{
		Addr:    h.GetAddress(),
		Handler: mux,
	}

	go h.serve(srv)

	return
}

func (h *Http) serve(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lr.Log.Fatal().Err(err)
	}
}

// 访问主机系统和网络
func accessHostSystem(vmIP string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s", vmIP))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// This function handles connections made to a valid path
func (h *Http) valid(w http.ResponseWriter, req *http.Request) {
	response, err := accessHostSystem(vmIP)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	//response := fmt.Sprintf("%s%s", head, body)

	fmt.Fprint(w, response)
}
