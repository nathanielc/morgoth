import logging

from influxdb import client

from morgoth.data.engine import Engine
from morgoth.data.influx.reader import InfluxReader
from morgoth.data.influx.writer import InfluxWriter


logger = logging.getLogger(__name__)

urllib3_logger = logging.getLogger('urllib3')
urllib3_logger.setLevel(logging.ERROR)

class InfluxEngine(Engine):
    def __init__(self,
            app,
            host,
            port,
            user,
            password,
            database,
            writer_options
        ):
        super(InfluxEngine, self).__init__(app)
        self._host = host
        self._port = port
        self._user = user
        self._password = password
        self._database = database
        self._writer_options = writer_options
        self._reader = None
        self._writer = None

    @classmethod
    def from_conf(cls, conf, app):
        host = conf.get('host', 'localhost')
        port = int(conf.get('port', 8086))
        user = conf.get('user', 'morgoth')
        password = conf.get('password', 'morgoth')
        database = conf.get('database', 'morgoth')
        writer_options = InfluxWriter.get_options(conf)

        return InfluxEngine(
                app,
                host,
                port,
                user,
                password,
                database,
                writer_options,
            )

    def _get_client(self):
        return client.InfluxDBClient(
                self._host,
                self._port,
                self._user,
                self._password,
                self._database,
            )

    def initialize(self):

        db = self._get_client()
        self._reader = InfluxReader(db)
        self._writer = InfluxWriter(db, self._app, **self._writer_options)

    def get_reader(self):
        return self._reader

    def get_writer(self):
        return self._writer
