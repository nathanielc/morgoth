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

def load_outputs():
    """ Load the configured Outputs """
    from morgoth.outputs.output import Output
    dirs = [os.path.dirname(__file__)]
    dirs.extend(Config.get(['outputs', 'plugin_dirs'], []))

    pl = PluginLoader()
    mods = pl.find_modules(dirs)

    classes = pl.find_subclasses(mods, Output)

    conf_outputs = [ k for k in Config.outputs.keys() if k != 'plugin_dirs']
    outputs = []
    for output_name, output_class in classes:
        if output_name not in conf_outputs:
            continue
        try:
            logger.debug("Found Output %s", output_name)
            output = output_class.from_conf(Config.outputs.get(output_name, None))
            outputs.append(output)
        except Exception as e:
            logger.warning("Error creating output '%s': %s", output_name, e)

    return outputs

