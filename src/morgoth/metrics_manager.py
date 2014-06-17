
from morgoth.metric_supervisor import MetricSupervisor, NullMetricSupervisor
from collections import OrderedDict

import re

import logging
logger = logging.getLogger(__name__)

class MetricsManager(object):
    """
    Manages metric supervisors
    """
    _null_supervisor = NullMetricSupervisor()
    def __init__(self, app):
        self._app = app
        self._supervisors = OrderedDict()
        self._metrics = set()

        writer = self._app.engine.get_writer()

        # Load supervisors from conf
        metrics_conf = self._app.config.metrics
        if not metrics_conf.is_list:
            metrics_conf = [metrics_conf]
        logger.debug('Pattern matching order...')
        for metric_conf in metrics_conf:
            for pattern, conf in metric_conf.items():
                pattern = str(pattern)
                logger.debug(pattern)
                self._supervisors[pattern] = MetricSupervisor(writer, pattern, conf)

    def new_metrics(self, metrics):
        """
        Initialize a supervisor for the given metrics
        """
        for metric in metrics:
            self.new_metric(metric)

    def new_metric(self, metric):
        """
        Initialize a supervisor for a given metric
        """
        if metric not in self._metrics:

            supervisor = self._match_metric(metric)
            supervisor.add_metric(metric)

            supervisor.start()

            self._metrics.add(metric)


    def _match_metric(self, metric):
        """
        Determine which pattern matches the given metric

        @param metric: the name of the metric
        @return the MetricSupervisor for the given metric
        """
        for pattern, supervisor in self._supervisors.items():
            if re.match(pattern, metric):
                return supervisor

        # No config for the metric
        logger.warn("Metric '%s' has no matching configuration", metric)
        return self._null_supervisor

