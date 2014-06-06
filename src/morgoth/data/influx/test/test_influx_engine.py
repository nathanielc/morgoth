
from morgoth.config import Config
from morgoth.data.test.test_engine import EngineTestCase
from morgoth.data.influx.engine import InfluxEngine

import unittest
import random


class InfluxEngineTest(EngineTestCase, unittest.TestCase):

    engine_class = InfluxEngine

    engine_conf = """
    host: localhost
    user: morgoth
    password: morgoth
    database: morgoth
    smoething: %s
    """


if __name__ == '__main__':
    unittest.main()

