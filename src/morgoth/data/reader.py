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
from morgoth.utc import utc
import re

import logging
logger = logging.getLogger(__name__)

class Reader(object):
    """
    Class that provides read access to the metric data and anomalies
    """
    def __init__(self):

    def get_metrics(self, pattern=None):
        pass

    def get_data(self, metric, start=None, stop=None, step=None):
        pass

    def get_anomalies(self, metric, start=None, stop=None):
        pass
