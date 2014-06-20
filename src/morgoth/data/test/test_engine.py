
from dateutil import parser
from datetime import datetime, timedelta
from morgoth.config import Config
from morgoth.date_utils import utc
import numpy

import random

import unittest

import logging
logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.DEBUG)


class EngineTestType(type):
    def __new__(cls, name, bases, attrs):
        newattrs = {}
        for attrname, value in attrs.items():
            newattrs[attrname] = value
            if attrname.startswith('_test_04'):
                newattrs[attrname[1:]] = lambda self, attrname=attrname: self._do_test(attrname)

        return super(EngineTestType, cls).__new__(cls, name, bases, newattrs)

class EngineTestCase(object):
    __metaclass__ = EngineTestType

    engine_class = None
    engine_conf = None

    def _new_config(self):
        db_name = "test_engine_db_%d" % random.randint(0, 1000)
        return Config.loads(self.engine_conf % db_name)

    def _create_engine(self, engine_class, engine_conf, app=None):
        if app is None:
            app = MockApp()
        return engine_class.from_conf(engine_conf, app), app

    def _destroy_engine(self, engine_conf):
        pass

    def _do_test(self, test):
        engine_conf = self._new_config()
        engine, app = self._create_engine(self.engine_class, engine_conf)
        try:
            engine.initialize()
            getattr(self, test)(engine, app)
        finally:
            self._destroy_engine(engine_conf)
            pass


    def _test_initialize(self, engine, app):
        self.assertEqual(0, app.metrics_manager.new_metric_count)

    def _test_01(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 29, 1, tzinfo=utc)

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_01'
        writer.insert(start, metric, 42)
        self.assertEqual(1, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual([metric], metrics)

        data = reader.get_data(metric)
        self.assertEqual(1, len(data))
        self.assertEqual((start.isoformat(), 42), data[0])

    def _test_02(self, engine, app):

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_02'

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 30, 1, tzinfo=utc)
        count = 100
        stop = start + timedelta(seconds=count - 1)

        expected_data = []
        for i in range(count):
            cur = start + timedelta(seconds=i)
            value = i*i
            expected_data.append((cur.isoformat(), value))
            writer.insert(cur, metric, value)

        self.assertEqual(count, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual([metric], metrics)

        data = reader.get_data(metric)
        self.assertEqual(count, len(data))
        self.assertEqual(expected_data, data)


        data = reader.get_data(metric, start=start)
        self.assertEqual(count, len(data))
        self.assertEqual(expected_data, data)

        data = reader.get_data(metric, stop=stop)
        self.assertEqual(count, len(data))
        self.assertEqual(expected_data, data)

        data = reader.get_data(metric, start=start, stop=stop)
        self.assertEqual(count, len(data))
        self.assertEqual(expected_data, data)

        half = start + ((stop - start) / 2)

        data = reader.get_data(metric, start=start, stop=half)
        self.assertEqual(count / 2, len(data))
        self.assertEqual(expected_data[:count/2], data)


    def _test_03(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 29, 1, tzinfo=utc)

        count = 10
        expected_metrics = []
        for i in range(count):
            metric = 'test_engine_' + engine.__class__.__name__ + str(i)
            expected_metrics.append(metric)

            writer.insert(start, metric, 42)

        self.assertEqual(count, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual(set(expected_metrics), set(metrics))

        for metric in metrics:
            data = reader.get_data(metric)
            self.assertEqual(1, len(data))
            self.assertEqual((start.isoformat(), 42), data[0])


    def _test_04(self, engine, app):

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_04'

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 30, 1, tzinfo=utc)
        count = 100
        stop = start + timedelta(seconds=count)

        expected_data = []
        for i in range(count):
            cur = start + timedelta(seconds=i)
            value = i*i
            expected_data.append((cur.isoformat(), value))
            writer.insert(cur, metric, value)

        self.assertEqual(count, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual([metric], metrics)

        n_bins = 10

        hist, hist_count = reader.get_histogram(metric, n_bins, start, stop)

        self.assertEqual(count, hist_count)
        self.assertAlmostEqual(1, sum(hist), places=1)

        expected_hist = [
                0.3178217821782178,
                0.1297029702970297,
                0.1,
                0.0801980198019802,
                0.0801980198019802,
                0.0702970297029703,
                0.060396039603960394,
                0.0504950495049505,
                0.060396039603960394,
                0.0504950495049505
            ]

        self.assertEqual(len(expected_hist), len(hist))
        for i in range(len(expected_hist)):
            self.assertAlmostEqual(expected_hist[i], hist[i], places=2)


    def _test_05(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 29, 1, tzinfo=utc)

        metric_count = 10
        for i in range(metric_count):
            metric = 'test_engine_%s_test05_%d' % (engine.__class__.__name__, i)
            writer.insert(start, metric, 42)
            self.assertEqual(i + 1, app.metrics_manager.new_metric_count)

            writer.flush()
            metrics = reader.get_metrics(metric)
            self.assertEqual([metric], metrics)

            data = reader.get_data(metric)
            self.assertEqual(1, len(data))
            self.assertEqual((start.isoformat(), 42), data[0])

        metrics = reader.get_metrics()
        self.assertEqual(metric_count, len(metrics))

        metrics = reader.get_metrics(r'test05_[0-4]')
        self.assertEqual(5, len(metrics))

        metrics = reader.get_metrics(r'test05_[5-9]')
        self.assertEqual(5, len(metrics))

    def _test_06(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_06'
        start = datetime(2014, 5, 29, 1, tzinfo=utc)
        stop = datetime(2014, 5, 29, 2, tzinfo=utc)

        writer.record_anomalous(metric, start, stop)

        anomalies = reader.get_anomalies(metric)
        self.assertEqual(1, len(anomalies))
        self.assertEqual(start.isoformat(), anomalies[0]['start'])
        self.assertEqual(stop.isoformat(), anomalies[0]['stop'])

        anomalies = reader.get_anomalies(metric, start=start)
        self.assertEqual(1, len(anomalies))
        self.assertEqual(start.isoformat(), anomalies[0]['start'])
        self.assertEqual(stop.isoformat(), anomalies[0]['stop'])

        anomalies = reader.get_anomalies(metric, stop=stop)
        self.assertEqual(1, len(anomalies))
        self.assertEqual(start.isoformat(), anomalies[0]['start'])
        self.assertEqual(stop.isoformat(), anomalies[0]['stop'])


        anomalies = reader.get_anomalies(metric, start=start, stop=stop)
        self.assertEqual(1, len(anomalies))
        self.assertEqual(start.isoformat(), anomalies[0]['start'])
        self.assertEqual(stop.isoformat(), anomalies[0]['stop'])

        anomalies = reader.get_anomalies(metric, start=stop)
        self.assertEqual(0, len(anomalies))

        anomalies = reader.get_anomalies(metric, stop=start)
        self.assertEqual(0, len(anomalies))


    def _test_07(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_07'
        data = numpy.linspace(1, 100, 100)
        start = datetime(2014, 5, 29, 1, tzinfo=utc)
        stop = datetime(2014, 5, 29, 9, tzinfo=utc)
        step = (stop - start) / len(data)

        curr = start
        for d in data:
            writer.insert(curr, metric, d)
            curr += step


        writer.flush()
        _data = reader.get_data(metric)
        self.assertEqual(len(data), len(_data))

        percentile_50 = reader.get_percentile(metric, 50)
        self.assertAlmostEqual(50, percentile_50)
        percentile_75 = reader.get_percentile(metric, 75)
        self.assertAlmostEqual(75, percentile_75)
        percentile_90 = reader.get_percentile(metric, 90)
        self.assertAlmostEqual(90, percentile_90)
        percentile_95 = reader.get_percentile(metric, 95)
        self.assertAlmostEqual(95, percentile_95)
        percentile_99 = reader.get_percentile(metric, 99)
        self.assertAlmostEqual(99, percentile_99)

    def _test_08(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_08'

        data = numpy.linspace(1, 100, 100)
        start = datetime(2014, 5, 29, 1, tzinfo=utc)
        stop = datetime(2014, 5, 29, 9, tzinfo=utc)
        step = (stop - start) / len(data)

        curr = start
        for d in data:
            writer.insert(curr, metric, d)
            curr += step


        writer.flush()
        _data = reader.get_data(metric)
        self.assertEqual(len(data), len(_data))

        def _test_step_size(size):
            _data = reader.get_data(metric, step=step*size)
            self.assertEqual(len(data)/size, len(_data))
            self.assertAlmostEqual(sum(data), size*(sum([d[1] for d in _data])))

        sizes = [1, 2, 4]

        for size in sizes:
            _test_step_size(size)


        _data = reader.get_data(metric, start=start, step=step*2)
        self.assertEqual(len(data)/2, len(_data))

        _data = reader.get_data(metric, stop=stop, step=step*2)
        self.assertEqual(len(data)/2, len(_data))

        _data = reader.get_data(metric, start, stop, step=step*2)
        self.assertEqual(len(data)/2, len(_data))

    def _test_09(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        metric = 'test_engine_' + engine.__class__.__name__ + 'test_09'

        start = datetime(2014, 5, 29, 1, tzinfo=utc)
        stop = datetime(2014, 5, 29, 2, tzinfo=utc)

        writer.insert(start, metric, 42)
        self.assertEqual(1, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual([metric], metrics)

        data = reader.get_data(metric)
        self.assertEqual([(start.isoformat(), 42)], data)

        writer.record_anomalous(metric, start, stop)

        anomalies = reader.get_anomalies(metric)
        self.assertEqual(1, len(anomalies))
        self.assertEqual(start.isoformat(), anomalies[0]['start'])
        self.assertEqual(stop.isoformat(), anomalies[0]['stop'])

        # Delete metric
        writer.delete_metric(metric)

        metrics = reader.get_metrics()
        self.assertEqual([], metrics)

        data = reader.get_data(metric)
        self.assertEqual([], data)

        anomalies = reader.get_anomalies(metric)
        self.assertEqual(0, len(anomalies))




class MockMetricsManager(object):
    new_metric_count = 0
    def new_metric(self, *args, **kwargs):
        self.new_metric_count += 1


class MockApp(object):
    config = Config.loads("""
    metrics:
        .*:
            key:value
    """)
    def __init__(self):
        self.metrics_manager = MockMetricsManager()



