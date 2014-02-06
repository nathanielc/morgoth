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
    def __init___(self):
        pass

    @classmethod
    def from_conf(cls, conf):
        """
        Create a notifier from the given conf

        @param conf: a conf object
        """
        return Notifier()


    def notify(self, window):
        """
        Notify that the window is anomalous
        """
        logger.info("Window: %s is anomalous" % window)

