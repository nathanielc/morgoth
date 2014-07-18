
from morgoth.config import Config
from morgoth.data.test.int_test_engine import EngineTestCase
from morgoth.data.influx.engine import InfluxEngine

import unittest
import random
from influxdb import client


class InfluxEngineTest(EngineTestCase, unittest.TestCase):

    engine_class = InfluxEngine

    engine_conf = """
    host: localhost
    user: morgoth
    password: morgoth
    database: %s
    """

    def _create_engine(self, engine_class, engine_conf, app=None):
        conn = client.InfluxDBClient(
                engine_conf.get('host', 'localhost'),
                engine_conf.get('port', 8086),
                'root',
                'root',
                engine_conf.database
            )

        conn.create_database(engine_conf.database)
        conn.add_database_user(engine_conf.user, engine_conf.password)
        return super(InfluxEngineTest, self)._create_engine(engine_class, engine_conf, app)


    def _destroy_engine(self, engine_conf):
        conn = client.InfluxDBClient(
                engine_conf.get('host', 'localhost'),
                engine_conf.get('port', 8086),
                'root',
                'root',
                engine_conf.database
            )

        try:
            conn.delete_database(engine_conf.database)
        except:
            pass

if __name__ == '__main__':
    unittest.main()

