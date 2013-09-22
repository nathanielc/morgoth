#!/usr/bin/python

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

