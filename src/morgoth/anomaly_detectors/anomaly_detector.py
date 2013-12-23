

class AnomalyDetector(object):
    def __init__(self):
        pass

    def new_metric(self, metric):
        """ Called when a new metric is just created """

    def new_value(self, metric, value):
        """ Called when a given metric receives data """


    def is_anomalous(self, start, end):
        """
        Return whether the given time range is anomalous

        @param start: datetime - marks the start time to analyze
        @param end: datetime - marks the end time to analyze
        @return window - a window object indicating whether it is anomalous
        """
        raise NotImplementedError("%s.is_anomalous is not implemented" % self.__class__.__name__)

