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
    """
    The entry point for launching the morgoth application
    """

    def __init__(self, config):
        """
        Initialize the morgoth append

        @param config: the app configuration object
        @type config: morgoth.config.Config
        """
        self._started_event = Event()
        self._started_event.clear()
        self._finished_event = Event()
        self._finished_event.set()
        self._fittings = []
        self._config = config
        self._engine = None
        self._metrics_manager = None


    def _handler(self):
        self._finished_event.clear()
        self._logger.info("Caught signal, shutting down")
        for fitting in self._fittings:
            fitting.stop()
        self._logger.debug("All fittings have been shutdown")
        self._finished_event.set()

    @property
    def started_event(self):
        return self._started_event

    @property
    def engine(self):
        """
        Data Engine instance
        """
        return self._engine

    @property
    def metrics_manager(self):
        """
        Metrics Manager
        """
        return self._metrics_manager

    @property
    def config(self):
        """
        Morgoth Application configuration
        """
        return self._config

    def run(self):
        try:
            import logging
            self._logger = logging.getLogger(__name__)


            self._logger.info('Setup signal handlers')
            import signal
            gevent.signal(signal.SIGINT, self._handler)

            self._logger.info('Configure detectors and notifiers')
            from morgoth.detectors import configure_detectors
            from morgoth.notifiers import configure_notifiers
            configure_detectors(self)
            configure_notifiers(self)

            self._logger.info('Setup data engine')
            from morgoth.data import load_data_engine
            self._engine = load_data_engine(self)
            self._engine.initialize()


            self._logger.info('Setup metrics manager')
            from morgoth.metrics_manager import MetricsManager
            self._metrics_manager = MetricsManager(self)
            self._logger.info('Inform the metric manager of existing metrics on startup')
            reader = self.engine.get_reader()
            metrics = reader.get_metrics()
            self._logger.debug(metrics)
            self._metrics_manager.new_metrics(metrics)

            self._logger.info('Start fittings')
            from morgoth.fittings import load_fittings
            self._fittings = load_fittings(self)

            spawned = []
            for fitting in self._fittings:
                spawn = gevent.spawn(fitting.start)
                spawned.append(spawn)

            self._started_event.set()
            self._logger.info('Startup complete')

            for spawn in spawned:
                spawn.join()

            self._logger.info("All fittings have stopped")
            self._finished_event.wait()
            self._logger.info("Finished event set")

        except Exception as e:
            self._logger.critical('Error launching morgoth')
            self._logger.exception(e)



def main(config_path):
    from morgoth.compat import patch
    #patch()


    from morgoth.config import Config
    config = Config.load(config_path)

    from morgoth import logger
    logger.init(config)

    app = App(config)
    app.run()

    return 0
