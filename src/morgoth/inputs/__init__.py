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

def load_inputs():
    """ Load the configured Inputs """
    from morgoth.inputs.input import Input
    dirs = [os.path.dirname(__file__)]
    dirs.extend(Config.get(['inputs', 'plugin_dirs'], []))

    pl = PluginLoader()
    mods = pl.find_modules(dirs)

    classes = pl.find_subclasses(mods, Input)

    conf_inputs = [ k for k in Config.inputs.keys() if k != 'plugin_dirs']
    inputs = []
    for input_name, input_class in classes:
        if input_name not in conf_inputs:
            continue
        try:
            logger.debug("Found Input %s", input_name)
            input = input_class.from_conf(Config.inputs.get(input_name, None))
            inputs.append(input)
        except Exception as e:
            logger.warning("Error creating input '%s': %s", input_name, e)

    return inputs

