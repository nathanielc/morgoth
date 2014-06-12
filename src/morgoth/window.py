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

from bson.objectid import ObjectId
import os

__dir__ = os.path.dirname(__file__)

class Window(object):
    """
    Represents a window in data for a given metric
    """

    def __init__(self, metric, start, stop, detector_name=None):
        self._metric = metric
        self._start = start
        self._stop = stop
        self._id = None
        self._anomalous = None
        self._detector_name = detector_name

    @property
    def metric(self):
        """
        Return the metric this window applies to
        """
        return self._metric

    @property
    def detector_name(self):
        """
        Return the detector for this window
        """
        return self._detector_name

    @property
    def start(self):
        """
        Return the start time in UTC of the window
        """
        return self._start

    @property
    def stop(self):
        """
        Return the stop time in UTC of the window
        """
        return self._stop

    @property
    def id(self):
        if self._id is None:
            self._id = ObjectId()
        return self._id


    @property
    def anomalous(self):
        """ Return whether the window is anomalous
            NOTE: this property is `None` if it has not been determined
        """
        return self._anomalous

    @anomalous.setter
    def anomalous(self, value):
        self._anomalous = value

    @property
    def range(self):
        return self._start, self._stop

    def __repr__(self):
        return self.__str__()

    def __str__(self):
        return "{Window[%s|%s|%s|anomalous:%s|detector:%s]}" % (
                self._metric,
                self._start,
                self._stop,
                self.anomalous,
                self._detector_name
            )
