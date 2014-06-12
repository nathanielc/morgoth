

from morgoth.notifiers.email_notifier import EmailNotifier
from morgoth.config import Config

import unittest

class EmailNotifierTest(unittest.TestCase):
    conf = Config.loads("""
    toaddrs: user@example.com
    fromaddr: other@example.com
    """)
    def test(self):
        """
        Send test email notification
        """
        email_notifier = EmailNotifier.from_conf(self.conf, None)
        email_notifier.notify('test.metric', ['Windows0'])


if __name__ == '__main__':
    unittest.main()
