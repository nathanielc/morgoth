
from morgoth.data.engine import Engine
from morgoth.data.mongodb.reader import MongoReader
from morgoth.data.mongodb.writer import MongoWriter
from pymongo.read_preferences import ReadPreference
from pymongo.errors import OperationFailure
from pymongo import HASHED

import pymongo

import logging
logger = logging.getLogger(__name__)

class MongoEngine(Engine):

    def __init__(self, app, host, port, database, use_sharding, writer_options):
        super(MongoEngine, self).__init__(app)
        self._host = host
        self._port = port
        self._database = database
        self._use_sharding = use_sharding
        self._writer_options = writer_options
        self._writer = None
        self._reader = None

    @classmethod
    def from_conf(cls, conf, app):
        host = conf.get('host', 'localhost')
        port = int(conf.get('port', 27017))
        database = conf.get('database', 'morgoth')
        use_sharding = conf.get(['use_sharding'], True)
        writer_options = MongoWriter.get_options(conf)
        writer_options['refresh_interval'] = conf.get('refresh_interval', 60)

        return MongoEngine(
                app,
                host,
                port,
                database,
                use_sharding,
                writer_options,
            )

    def _get_client(self, **options):
        return pymongo.MongoClient(tz_aware=True, **options)

    def initialize(self):
        conn = self._get_client()
        self._reader = MongoReader(conn[self._database])
        self._writer = MongoWriter(self._app, conn[self._database], **self._writer_options)

        super(MongoEngine, self).initialize()

        if self._use_sharding:
            try:
                conn.admin.command('enableSharding', self._database)
            except OperationFailure as e:
                if not e.message.endswith('already enabled'):
                    logger.error(
                        'Error: sharding enabled for morgoth but unable to enable sharding on mongo. See use_sharding config'
                    )
                    raise e

        cols = [
            ('meta', [('_id', HASHED)]),
            ('metrics', [('metric', 1), ('time', 1)]),
            ('windows', [('metric', 1), ('ad', 1)])
        ]

        for col, key in cols:
            try:
                conn[self._database][col].ensure_index(key)
                if self._use_sharding:
                    conn.admin.command(
                        'shardCollection',
                        '%s.%s' % (self._database, col),
                        key=key)
            except OperationFailure as e:
                if not e.message.endswith('already sharded'):
                    logger.error(
                        'Error: sharding enabled for morgoth but unable to shard %s collection. See use_sharding config',
                        col,
                    )
                    raise e

    def get_reader(self):
        return self._reader

    def get_writer(self):
        return self._writer
