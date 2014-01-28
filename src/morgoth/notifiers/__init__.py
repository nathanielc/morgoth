
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

