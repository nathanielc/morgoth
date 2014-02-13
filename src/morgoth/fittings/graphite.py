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

import socket
from morgoth.config import Config
from morgoth.utc import utc
from morgoth.data.writer import Writer
from morgoth.fittings.fitting import Fitting
from gevent.server import StreamServer
from datetime import datetime
import gevent

import logging
logger = logging.getLogger(__name__)

class Graphite(Fitting):
    def __init__(self, host, port, stop_timeout, max_pool_size):
        super(Graphite, self).__init__()
        self._port = port
        self._host = host
        self._stop_timeout = stop_timeout
        self._max_pool_size = max_pool_size
        self._writer = Writer()

    @classmethod
    def from_conf(cls, conf):
        host = conf.get('host', '')
        port = conf.get('port', 2003)
        stop_timeout = conf.get('stop_timeout', 10)
        max_pool_size = conf.get('max_pool_size', 1000)
        graphite = Graphite(
                host,
                port,
                stop_timeout,
                max_pool_size
            )
        return graphite


    def start(self):
        logger.info("Starting graphite fitting plugin...")

        # Start gevent server for quick processing of metrics
        self._server = StreamServer(
            listener=(self._host, self._port),
            handle=self._process,
            spawn=self._max_pool_size
        )
        self._server.serve_forever()

    def stop(self):
        self._server.stop(self._stop_timeout)
        self._writer.close()
        logger.debug("Graphite is stopped")

    def _process(self, socket, address):
        """
        Process a request
        """
        logger.debug("New connection from %s" % str(address))
        data = ""
        fileobj = socket.makefile()
        while True:
            line = fileobj.readline()
            if not line:
                logger.debug("Client disconnected %s" % str(address))
                break
            metric, value, timestamp = line.split()
            value = float(value)
            dt_utc = datetime.fromtimestamp(int(timestamp), utc)
            self._writer.insert(dt_utc, metric, value)
            #logger.debug("%s %f %s" % (metric, float(value), dt_utc))
        socket.close()

