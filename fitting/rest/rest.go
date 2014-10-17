package rest

import (
	app "github.com/nvcook42/morgoth/app/types"
	"github.com/nvcook42/morgoth/engine"
	"net"
	"net/http"
	log "github.com/cihub/seelog"
	"github.com/ant0ine/go-json-rest/rest"
	"fmt"
)

type RESTFitting struct {
	port uint
	reader engine.Reader
	listener net.Listener
}

func (self *RESTFitting) Name() string {
	return "REST"
}

func (self *RESTFitting) Start(app app.App) {
	self.reader = app.GetReader()
	handler := rest.ResourceHandler{}
	handler.SetRoutes(
		&rest.Route{"GET", "/healthy", self.healthy},
	)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", self.port))
	if err != nil {
		log.Error("Error starting REST fitting %s", err.Error())
		return
	}
	self.listener = listener
	err = http.Serve(self.listener, &handler)
	if err != nil {
		log.Debug(err)
		return
	}
	log.Info("REST fitting is done")
}

func (self * RESTFitting) healthy(w rest.ResponseWriter, req *rest.Request) {
	w.WriteJson(42)
}

func (self *RESTFitting) Stop() {
	self.listener.Close()
}
