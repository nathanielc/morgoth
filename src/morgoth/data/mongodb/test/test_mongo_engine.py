
from morgoth.config import Config
from morgoth.data.test.test_engine import EngineTestCase
from morgoth.data.mongodb.engine import MongoEngine

import unittest
import random


class MongoEngineTest(EngineTestCase, unittest.TestCase):

    engine_class = MongoEngine

    engine_conf = """
    host: localhost
    port: 27017
    use_sharding: False
    database: %s
    """



if __name__ == '__main__':
    unittest.main()

