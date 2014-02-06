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


_NOTIFIERS = {}

def register_notifier(name, n_class):

    """
    Register an Notifier by a given name

    @param name: the name to identify the Notifier
    @param n_class: the class of the Notifier
    """
    if name in _NOTIFIERS:
        raise ValueError('Notifier of name "%s" already exists' % name)
    _NOTIFIERS[name] = n_class

def get_notifier(name):
    """
    Return the Notifier class by a given name

    @return Notifier class
    """
    if name not in _NOTIFIERS:
        raise ValueError("Notifier of name '%s' doesn't exist" % name)
    return _NOTIFIERS[name]

def load_notifiers():
    """ Load the configured notifiers"""
    from morgoth.notifiers.notifier import Notifier
    dirs = [os.path.dirname(__file__)]
    dirs.extend(Config.get(['notifier', 'plugin_dirs'], []))

    pl = PluginLoader()
    mods = pl.find_modules(dirs)

    classes = pl.find_subclasses(mods, Notifier)

    for n_name, n_class in classes:
        logger.debug("Found Notifier %s", n_name)
        try:
            register_notifier(n_name, n_class)
        except ValueError as e:
            logger.warning("Found duplicate Notifiers with name %s", n_name)

