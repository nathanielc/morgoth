#
# Copyright 2014 Nathaniel Cook
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from gevent import monkey; monkey.patch_all()
import pymongo
from pymongo.read_preferences import ReadPreference

class MongoClients(object):
    """
    Class that provides easy access to different types of mongo clients
    """
    Normal = pymongo.MongoClient(tz_aware=True)
    WriteOptimized = pymongo.MongoClient(tz_aware=True, w=0)
    SecondaryPreferred = pymongo.MongoClient(tz_aware=True, read_preference=ReadPreference.SECONDARY_PREFERRED)

