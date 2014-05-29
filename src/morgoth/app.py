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
        self._finish_event = Event()
        self._finish_event.set()
        self._fittings = []
        self._config = config


    def handler(self):
        self._finish_event.clear()
        self._logger.info("Caught signal, shutting down")
        for fitting in self._fittings:
            fitting.stop()
        self._logger.debug("All fittings have been shutdown")
        self._finish_event.set()

    def run(self):
        try:
            import logging
            self._logger = logging.getLogger(__name__)


            # Setup signal handlers
            import signal
            gevent.signal(signal.SIGINT, self.handler)

            # Setup data engine
            from morgoth.data import load_data_engine
            engine = load_data_engine(self._config)
            engine.initialize()

            # Configure detectors and notifiers
            from morgoth.detectors import configure_detectors
            from morgoth.notifiers import configure_notifiers
            configure_detectors(self._config)
            configure_notifiers(self._config)

            # Start fittings
            from morgoth.fittings import load_fittings
            self._fittings = load_fittings(self._config)

            spawned = []
            for fitting in self._fittings:
                spawn = gevent.spawn(fitting.start)
                spawned.append(spawn)

            for spawn in spawned:
                spawn.join()

            self._logger.info("All fittings have stopped")
            self._finish_event.wait()
            self._logger.info("Finished event set")

        except Exception as e:
            self._logger.critical('Error launching morgoth')
            self._logger.exception(e)



def main(config_path):
    from morgoth import logger
    logger.init()

    from morgoth.config import Config
    config = Config.load(config_path)

    app = App(config)
    app.run()

    return 0
