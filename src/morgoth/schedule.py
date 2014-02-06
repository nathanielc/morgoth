
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
    def __init__(self, period, callback):
        """
        Create a simple periodic schedule

        Every `period` seconds callback will be called
        @param period: timedelta object
        @param callback: callable action
        """
        self._period = period.total_seconds()
        self._callback = callback
        self._running = False

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
        self.start_at(when)

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
        self._sched.cancel(self._event)
        assert self._sched.empty()
        self._event = None

    def _next(self):
        gevent.spawn_later(self._period, self._next)
        gevent.spawn(self._callback)

