
from datetime import datetime
import gevent
import time

import logging
logger = logging.getLogger(__name__)

class ScheduleError(Exception):
    pass


class Schedule(object):
    def __init__(self, period, callback):
        """
        Create a simple periodic schedule

        Every `period` seconds callback will be called
        @param period: timedelta object
        @param callback: callable action
        """
        self._period = period.total_seconds()
        logger.debug("period %d" % self._period)
        self._callback = callback
        self._running = False

    def start_at(self, when_utc):
        """
        Start the schedule at a specific time

        @param when_utc: the datetime object for when the schedule should start
        """
        delta = when_utc - datetime.utcnow()
        gevent.spawn_later(delta.total_seconds(), self.start)

    def start(self):
        """
        Start the schedule now
        """
        logger.debug("start")
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
        logger.debug("_next")
        gevent.spawn_later(self._period, self._next)
        gevent.spawn(self._callback)

