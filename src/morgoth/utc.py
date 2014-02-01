

from datetime import datetime
from dateutil.tz import tzutc

utc = tzutc()

def now():
    """
    Return a datetime object representing the current time in UTC
    """
    return datetime.now(utc)

