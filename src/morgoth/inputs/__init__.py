
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

