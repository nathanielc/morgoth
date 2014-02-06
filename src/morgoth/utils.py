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
import re

td_pattern = re.compile(r'^(\d+(\.\d+)?)([smhdwy])')

def timedelta_from_str(s):
    """
    Return timedelta object represented by the given string

    @param s: string matching td_pattern
    """
    match = td_pattern.match(s)
    if not match:
        raise ValueError('String %s not a valid time delta string' % s)
    num, _, unit = match.groups()
    num = float(num)
    if unit == 's':
        return timedelta(seconds=num)
    elif unit == 'm':
        return timedelta(minutes=num)
    elif unit == 'h':
        return timedelta(hours=num)
    elif unit == 'd':
        return timedelta(days=num)
    elif unit == 'w':
        return timedelta(weeks=num)
    elif unit == 'y':
        return timedelta(days=num * 365)
    else:
        assert False # regex pattern should garauntee a match

