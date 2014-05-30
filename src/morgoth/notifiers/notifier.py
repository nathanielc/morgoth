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


import logging

logger = logging.getLogger(__name__)

class Notifier(object):
    """
    A Notifier base class. Notifiers are plugable so any custom
    notifier can be written

    Notifiers 'notify' of anomalies in any desired manner.
    """
    def __init___(self):
        pass

    @classmethod
    def from_conf(cls, conf, app):
        """
        Create a notifier from the given conf

        @param conf: a conf object
        @param app: reference to the current morgoth application
        """
        return Notifier()


    def notify(self, metric, windows):
        """
        Notify that the window is anomalous

        @param metric: the metric that is considered anomalous
        @param windows: list of window objects from each detector
            where the consensous is that the window is anomalous
        """
        raise NotImplementedError('%s.notify is not implemented' % (self.__class__.__name__))

