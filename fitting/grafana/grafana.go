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

	glog.V(1).Info("Create default dashboard")
	self.defaultDashboard()
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

	config := []byte(fmt.Sprintf(`
define(['settings'],
function (Settings) {
  return new Settings({
    datasources: {
      influxdb: {
        type: 'influxdb',
        url: "http://%s:%d/db/%s",
        username: "%s",
        password: "%s",
        default: true
      },
      grafana: {
        type: 'influxdb',
        url: "http://%s:%d/db/%s",
        username: "%s",
        password: "%s",
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
    window_title_prefix: 'Morgoth Dashboard- ',
    plugins: {
      panels: [],
      dependencies: [],
    }
  });
});`,
		//Influxdb metric config
		self.conf.InfluxDBConf.Host,
		self.conf.InfluxDBConf.Port,
		self.conf.InfluxDBConf.Database,
		self.conf.InfluxDBConf.User,
		self.conf.InfluxDBConf.Password,

		//Influxdb GrafanaDB config
		self.conf.InfluxDBConf.Host,
		self.conf.InfluxDBConf.Port,
		self.conf.GrafanaDB,
		self.conf.InfluxDBConf.User,
		self.conf.InfluxDBConf.Password,
	))

	err := ioutil.WriteFile(filepath, config, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (self *GrafanaFitting) defaultDashboard() error {
	filepath := path.Join(self.conf.Dir, "app", "dashboards", "default.json")

	dash := []byte(`
{
  "id": null,
  "title": "Morgoth",
  "originalTitle": "Morgoth",
  "tags": [],
  "style": "dark",
  "timezone": "browser",
  "editable": true,
  "hideControls": false,
  "sharedCrosshair": false,
  "rows": [
    {
      "title": "test",
      "height": "250px",
      "editable": true,
      "collapse": false,
      "panels": [
        {
          "id": 4,
          "span": 12,
          "type": "graph",
          "x-axis": true,
          "y-axis": true,
          "scale": 1,
          "y_formats": [
            "short",
            "short"
          ],
          "grid": {
            "max": null,
            "min": null,
            "leftMax": null,
            "rightMax": null,
            "leftMin": null,
            "rightMin": null,
            "threshold1": null,
            "threshold2": null,
            "threshold1Color": "rgba(216, 200, 27, 0.27)",
            "threshold2Color": "rgba(234, 112, 112, 0.22)"
          },
          "resolution": 100,
          "lines": true,
          "fill": 1,
          "linewidth": 2,
          "points": false,
          "pointradius": 5,
          "bars": false,
          "stack": false,
          "spyable": true,
          "options": false,
          "legend": {
            "show": true,
            "values": false,
            "min": false,
            "max": false,
            "current": false,
            "total": false,
            "avg": false
          },
          "interactive": true,
          "legend_counts": true,
          "timezone": "browser",
          "percentage": false,
          "nullPointMode": "connected",
          "steppedLine": false,
          "tooltip": {
            "value_type": "cumulative",
            "query_as_alias": true,
            "shared": false
          },
          "targets": [
            {
              "target": "randomWalk('random walk')",
              "function": "first",
              "column": "value",
              "series": "/^m\\..*/",
              "query": "select first(value) from /^m\\..*/ where $timeFilter group by time($interval) order asc",
              "alias": ""
            },
            {
              "target": "",
              "function": "count",
              "column": "value",
              "series": "/^a\\..*/",
              "query": "select count(value) from /^a\\..*/ where $timeFilter group by time($interval) order asc",
              "alias": ""
            }
          ],
          "aliasColors": {},
          "aliasYAxis": {},
          "title": "Morgoth",
          "datasource": "graphite",
          "renderer": "flot",
          "annotate": {
            "enable": false
          },
          "seriesOverrides": [
            {
              "alias": "/^a\\..*/",
              "lines": false,
              "yaxis": 2,
              "points": true,
              "pointradius": 5,
              "stack": true
            }
          ],
          "links": []
        }
      ]
    }
  ],
  "nav": [
    {
      "type": "timepicker",
      "collapse": false,
      "enable": true,
      "status": "Stable",
      "time_options": [
        "5m",
        "15m",
        "1h",
        "6h",
        "12h",
        "24h",
        "2d",
        "7d",
        "30d"
      ],
      "refresh_intervals": [
        "5s",
        "10s",
        "30s",
        "1m",
        "5m",
        "15m",
        "30m",
        "1h",
        "2h",
        "1d"
      ],
      "now": true,
      "notice": false
    }
  ],
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "templating": {
    "list": []
  },
  "annotations": {
    "list": []
  },
  "version": 6,
  "hideAllLegends": false
}`)

	err := ioutil.WriteFile(filepath, dash, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (self *GrafanaFitting) Stop() {
	self.listener.Close()
}
