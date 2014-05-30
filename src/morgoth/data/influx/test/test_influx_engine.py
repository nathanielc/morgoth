
from morgoth.config import Config
from morgoth.data.test.test_engine import EngineTestCase
from morgoth.data.influx.engine import InfluxEngine

import unittest
import random


class InfluxEngineTest(EngineTestCase):

    engine_class = InfluxEngine

    engine_conf = """
    host: localhost
    port: 27017
    user: root,
    password: root,
    database: %s
    """

    def _new_config(self):
        db_name = "test_influx_engine_db_%d" % random.randint(0, 100)
        return Config.loads(self.engine_conf % db_name)


#    def test_initialize(self):
#        engine, app = self._create_engine(self.engine_class, self._new_config())
#        self._test_initialize(engine, app)
#
#    def test_01(self):
#        engine, app = self._create_engine(self.engine_class, self._new_config())
#        engine.initialize()
#        self._test_01(engine, app)
#
#    def test_02(self):
#        engine, app = self._create_engine(self.engine_class, self._new_config())
#        engine.initialize()
#        self._test_02(engine, app)
#
#    def test_03(self):
#        engine, app = self._create_engine(self.engine_class, self._new_config())
#        engine.initialize()
#        self._test_03(engine, app)
#
#    def test_04(self):
#        engine, app = self._create_engine(self.engine_class, self._new_config())
#        engine.initialize()
#        self._test_04(engine, app)



if __name__ == '__main__':
    unittest.main()

