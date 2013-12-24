
from morgoth.data.mongo_clients import MongoClients
from bson.objectid import ObjectId
import os

__dir__ = os.path.dirname(__file__)

class Window(object):

    def __init__(self, metric,  start, end, trainer=False):
        self._db = MongoClients.Normal.morgoth
        self._metric = metric
        self._start = start
        self._end = end
        self._trainer = False
        self.__id = None
        self._anomalous = None

    @property
    def _id(self):
        if self.__id is None:
            self.__id = ObjectId()
        return self.__id

    @property
    def trainer(self):
        return self._trainer

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
        return self._start, self._end

    def __repr__(self):
        return self.__str__()

    def __str__(self):
        return "{Window[%s:%s]anomalous:%s}" % (self._start, self._end, self.anomalous)
