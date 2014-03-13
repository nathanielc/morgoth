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
@namespace morgoth.detectors.mgof.mgof

Multinomial Goodness of Fit anomaly detection algorithm
"""

from datetime import timedelta
from morgoth.detectors.detector import Detector
from morgoth.detectors.mgof.mgof_window import MGOFWindow
from morgoth.utils import timedelta_from_str
from scipy.stats import chi2
import numpy


import logging
logger = logging.getLogger(__name__)


class MGOF(Detector):
    """
    Multinomial Goodness Of Fit
    """
    def __init__(self,
                 windows,
                 n_bins=20,
                 normal_count=1,
                 chi2_percentage=0.95):
        """
        Create an anomaly detector.

        This anomaly detector uses a multinomial goodness-of-fit test.

        The algorithm is adapted from this <a href='http://www.hpl.hp.com/techreports/2011/HPL-2011-8.html'>paper</a>

        @param windows: list of (timedelta, timedelta) tuples
            The first element is the offset of the window
            The second element is the duration of the window
        @param period: timedelta object of how much time to wait between
            checking windows for anomalies
        @param duration: the length of the window to analyze, timedelta object
        @param n_bins: the number of discrete values to use in analyzing
            the metric data
        @param normal_count: the number of windows of a given pattern
            that must be found in order to count that pattern as normal
        @param chi2_percentage: a value between 0-1 used to determine the
            statistical probalility of a window matching a given pattern
            Increasing this percentage will make the algorithm less tolerant of
            differences between windows.
        """
        super(MGOF, self).__init__()
        self._windows = windows
        self._n_bins = n_bins
        self._normal_count = normal_count
        self._chi2_percentage = chi2_percentage
        if len(self._windows) <= self._normal_count:
            logger.warn(
                "The normal_count of %d and the number of windows "
                "%d doesn't allow for any bad training windows",
                self._normal_count, len(self._windows)
            )

    @classmethod
    def from_conf(cls, conf):
        windows = []
        for window in conf.windows:
            offset = timedelta()
            if 'offset' in window:
                offset = timedelta_from_str(window.offset)
            if 'duration' in window:
                duration = timedelta_from_str(window.duration)
                windows.append((
                    offset,
                    duration,
                ))
            if 'range' in window:
                range = timedelta_from_str(window.range)
                interval = timedelta_from_str(window.interval)
                start = offset + range
                while start > offset:
                    windows.append((
                        start,
                        interval,
                    ))
                    start -= interval

        return MGOF(
            windows,
            conf.get(['n_bins'], 20),
            conf.get(['normal_count'], 1),
            conf.get(['chi2_percentage'], 0.95)
        )

    def _relative_entropy(self, q, p):
        """
        Calculate the relative entropy of two probability distributions
        """
        assert len(q) == len(p)
        return numpy.sum(q * numpy.log(q / p))

    def is_anomalous(self, metric, start, stop):
        """

        Use the MGOF alogrithm to determine if the given
        metric is anomalous in the specified time range

        @example is_anomalous.py
        """

        windows = []

        for offset, duration in self._windows:
            s = start - offset
            e = s + duration
            w = MGOFWindow(metric, s, e, self._n_bins, trainer=True)
            windows.append(w)

        window = MGOFWindow(metric, start, stop, self._n_bins, trainer=False)
        windows.append(window)

        threshold = chi2.ppf(self._chi2_percentage, self._n_bins - 1)

        m = 0
        prob_distrs = []
        prob_counts = []
        for w in windows:
            anomalous = False
            p, count = w.prob_dist
            p = numpy.array(p)
            if count < self._n_bins:
                #Skip the window, too small
                #logger.debug("Skipped %s %d" % (w, count))
                if not w.trainer:
                    logger.warn("Skipped non training window %s" % w)
                continue
            if m == 0:
                prob_distrs.append(p)
                prob_counts.append(1)
                m = 1
            else:
                min_re = None
                count_index = None
                i = 0
                for i in range(len(prob_distrs)):
                    re = self._relative_entropy(p, prob_distrs[i])
                    # Is the relative entropy statistically significant?
                    if 2 * count * re < threshold:
                        if re < min_re or min_re is None:
                            min_re = re
                            count_index = i

                # Have we seen the prob dist before?
                if count_index is not None:
                    prob_counts[count_index] += 1
                    # Have we seen this prob dist enough?
                    if prob_counts[count_index] <= self._normal_count:
                        anomalous = True
                else:
                    anomalous = True
                    m += 1
                    prob_distrs.append(p)
                    prob_counts.append(1)

            w.anomalous = anomalous
            #logger.debug("Analyzed %s" % w)

        return window
    def __repr__(self):
        return 'MGOF[n_bins:%d,normal_count:%d,chi2_percentage:%0.2f]' % (
                    self._n_bins,
                    self._normal_count,
                    self._chi2_percentage
                )

