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


from morgoth.notifiers.notifier import Notifier

import smtplib

import logging
logger = logging.getLogger(__name__)

class EmailNotifier(Notifier):
    """
    Sends an email with information about each anomaly
    """
    def __init__(self,
            app,
            toaddrs,
            fromaddr,
            host,
            port,
            tls=False,
            username=None,
            password=None):
        """
        Create an email notifier

        @param toaddrs: list of addresses to send to
        @param fromaddr: from address
        @param host: the smtp host
        @param port: the smtp port
        @param tls: whether to use a tls connection
        @param username: if tls is true the username to login
        @param password: if tls is true the password to login
        """
        super(EmailNotifier, self).__init__(app)
        self._toaddrs = toaddrs
        self._fromaddr = fromaddr
        self._host = host
        self._port = port
        self._tls = tls
        self._username = username
        self._password = password

    @classmethod
    def from_conf(cls, conf, app):
        """
        Create a notifier from the given conf

        @param conf: a conf object
        """
        toaddrs = conf.toaddrs.split(',')
        fromaddr = conf.fromaddr
        host = conf.get('host', 'localhost')
        port = int(conf.get('port', 25))
        tls = bool(conf.get('tls', False))
        username = conf.get('username', None)
        password = conf.get('password', None)
        return EmailNotifier(
                app,
                toaddrs,
                fromaddr,
                host,
                port,
                tls,
                username,
                password,
            )


    def notify(self, metric, windows):
        """
        Notify that the window is anomalous
        """
        headers = 'From: %(from)s\r\nTo: %(to)s\r\nSubject: %(subject)s\r\n\r\n' % {
            'from' : self._fromaddr,
            'to' : ', '.join(self._toaddrs),
            'subject' : 'morgoth anomaly %s' % metric
        }
        msg = '''Detected anomaly in metric %(metric)s:
Windows:

    %(windows)s
        ''' % {
            'metric' : metric,
            'windows': '    \n'.join([str(w) for w in windows])
        }
        server = smtplib.SMTP('%s:%d' % (self._host, self._port))
        if self._tls:
            server.starttls()
            server.login(self._username, self._password)
        server.sendmail(self._fromaddr, self._toaddrs, headers + msg)
        server.quit()


    def __repr__(self):
        return 'EmailNotifier[to:%s,from:%s,host:%s,port:%d]' % (
                self._toaddrs,
                self._fromaddr,
                self._host,
                self._port
            )


