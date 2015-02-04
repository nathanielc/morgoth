package grafana

import (
	"fmt"
	"github.com/nvcook42/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	app "github.com/nvcook42/morgoth/app/types"
	"net"
	"net/http"
	"os"
	"path"
	"io"
	"io/ioutil"
)

type GrafanaFitting struct {
	conf     GrafanaConf
	listener net.Listener
}

func (self *GrafanaFitting) Name() string {
	return "Grafana"
}

func (self *GrafanaFitting) Start(app app.App) {
	glog.V(1).Info("Starting grafana fitting", self.conf)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", self.conf.Port))
	if err != nil {
		glog.Error("Error starting Grafana fitting: ", err.Error())
		return
	}
	self.listener = listener
	err = self.prepare()
	if err != nil {
		glog.Error("Error starting Grafana fitting: ", err.Error())
		return
	}

	http.Serve(self.listener, http.FileServer(http.Dir(self.conf.Dir)))
}

func (self *GrafanaFitting) prepare() error {
	glog.V(1).Info("Downloading grafana")
	filepath, err := self.download()
	if err != nil {
		return err
	}

	glog.V(1).Info("Extracting grafana")
	err = untar(filepath, self.conf.Dir)
	if err != nil {
		return err
	}

	glog.V(1).Info("Configuring grafana")
	self.configure()
	if err != nil {
		return err
	}
	return nil
}

func (self *GrafanaFitting) download() (string, error) {
	if _, err := os.Stat(self.conf.Dir); os.IsNotExist(err) {
		err = os.Mkdir(self.conf.Dir, 0755)
		if err != nil {
			return "", err
		}
	}
	name := path.Base(self.conf.URL)
	filepath := path.Join(self.conf.Dir, name)
	if stat, err := os.Stat(filepath); err == nil && stat.Size() > 0 {
		//We already downloaded the tar skipping
		return filepath, nil
	}
	// Start download
	resp, err := http.Get(self.conf.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	io.Copy(file, resp.Body)

	return filepath, nil
}

func (self *GrafanaFitting) configure() error {
	filepath := path.Join(self.conf.Dir, "config.js")

	config := []byte(`
define(['settings'],
function (Settings) {
  return new Settings({
    datasources: {
      influxdb: {
        type: 'influxdb',
        url: "http://localhost:8086/db/morgoth",
        username: 'morgoth',
        password: 'morgoth',
      },
      grafana: {
        type: 'influxdb',
        url: "http://localhost:8086/db/grafana",
        username: 'morgoth',
        password: 'morgoth',
        grafanaDB: true
      },
    },
    search: {
      max_results: 20
    },
    default_route: '/dashboard/file/default.json',
    unsaved_changes_warning: true,
    playlist_timespan: "1m",
    admin: {
      password: ''
    },
    window_title_prefix: 'Grafana - ',
    plugins: {
      panels: [],
      dependencies: [],
    }
  });
});`)

	err := ioutil.WriteFile(filepath, config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (self *GrafanaFitting) Stop() {
	self.listener.Close()
}
