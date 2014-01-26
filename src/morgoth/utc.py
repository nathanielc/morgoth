from datetime import datetime, timedelta, tzinfo

ZERO = timedelta(0)

# A UTC class.

class UTC(tzinfo):
    """UTC"""
    def utcoffset(self, dt):
        return ZERO

    def tzname(self, dt):
        return "UTC"

    def dst(self, dt):
        return ZERO

utc = UTC()

def now():
    """
    Return a datetime object representing the current time in UTC
    """
    return datetime.now(utc)

