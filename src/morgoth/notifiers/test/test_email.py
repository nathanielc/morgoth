

import morgoth.notifiers.email_notifier
from morgoth.notifiers.email_notifier import EmailNotifier
from morgoth.config import Config

import unittest

class EmailNotifierTest(unittest.TestCase):
    conf = Config.loads("""
    toaddrs: user@example.com
    fromaddr: other@example.com
    host: smtphost
    port: 25
    """)
    def test(self):
        """
        Send test email notification
        """
        mock_smtp = MockSMTPLib(self)
        morgoth.notifiers.email_notifier.smtplib = mock_smtp

        email_notifier = EmailNotifier.from_conf(self.conf, None)
        email_notifier.notify('test.metric', ['Windows0'])

        self.assertNotEqual(None, mock_smtp.server)
        self.assertTrue(mock_smtp.server.quit_called)

class MockSMTPLib:
    def __init__(self, test_case):
        self.test_case = test_case
        self.server = None

    def SMTP(self, host_port):
        self.test_case.assertEqual('smtphost:25', host_port)
        self.server = MockServer(self.test_case)
        return self.server

class MockServer:
    def __init__(self, test_case):
        self.test_case = test_case
        self.quit_called = False

    def sendmail(self, fromaddr, to, msg):
        self.test_case.assertEqual('other@example.com', fromaddr)
        self.test_case.assertEqual(1, len(to))
        self.test_case.assertIn('user@example.com', to)

    def quit(self):
        self.quit_called = True


if __name__ == '__main__':
    unittest.main()
