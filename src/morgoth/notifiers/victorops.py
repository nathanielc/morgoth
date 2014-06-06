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

import json
import logging
import re
import urllib2

from morgoth.notifiers.notifier import Notifier
from morgoth.utc import to_epoch


logger = logging.getLogger(__name__)

class VictorOps(Notifier):
    """
    Push anomalies to a VictorOps HTTP REST API
    """
    def __init__(self, app, url, routing_key, routing_key_pattern):
        super(VictorOps, self).__init__(app)
        self._url = url.rstrip('/')
        self._routing_key = routing_key
        self._routing_key_ptn = routing_key_pattern

    @classmethod
    def from_conf(cls, conf, app):
        """
        Create a notifier from the given conf

        @param conf: a conf object
        @param app: reference to the current morgoth application
        """
        url = conf.url
        routing_key = conf.get('routing_key', None)
        routing_key_pattern = conf.get('routing_key_pattern', None)
        return VictorOps(
                app,
                url,
                routing_key,
                routing_key_pattern,
            )


    def notify(self, metric, windows):
        """
        Notify that the window is anomalous

        @param metric: the metric that is considered anomalous
        @param windows: list of window objects from each detector
            where the consensous is that the window is anomalous
        """
        start = to_epoch(windows[0].start)
        stop = to_epoch(windows[0].stop)

        msg = ['%s to %s anomalous!' % (windows[0].start, windows[0].stop)]
        for window in windows:
            msg.append('%s: %s' % (window.detector_name, window.anomalous))


        data = {
            'message_type' : 'CRITICAL',
            'entity_id' : metric,
            'state_start_time' : start,
            'entity_is_host' : False,
            'monitoring_tool' : 'morgoth',
            'entity_display_name' : metric,
            'state_message' : '\n'.join(msg)

        }
        data = json.dumps(data)

        routing_key = self._routing_key
        if self._routing_key_ptn:
            matches = re.match(self._routing_key_ptn, metric)
            if matches and len(matches.groups()) > 0:
                routing_key = matches.groups()[0]

        url = self._url + '/' + routing_key

        req = urllib2.Request(url, data, {'Content-Type': 'application/json'})

        f = urllib2.urlopen(req)
        response = f.read()
        f.close()

        return
