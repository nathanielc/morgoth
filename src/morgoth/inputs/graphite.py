
import socket
from morgoth.config import Config
from morgoth.utc import utc
from morgoth.data.writer import Writer
from morgoth.inputs.input import Input
from gevent.server import StreamServer
from datetime import datetime
import gevent

import logging
logger = logging.getLogger(__name__)

class Graphite(Input):
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
        logger.info("Starting graphite input plugin...")

        # Start gevent server for quick processing of metrics
        self._server = StreamServer(
            listener=(self._host, self._port),
            handle=self._process,
            spawn=self._max_pool_size
        )
        self._server.serve_forever()

    def stop(self):
        self._server.stop(self._stop_timeout)
        logger.debug("Server stopped")
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

