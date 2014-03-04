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

from datetime import datetime, timedelta
from morgoth.utc import now, utc
import gevent
import time

import logging
logger = logging.getLogger(__name__)

class ScheduleError(Exception):
    pass


class Schedule(object):
    DAY = timedelta(days=1)
    def __init__(self, period, callback, delay=timedelta()):
        """
        Create a simple periodic schedule

        Every `period` seconds callback will be called
        @param period: timedelta object
        @param callback: callable action
        """
        self._period = period.total_seconds()
        self._callback = callback
        self._delay = delay
        self._running = False
        self._spawned = None

    def start_aligned(self):
        """
        Start the schedule aligned with midnight
        """
        current = now()
        when = None
        midnight = datetime(
                year=current.year,
                month=current.month,
                day=current.day,
                tzinfo=utc
            )
        period_td = timedelta(seconds=self._period)
        if period_td > self.DAY:
            when =  midnight + self.DAY
            logger.warn("Schedule period greater than a day, aligning with the next midnight")
        else:
            when = midnight
            while when < current:
                when += period_td
        # Start the schedule
        self.start_at(when + self._delay)

    def start_at(self, when_utc):
        """
        Start the schedule at a specific time

        @param when_utc: the datetime object for when the schedule should start
        """
        logger.debug("Starting schedule at %s" % when_utc)
        delta = when_utc - now()
        gevent.spawn_later(delta.total_seconds(), self.start)

    def start(self):
        """
        Start the schedule
        """
        if self._running:
            raise ScheduleError("The schedule is already started")
        self._running = True
        self._next()

    def stop(self):
        """
        Stop the schedule
        """
        self._running = False
        if self._spawned:
            self._spawned.kill()
            self._spawned = None

    def _next(self):
        if self._running:
            self._spawned = gevent.spawn_later(self._period, self._next)
            gevent.spawn(self._callback)

