
import random
import numpy
import unittest
from morgoth.collector import Collector
from morgoth.window import Window
from datetime import datetime


class WindowTestCase(unittest.TestCase):
    def create_metric_data(self, metric):
        self.delete_metric_data(metric)
        c = Collector()
        data = []
        for h in range(3):
            for m in range(60):
                for s in range(60):
                    value = h*60*60 + m*60 + s
                    c.insert(
                        datetime(2013, 9, 21, h, m ,s),
                        metric,
                        value
                        )
                    data.append(value)
        return data
    def delete_metric_data(self, metric):
        c = Collector()
        c.delete_metric(metric)

    def calc_prob_distr(self, data, n_bins):
        d_min = min(data)
        d_max = max(data) * 1.01
        step = (d_max - d_min) / float(n_bins)

        discrete = [int((v - d_min) / step) for v in data]
        P = numpy.ones(n_bins)
        for v in discrete:
            P[v] += 10

        P /= len(data) * 10 + n_bins


        return P

    def test_window_simple(self):
        metric = 'test_window_simple'
        data = self.create_metric_data(metric)

        for n_bins in range(10,100,10):
            print n_bins

            expected_P = self.calc_prob_distr(data, n_bins)

            w = Window(metric, datetime(2013, 9 ,21, 0), datetime(2013, 9, 21, 3), n_bins)
            P, count = w.get_prob_distr()

            self.assertEqual(len(data), count)
            self.assertAlmostEqual(1, sum(P))

            self.assertEqual(n_bins, len(P))
            for i in range(n_bins):
                self.assertAlmostEqual(expected_P[i], P[i])

if __name__ == '__main__':
    unittest.main()
