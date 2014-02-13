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

from morgoth.config import Config
from morgoth.plugin_loader import PluginLoader
import os

import logging

logger = logging.getLogger(__name__)

def load_fittings():
    """ Load the configured Fittings """
    from morgoth.fittings.fitting import Fitting
    dirs = [os.path.dirname(__file__)]
    dirs.extend(Config.get(['fittings', 'plugin_dirs'], []))

    pl = PluginLoader()
    mods = pl.find_modules(dirs)

    classes = pl.find_subclasses(mods, Fitting)

    conf_fittings = [ k for k in Config.fittings.keys() if k != 'plugin_dirs']
    fittings = []
    for fitting_name, fitting_class in classes:
        if fitting_name not in conf_fittings:
            continue
        try:
            logger.debug("Found Fitting %s", fitting_name)
            fitting = fitting_class.from_conf(Config.fittings.get(fitting_name, None))
            fittings.append(fitting)
        except Exception as e:
            logger.warning("Error creating fitting '%s': %s", fitting_name, e)

    return fittings

