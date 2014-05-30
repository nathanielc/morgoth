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

from datetime import datetime
from morgoth.error import MorgothError
from morgoth.window import Window
import os

import logging
logger = logging.getLogger(__name__)

class MGOFWindow(Window):
    """
    A window object specific to the MGOF algorithm
    """
    map = None
    reduce_code = None
    finalize = None

    def __init__(self, reader, metric,  start, stop, n_bins, trainer=False):
        super(MGOFWindow, self).__init__(metric, start, stop, 'MGOF')
        self._reader = reader
        self._n_bins = n_bins
        self._trainer = trainer
        self._prob_dist = None

    @Window.id.getter
    def id(self):
        """
        Return unique id for this window
        """
        if self._id is None:
            self._id = "%s|%s|%s|%d|mgof" % (
                    self._metric,
                    self._start,
                    self._stop,
                    self._n_bins
                )
        return self._id

    @property
    def trainer(self):
        return self._trainer

    @property
    def prob_dist(self):
        """
        The probability distribution of the window

        @return tuple (
            list of probabalities for each discete value,
            number of data points
        )
        """
        if self._prob_dist is None:
            self._prob_dist = self._reader.get_histogram(self._metric, self._n_bins, self._start, self._stop)
        if len(self._prob_dist[0]) != self._n_bins:
            raise ValueError('Probability distribution does not have the length n_bins, something with the meta got corrupted')
        return self._prob_dist

