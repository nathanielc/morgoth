
from gevent import monkey; monkey.patch_all()
from gevent.event import Event
import gevent


class App(object):

    def __init__(self):
        self._finish_event = Event()
        self._finish_event.set()

        self._initialize_db()

    def _initialize_db(self):
        from config import Config
        from pymongo import MongoClient, HASHED
        from pymongo.errors import OperationFailure

        db_name = Config.get(['mongo', 'database_name'], 'morgoth')
        use_sharding = Config.get(['mongo', 'use_sharding'], True)

        if use_sharding:
            conn = MongoClient()
            try:
                conn.admin.command('enableSharding', db_name)
            except OperationFailure as e:
                if not e.message.endswith('already enabled'):
                    raise e
            try:
                conn[db_name].meta.ensure_index([('_id', HASHED)])
                conn.admin.command(
                    'shardCollection',
                    '%s.meta' % db_name,
                    key={'_id': HASHED})
            except OperationFailure as e:
                if not e.message.endswith('already sharded'):
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

        # Setup signal handlers
        import signal
        gevent.signal(signal.SIGINT, self.handler)


        # Initialize collector
        from collector import Collector
        collector = Collector()


        # Initialize notifiers
        from notifiers import Notifier
        notifier = Notifier()

        # Initialize watchers
        #watcher = Watcher()

        # Initialize anomaly detectors - we may do this lazily not sure

        # Start input plugins
        from inputs import Graphite
        in_g = Graphite()
        self._inputs = []
        t = gevent.spawn(in_g.start)
        self._inputs.append(in_g)


        t.join()
        self._logger.info("All inputs have stopped")
        self._finish_event.wait()
        self._logger.info("Finished event set")

def main(args):
    import logger
    logger.init()

    from config import Config
    Config.load()

    app = App()
    app.run()

    return 0
