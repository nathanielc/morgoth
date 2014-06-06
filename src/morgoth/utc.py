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

import calendar

from datetime import datetime
from dateutil.tz import tzutc

utc = tzutc()

def now():
    """
    Return a datetime object representing the current time in UTC
    """
    return datetime.now(utc)


def to_epoch(dt_utc):
    """
    Convert a datetime object to the number of seconds since the epoch

    @param dt_utc: datetime object in the UTC timezone
    """

    return calendar.timegm(dt_utc.timetuple())

def from_epoch(epoch):
    """
    Create a datetime object in UTC timezone from an epoch timestamp

    @param epoch: number of seconds since the epoch
    """
    dt = datetime.fromtimestamp(epoch)
    print dt.isoformat()
    dt = dt.replace(tzinfo=utc)
    return dt

