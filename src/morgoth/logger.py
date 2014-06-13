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

import logging

def init(config):
    root_logger = logging.getLogger()
    root_logger.setLevel(_get_level(config))
    ch = logging.StreamHandler()
    formatter = logging.Formatter('[%(asctime)s|%(name)s:%(lineno)d][%(levelname)s] %(message)s')
    ch.setFormatter(formatter)
    root_logger.addHandler(ch)

    print "Initialized Logging"

def _get_level(config):
    """
    Return the logging level from the config
    """
    level = config.get(['logging', 'level'], 'INFO')
    return getattr(logging, level.upper())
