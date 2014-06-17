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
"""
@namespace morgoth.detectors.threshold

Simple detector that marks the window as anomalous if the data crosses
a given threshold
"""

from morgoth.detectors.detector import Detector
from morgoth.window import Window

import numpy


import logging
logger = logging.getLogger(__name__)


class Threshold(Detector):
    """
    Threshold detector
    """
    def __init__(self,
            app,
            threshold,
            percentile):
        """
        Create a Threshold detector

        @param threshold: the threshold to use
        @param percentile: which percentile to calculate and compare to the threshold

        Example:
            threshold of 1000 and a percentile of 50 would mean that 50th percentile (median)
            of the data must be greater than the threshold in order to consider the data to be anomalous.
        """
        self._app = app
        self._threshold = threshold
        self._percentile = percentile
        self._reader = self._app.engine.get_reader()


    @classmethod
    def from_conf(cls, conf, app):
        """
        Create a Threshold detector from the conf
        """
        threshold = conf.threshold
        percentile = conf.percentile

        return Threshold(app, threshold, percentile)



    def is_anomalous(self, metric, start, stop):
        """
        Check the data against the threshold
        """
        window = Window(metric, start, stop, self.__class__.__name__)
        percentile = self._reader.get_percentile(metric, self._percentile, start, stop)
        if percentile is None:
            logger.warn('Unable to get percentile for metric %s in %s', metric, window)
            window.anomalous = False
            return window.anomalous, window
        logger.debug('%dth percentile for %s is "%f"' , self._percentile, metric, percentile)

        window.anomalous = percentile > self._threshold
        return window.anomalous, window

    def __repr__(self):
        return 'Threshold[threshold=%f,percentile=%d]' % (
                self._threshold,
                self._percentile,
            )



