
from mongo_client_factory import MongoClientFactory
from bson.code import Code
from bson.son import SON
from datetime import datetime
from morgoth.error import MorgothError
import os

__dir__ = os.path.dirname(__file__)

class Window(object):
    map = None
    reduce = None
    finalize = None

    def __init__(self, metric,  start, end, n_bins, trainer=False):
        self._db = MongoClientFactory.create().morgoth
        self._metric = metric
        self._start = start
        self._end = end
        self._n_bins = n_bins
        self._trainer = False
        self.__id = None
        self._anomalous = None
        self._prob_dist = None

        # Initialize js code
        if self.map is None:
            with open(os.path.join(__dir__, 'window.map.js')) as f:
                self.map = f.read()
            with open(os.path.join(__dir__, 'window.reduce.js')) as f:
                self.reduce = Code(f.read())
            with open(os.path.join(__dir__, 'window.finalize.js')) as f:
                self.finalize = Code(f.read())

    @property
    def _id(self):
        if self.__id is None:
            self.__id = "%s:%d:%s:%s" % (self._metric, self._n_bins, self._start.isoformat(), self._end.isoformat())
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
    def prob_dist(self):
        """
        The probability distribution of the window

        @return (list of probabalities for each discete value,
                    number of data points)
        """
        if self._prob_dist is None:
            data = self._db.windows.find_one({'_id': self._id})
            meta = self._db.meta.find_one({'_id' : self._metric})
            if data is None or data['value']['version'] != meta['version']:
                self._update()
                data = self._db.windows.find_one({'_id': self._id})
                if data is None:
                    return [0] * self._n_bins, 0
            self._prob_dist = data['value']['P'], data['value']['count']
        return self._prob_dist

    @property
    def range(self):
        return self._start, self._end

    def _update(self):
        """
        Updates the window data
        """
        meta = self._db.meta.find_one({'_id' : self._metric})
        if meta is None:
            raise MorgothError("No such metric '%s'" % self._metric)
        m_max = meta['max']
        m_min = meta['min']
        version = meta['version']

        step_size = ((m_max * 1.01) - m_min) / float(self._n_bins)

        map_values = {
            'id' : self._id,
            'step_size' : step_size,
            'm_min' : m_min,
            'n_bins' :self._n_bins,
            'version': version,
            'start_h' : self._start.hour,
            'start_m' : self._start.minute,
            'start_s' : self._start.second,
            'end_h' : self._end.hour,
            'end_m' : self._end.minute,
            'end_s' : self._end.second,
        }

        map = Code(self.map % map_values)

        start_hour = datetime(
            self._start.year,
            self._start.month,
            self._start.day,
            self._start.hour)

        end_hour = datetime(
            self._end.year,
            self._end.month,
            self._end.day,
            self._end.hour)


        self._db.metrics.map_reduce(map, self.reduce,
            out=SON([('merge', 'windows'), ('db', 'morgoth')]),
            query={
                'metric' : self._metric,
                'time' : { '$gte' : start_hour, '$lte' : end_hour},
            },
            finalize=self.finalize
        )

    def __repr__(self):
        return self.__str__()

    def __str__(self):
        return "Window [%s:%s] anomalous: %s" % (self._start, self._end, self.anomalous)
