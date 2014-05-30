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

class Detector(object):
    """
    This class is responsible for dectecting anomalies in metrics

    A single instance of this class will watch many metrics
    """
    def __init__(self):
        self._metrics = set()

    @classmethod
    def from_conf(cls, conf, app):
        """
        Create a Detector from the given conf

        @param conf: a conf object
        @param app: reference to the current morgoth application
        """
        raise NotImplementedError("%s.from_conf is not implemented" % cls.__name__)

    def start(self):
        """
        Start watching metrics for anomalies
        """
        pass

    def watch_metric(self, metric):
        """ Called when a this AD should watch a new metric """
        self._metrics.add(metric)

    def new_value(self, metric, value):
        """
        Called when a given metric receives data

        NOTE: This method is not yet called. It will be supported
        eventually
        """
        pass

    def is_anomalous(self, metric, start, end):
        """
        Return whether the given time range is anomalous

        @param metric: the name of the metric to analyze
        @param start: datetime - marks the start time to analyze
        @param end: datetime - marks the end time to analyze
        @return (bool, window) - a tuple with a bool indicating whether
            the window is considered anomalous and a window object
            with information about why each the detector labeled the window as it did
        """
        raise NotImplementedError("%s.is_anomalous is not implemented" % self.__class__.__name__)
