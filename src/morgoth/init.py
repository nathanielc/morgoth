#!/usr/bin/python
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

# Initialize the MongoDB databases and collections


from pymongo import MongoClient
from pymongo.errors import OperationFailure

conn = MongoClient()
try:
    conn.admin.command('enableSharding', 'morgoth')
except OperationFailure:
    pass
try:
    conn.admin.command(
        'shardCollection',
        'morgoth.metrics',
        key={'_id': 1})
except OperationFailure:
    pass

conn.morgoth.metrics.ensure_index('time')
conn.morgoth.metrics.ensure_index('metric')
conn.morgoth.windows.ensure_index('value.metric')

