
from datetime import datetime, timedelta
from morgoth.config import Config
from morgoth.utc import utc

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
            if attrname.startswith('_test'):
                newattrs[attrname[1:]] = lambda self: EngineTestType._test_wrapper(self, attrname)


        logger.debug(newattrs.keys())

        return super(EngineTestType, cls).__new__(cls, name, bases, newattrs)

    @staticmethod
    def _test_wrapper(instance, test_name):
        logger.debug("Instance: %s", instance)
        logger.debug("Test Name: %s", test_name)
        engine, app = instance._create_engine(instance.engine_class, instance._new_config())

        test_method = getattr(instance, test_name)
        test_method(instance, engine, app)


class EngineTestCase(unittest.TestCase):
    __metaclass__ = EngineTestType


    def _create_engine(self, engine_class, engine_conf, app=None):
        if app is None:
            app = MockApp()
        return engine_class.from_conf(engine_conf, app), app

    def _test_initialize(self, engine, app):
        self.assertEqual(0, app.metrics_manager.new_metric_count)

    def _test_01(self, engine, app):

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 29, 1, tzinfo=utc)

        metric = 'test_engine_' + engine.__class__.__name__ + str(random.randint(0, 100))
        writer.insert(start, metric, 42)
        self.assertEqual(1, app.metrics_manager.new_metric_count)

        writer.flush()
        metrics = reader.get_metrics()
        self.assertEqual([metric], metrics)

        data = reader.get_data(metric)
        self.assertEqual(1, len(data))
        self.assertEqual((start.isoformat(), 42), data[0])

    def _test_02(self, engine, app):

        metric = 'test_engine_' + engine.__class__.__name__ + str(random.randint(0, 100))

        writer = engine.get_writer()
        reader = engine.get_reader()

        start = datetime(2014, 5, 30, 1, tzinfo=utc)
        count = 100
        stop = start + timedelta(seconds=count -1)

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

        metric = 'test_engine_' + engine.__class__.__name__ + str(random.randint(0, 100))

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
        self.assertAlmostEqual(1, sum(hist))

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
            self.assertAlmostEqual(expected_hist[i], hist[i])




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



