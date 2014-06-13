
# Compatibility issues with python2.6 and python2.7
import sys



def total_seconds(td):
    return (td.microseconds + (td.seconds + td.days * 24 * 3600) * 10**6) / 10**6


def patch_26():
    from datetime import datetime
    datetime.total_seconds = lambda self: total_seconds(self)

def patch():
    assert sys.version_info[0] == 2

    if sys.version_info[1] < 7:
        patch_26()
