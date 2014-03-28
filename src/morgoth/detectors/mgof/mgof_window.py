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

from bson.code import Code
from bson.son import SON
from datetime import datetime
from morgoth.error import MorgothError
from morgoth.window import Window
import os

import logging
logger = logging.getLogger(__name__)

__dir__ = os.path.dirname(__file__)

class MGOFWindow(Window):
    """
    A window object specific to the MGOF algorithm
    """
    map = None
    reduce_code = None
    finalize = None

    def __init__(self, metric,  start, stop, n_bins, trainer=False):
        super(MGOFWindow, self).__init__(metric, start, stop)
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
            data = self._db.windows.find_one({'_id': self.id})
            meta = self._db.meta.find_one({'_id' : self._metric})

            if data is None or data['value']['version'] != meta['version']:
                self._update()
                data = self._db.windows.find_one({'_id': self.id})
                if data is None:
                    return [0] * self._n_bins, 0
            self._prob_dist = data['value']['prob_dist'], data['value']['count']
        if len(self._prob_dist[0]) != self._n_bins:
            raise ValueError('Probability distribution does not have the length n_bins, something with the meta got corrupted')
        return self._prob_dist

    def _update(self):
        """
        Updates the window data
        """
        meta = self._db.meta.find_one({'_id' : self._metric})
        if meta is None:
            raise MorgothError("No such metric '%s'" % self._metric)
        m_max = meta['max']
        m_min = meta['min']
        version = meta['version']

        step_size = ((m_max * 1.01) - m_min) / float(self._n_bins)

        map_values = {
            'id' : self.id,
            'step_size' : step_size,
            'm_min' : m_min,
            'n_bins' :self._n_bins,
        }

        finalize_values = {
            'start' : self._start.isoformat(),
            'stop' : self._stop.isoformat(),
            'version': version,
            'metric' : self._metric,
        }

        map_code = Code(self.map % map_values)
        finalize_code = Code(self.finalize % finalize_values)


        query = {
            'metric' : self._metric,
            'time' : { '$gte' : self._start, '$lt' : self._stop},
        }
        self._db.metrics.map_reduce(map_code, self.reduce_code,
            out=SON([('merge', 'windows')]),
            query=query,
            finalize=finalize_code
        )

# Initialize js code
if MGOFWindow.map is None:
    with open(os.path.join(__dir__, 'window.map.js')) as f:
        MGOFWindow.map = f.read()
    with open(os.path.join(__dir__, 'window.reduce.js')) as f:
        MGOFWindow.reduce_code = Code(f.read())
    with open(os.path.join(__dir__, 'window.finalize.js')) as f:
        MGOFWindow.finalize = f.read()

