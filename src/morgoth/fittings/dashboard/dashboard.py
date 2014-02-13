#
# Copyright 2014 Nathaniel Cook
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


from flask import Flask, request, jsonify, make_response, current_app
from gevent.pywsgi import WSGIServer
from morgoth.fittings.fitting import Fitting

import logging
logger = logging.getLogger(__name__)

app = Flask(__name__)

class Dashboard(Fitting):
    def __init__(self, host, port):
        super(Dashboard, self).__init__()
        self._host = host
        self._port = port


    @classmethod
    def from_conf(cls, conf):
        host = ''
        port = conf.get('port', 8080)
        return Dashboard(host, port)


    def start(self):
        logger.info("Starting Dashboard fitting...")
        self._server = WSGIServer((self._host, self._port), app, log=None)
        self._server.serve_forever()

    def stop(self):
        self._server.stop()

##############################
# Flask methods
##############################


@app.route('/')
def root():
    return app.send_static_file('dashboard.html')

