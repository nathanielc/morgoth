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

from morgoth.app import App
import gevent
import os
import tempfile
import unittest

class AppTestCase(unittest.TestCase):
    """
    Base class for app (end to end) tests
    """

    conf = None
    def set_up_app(self, conf):
        tdir = tempfile.mkdtemp()
        config_path = os.path.join(tdir, 'morgoth.yml')
        with open(config_path, 'w') as f:
            f.write(conf)

        app = self._run_app(config_path)
        return (app, tdir, config_path)

    def _run_app(self, config_path):
        from morgoth import logger
        logger.init()

        from morgoth.config import Config
        config = Config.load(config_path)

        app = App(config)
        gevent.spawn(app.run)
        gevent.sleep(0)
        app.started_event.wait()
        return app



    def tear_down_app(self, app, tdir, config_path):
        os.remove(config_path)
        os.rmdir(tdir)
