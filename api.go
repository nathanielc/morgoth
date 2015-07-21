package morgoth

import (
	"fmt"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/ant0ine/go-json-rest/rest"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"net"
	"net/http"
)

type APIServer struct {
	app      *App
	port     uint
	listener net.Listener
	handler  rest.ResourceHandler
}

func NewAPIServer(app *App, port uint) *APIServer {
	return &APIServer{
		app:  app,
		port: port,
	}
}

func (self *APIServer) Start() (err error) {
	self.handler = rest.ResourceHandler{
		EnableStatusService: true,
	}
	self.handler.SetRoutes(
		&rest.Route{"GET", "/status", self.status},
		&rest.Route{"GET", "/stats", self.stats},
	)
	self.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", self.port))
	if err != nil {
		glog.Errorf("Error binding APIServer: %s", err)
		return
	}
	go func() {
		err = http.Serve(self.listener, &self.handler)
		if err != nil {
			glog.Errorf("Error while APIServer was running: %s", err)
		}
	}()
	glog.V(1).Info("APIServer running")
	return
}

func (self *APIServer) Stop() {
	self.listener.Close()
}

func (self *APIServer) status(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(self.handler.GetStatus())
}

func (self *APIServer) stats(w rest.ResponseWriter, req *rest.Request) {
	stats := self.app.Stats
	w.WriteJson(stats)
}
