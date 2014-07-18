from datetime import datetime

from morgoth.window import Window
from morgoth.date_utils import utc
import logging
import unittest
import morgoth.notifiers.victorops

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
        morgoth.notifiers.victorops.urllib2 = MockURLLib2(self)

        vo = VictorOps.from_conf(self.conf, None)
        metric = 'test_metrics.002'
        vo.notify(metric, [Window(metric, datetime(2014, 6, 6,), datetime(2014, 6, 6, 1))])
        # Need to add some assertions which require mocking urllib2...


class MockURLLib2:
    def __init__(self, test_case):
        self.test_case = test_case
        self.request = None

    def Request(self, url, data, headers):
        self.test_case.assertEqual(self.test_case.conf.url + 'test_metrics', url)
        self.test_case.assertNotEqual(0, len(data))
        self.request = object()
        return self.request

    def urlopen(self, req):
        self.test_case.assertEqual(self.request, req)
        return open('/dev/null')




if __name__ == '__main__':
    unittest.main()
