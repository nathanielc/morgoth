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

from morgoth.error import MorgothError
from morgoth.plugin_loader import PluginLoader
from types import ListType
import os

import logging

logger = logging.getLogger(__name__)


def load_data_engine(app):
    """
    Load the configured Data Engines

    Note: It is possible to configure more than one data engine.
    This behavior is not supported and will result in an error.

    @param app: the morgoth application instance
    @type app: morgoth.app.App
    """
    from morgoth.data.engine import Engine
    dirs = [os.path.dirname(__file__)]
    dirs.extend(app.config.get(['plugin_dirs', 'data_engines'], []))

    pl = PluginLoader(app, dirs, Engine)
    engine = None
    try:
        engines = pl.load(app.config.data_engine)
        if len(engines) > 1:
            raise MorgothError('Only one data engine is supported. Please only configure one')
        engine = engines[0]

    except KeyError:
        logger.error('No data engine config found')
    except Exception as e:
        logger.error("Error creating data engine %s", e)
        raise e

    return engine

