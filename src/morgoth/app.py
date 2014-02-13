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

from gevent import monkey; monkey.patch_all()
from gevent.event import Event
import gevent


class App(object):

    def __init__(self):
        self._finish_event = Event()
        self._finish_event.set()
        self._fittings = []


    def _initialize_db(self):
        from morgoth.config import Config
        from morgoth.data.mongo_clients import MongoClients
        from pymongo.errors import OperationFailure
        from pymongo import HASHED

        db_name = Config.get(['mongo', 'database_name'], 'morgoth')
        use_sharding = Config.get(['mongo', 'use_sharding'], True)

        conn = MongoClients.Normal


        if use_sharding:
            try:
                conn.admin.command('enableSharding', db_name)
            except OperationFailure as e:
                if not e.message.endswith('already enabled'):
                    self._logger.error('Error: sharding enabled for morgoth but unable to enable sharding on mongo. See use_sharding config')
                    raise e

        cols = [
                ('meta', [('_id', HASHED)]),
                ('metrics', [('metric', 1), ('time', 1)]),
                ('windows', [('metric', 1), ('ad', 1)])
            ]

        for col, key in cols:
            try:
                conn[db_name][col].ensure_index(key)
                if use_sharding:
                    conn.admin.command(
                        'shardCollection',
                        '%s.%s' % (db_name, col),
                        key=key)
            except OperationFailure as e:
                if not e.message.endswith('already sharded'):
                    self._logger.error('Error: sharding enabled for morgoth but unable to shard %s collection. See use_sharding config' % col)
                    raise e


    def handler(self):
        self._finish_event.clear()
        self._logger.info("Caught signal, shutting down")
        for fitting in self._fittings:
            fitting.stop()
        self._logger.debug("All fittings have been shutdown")
        self._finish_event.set()

    def run(self):
        import logging
        self._logger = logging.getLogger(__name__)

        self._initialize_db()

        # Setup signal handlers
        import signal
        gevent.signal(signal.SIGINT, self.handler)

        # Load the Detectors
        from morgoth.detectors import load_detectors
        load_detectors()

        # Initialize notifiers
        from morgoth.notifiers import load_notifiers
        load_notifiers()


        # Initialize the meta data
        from morgoth.meta import Meta
        Meta.load()

        # Start fittings
        from morgoth.fittings import load_fittings
        self._fittings = load_fittings()

        spawned = []
        for fitting in self._fittings:
            spawn = gevent.spawn(fitting.start)
            spawned.append(spawn)

        for spawn in spawned:
            spawn.join()

        self._logger.info("All fittings have stopped")
        self._finish_event.wait()
        self._logger.info("Finished event set")

def main(args):
    from morgoth  import logger
    logger.init()

    from morgoth.config import Config
    Config.load()

    app = App()
    app.run()

    return 0
