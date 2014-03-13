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

import unittest
from morgoth.app import main
import gevent

class AppTestCase(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.tdir = tempfile.mkdtemp()
        cls.config_path = os.path.join(cls.tdir, 'morgoth.yml')
        with open(cls.config_path, 'w') as f:
            f.write(cls.conf)
        gevent.spawn(main(config_path=cls.config_path))
        gevent.sleep(0)

    @classmethod
    def tearDownClass(cls):
        os.remove(cls.config_path)
        os.rmdir(cls.tdir)
