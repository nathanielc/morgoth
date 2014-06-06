from datetime import datetime

from morgoth.window import Window
from morgoth.utc import utc
import logging
import unittest

from morgoth.config import Config
from morgoth.notifiers.victorops import VictorOps

logger = logging.getLogger(__name__)
logging.basicConfig()

class TestVictorOps(unittest.TestCase):

    conf = Config.loads("""
    url: https://alert.victorops.com/integrations/generic/20131114/alert/API_TOCKEN/
    routing_key: morgoth
    routing_key_pattern: (\w+)\..*
    """)


    def test_01(self):
        vo = VictorOps.from_conf(self.conf, None)
        metric = 'test_metrics.002'
        vo.notify(metric, [Window(metric, datetime(2014, 6, 6,), datetime(2014, 6, 6, 1))])
        # Need to add some assertions which require mocking urllib2...




if __name__ == '__main__':
    unittest.main()
