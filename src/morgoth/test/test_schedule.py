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

from datetime import timedelta
from morgoth.schedule import Schedule, ScheduleError
from morgoth.date_utils import total_seconds

import unittest
import gevent
import time

class TestSchedule(unittest.TestCase):
    """
    Unittest for morgoth.schedule.Schedule
    """
    def test_schedule_00(self):
        self.period = timedelta(seconds=1)
        self.last = time.time()
        self.called = 0
        sched = Schedule(self.period, self._callback)
        sched.start()

        self.assertRaises(ScheduleError, sched.start)

        i = 0
        max_calls = 3
        while self.called < max_calls:
            gevent.sleep(total_seconds(self.period) / 2)
            if i > max_calls * 3:
                raise AssertionError('Callback not called often enough')
            i += 1

        sched.stop()
        time.sleep(total_seconds(self.period) / 2)
        self.assertEqual(max_calls, self.called)


    def _callback(self):
        self.assertGreater(time.time(), self.last)
        self.last += total_seconds(self.period)
        self.called += 1


if __name__ == '__main__':
    unittest.main()
