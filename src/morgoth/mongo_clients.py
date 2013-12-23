
from gevent import monkey; monkey.patch_all()
import pymongo
from pymongo.read_preferences import ReadPreference

class MongoClients(object):
    Normal = pymongo.MongoClient(tz_aware=True)
    WriteOptimized = pymongo.MongoClient(tz_aware=True, w=0)
    SecondaryPreferred = pymongo.MongoClient(tz_aware=True, read_preference=ReadPreference.SECONDARY_PREFERRED)

