
from gevent import monkey; monkey.patch_all()
from gevent.event import Event
import gevent


class App(object):

    def __init__(self):
        self._finish_event = Event()
        self._finish_event.set()
        self._inputs = []


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
        for i in self._inputs:
            i.stop()
        self._logger.debug("All inputs have been shutdown")
        self._finish_event.set()

    def run(self):
        import logging
        self._logger = logging.getLogger(__name__)

        self._initialize_db()

        # Setup signal handlers
        import signal
        gevent.signal(signal.SIGINT, self.handler)

        # Load the ADs
        from morgoth.ads import load_ads
        load_ads()

        # Initialize notifiers
        from morgoth.notifiers import load_notifiers
        load_notifiers()


        # Initialize the meta data
        from morgoth.meta import Meta
        Meta.load()

        # Start input plugins
        from morgoth.inputs import load_inputs
        self._inputs = load_inputs()

        spawned = []
        for input in self._inputs:
            self._logger.info("Staring input '%s'...", input)
            spawn = gevent.spawn(input.start)
            spawned.append(spawn)

        for spawn in spawned:
            spawn.join()


        self._logger.info("All inputs have stopped")
        self._finish_event.wait()
        self._logger.info("Finished event set")

def main(args):
    import logger
    logger.init()

    from morgoth.config import Config
    Config.load()

    app = App()
    app.run()

    return 0
