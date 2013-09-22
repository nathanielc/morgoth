
from pymongo import MongoClient

class MongoClientFactory(object):
    @classmethod
    def create(cls):
        return MongoClient(tz_aware=True)
