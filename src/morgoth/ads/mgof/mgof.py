from datetime import datetime, timedelta
from morgoth.ads.anomaly_detector import AnomalyDetector
from morgoth.utils import timedelta_from_str
from morgoth.window import Window
from scipy.stats import chi2
import numpy


import logging
logger = logging.getLogger(__name__)

class MGOF(AnomalyDetector):
    """
    Multinomial Goodness Of Fit
    """
    def __init__(self,
            windows,
            n_bins=20,
            count_threshold=1,
            chi2_percentage=0.95):
        """
        Create an anomaly detector.

        This anomaly detector uses a multinomial goodness-of-fit test.

        The algorithim is adapted from this paper:
            http://www.hpl.hp.com/techreports/2011/HPL-2011-8.html

        @param metric: the name of the metric to analyze
        @param windows: list of (timedelta, timedelta) tuples
            The first element is the offset of the window
            The second element is the duration of the window
        @param n_bins: the number of discrete values to use in analyzing
            the metric data
        @param count_threshold: the number of windows of a given pattern
            that must be found in order to count the pattern as normal
        @param chi2_percentage: a value between 0-1 used to determine the
            statistical probalility of a window matching a given pattern
        """
        super(MGOF, self).__init__()
        self._windows = windows
        self._n_bins = n_bins
        self._count_threshold = count_threshold
        self._chi2_percentage = chi2_percentage
        if len(self._windows) <= self._count_threshold:
            logger.warn("The count_threshold of %d and the number of windows "
            "%d doesn't allow for any bad training windows"
            % (self._count_threshold, len(self._windows)))

    @classmethod
    def from_conf(cls, conf):
        windows = []
        for window in conf.windows.values():
            windows.append((
                timedelta_from_str(window.offset),
                timedelta_from_str(window.duration),
            ))
        return MGOF(
                windows,
                conf.get(['n_bins'], 20),
                conf.get(['count_threshold'], 1),
                conf.get(['chi2_percentage'], 0.95)
            )

    def _relative_entropy(self, q, p):
        assert len(q) == len(p)
        return numpy.sum(q * numpy.log(q / p))

    def is_anomalous(self, start, end):

        windows = []

        for offset, duration in self._windows:
            s = start - offset
            e = s + duration
            w = Window(self._metric, s, e, self._n_bins, trainer=True)
            windows.append(w)

        window = Window(self._metric, start, end, self._n_bins, trainer=False)
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
                logger.debug("Skipped %s %d" % (w, count))
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
                    if prob_counts[count_index] <= self._count_threshold:
                        anomalous = True
                else:
                    anomalous = True
                    m += 1
                    prob_distrs.append(p)
                    prob_counts.append(1)

            w.anomalous = anomalous
            logger.debug("Analyzed %s" % w)

        return window



