
from morgoth.data.mongo_clients import MongoClients
from bson.objectid import ObjectId
import os

__dir__ = os.path.dirname(__file__)

class Window(object):

    def __init__(self, metric,  start, stop):
        self._db = MongoClients.Normal.morgoth
        self._metric = metric
        self._start = start
        self._stop = stop
        self._id = None
        self._anomalous = None

    @property
    def metric(self):
        """
        Return the metric this window applies to
        """
        return self._metric

    @property
    def start(self):
        """
        Return the start time in UTC of the window
        """
        return self._start

    @property
    def stop(self):
        """
        Return the stop time in UTC of the window
        """
        return self._stop

    @property
    def id(self):
        if self._id is None:
            self._id = ObjectId()
        return self._id


    @property
    def anomalous(self):
        """ Return whether the window is anomalous
            NOTE: this property is `None` if it has not been determined
        """
        return self._anomalous

    @anomalous.setter
    def anomalous(self, value):
        self._anomalous = value

    @property
    def range(self):
        return self._start, self._stop

    def __repr__(self):
        return self.__str__()

    def __str__(self):
        return "{Window[%s|%s|%s|anomalous:%s]}" % (
                self._metric,
                self._start,
                self._stop,
                self.anomalous,
            )
