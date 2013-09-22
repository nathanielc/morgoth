from datetime import datetime, timedelta
from morgoth.window import Window
from scipy.stats import chi2
import numpy

class AnomalyDetector(object):
    def __init__(self, metric):
        self._metric = metric

    def _get_windows(self, start, end, n_bins):
        windows = []
        for i in range(4, -1, -1):
            delta = timedelta(weeks=i)
            w = Window(self._metric, start - delta, end - delta, n_bins, i != 0)
            windows.append(w)
        return windows


    def D(self, q, p):
        assert len(q) == len(p)
        d = 0
        for i in range(len(q)):
            d += q[i] * numpy.log(q[i] / p[i])
        return d

    def is_anomalous(self, start, end):
        n_bins = 20
        windows = self._get_windows(start, end, n_bins)

        T = chi2.ppf(0.95, n_bins -1)
        c_th = 1

        m = 0
        Ps = []
        cs = []
        for w in windows:
            anomalous = False
            P, count = w.get_P()
            if m == 0:
                Ps.append(P)
                m = 1
                cs.append(1)
            else:
                min_D = None
                c_i = None
                i = 0
                for i in range(len(Ps)):
                    d = self.D(P, Ps[i])
                    if 2 * count * d < T:
                        if d < min_D or min_D is None:
                            min_D = d
                            c_i = i

                if c_i is not None:
                    cs[c_i] += 1
                    if cs[c_i] <= c_th:
                        anomalous = True
                else:
                    anomalous = True
                    m += 1
                    Ps.append(P)
                    cs.append(1)

            if not w.trainer:
                w.anomalous = anomalous

        return windows[-1].anomalous



