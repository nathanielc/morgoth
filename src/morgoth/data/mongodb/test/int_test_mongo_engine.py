
from morgoth.config import Config
from morgoth.data.test.int_test_engine import EngineTestCase
from morgoth.data.mongodb.engine import MongoEngine

import unittest
import random
import pymongo


class MongoEngineTest(EngineTestCase, unittest.TestCase):

    engine_class = MongoEngine

    engine_conf = """
    host: localhost
    port: 27017
    use_sharding: False
    database: %s
    """

    def _destroy_engine(self, engine_conf):
        try:
            conn = pymongo.MongoClient(
                engine_conf.host,
                engine_conf.port,
                tz_aware=True
            )
            conn.drop_database(engine_conf.database)
        except:
            pass


if __name__ == '__main__':
    unittest.main()
