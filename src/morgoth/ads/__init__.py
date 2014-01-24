
from morgoth.config import Config
from morgoth.plugin_loader import PluginLoader
import os

_ADS = {}

def register_ad(name, ad_class):
    """
    Register an AD by a given name

    @param name: the name to identify the AD
    @param ad_class: the class of the AD
    """
    if name in _ADS:
        raise ValueError('AD of name "%s" already exists' % name)
    _ADS[name] = ad_class

def get_ad(name):
    """
    Return the AD class by a given name

    @return AD class
    """
    if name not in _ADS:
        raise ValueError("AD of name '%s' doesn't exist" % name)
    return _ADS[name]


def load_ads():
    """ Load the configured ADs """
    from morgoth.ads.anomaly_detector import AnomalyDetector
    dirs = [os.path.dirname(__file__)]
    dirs.extend(Config.get(['anomaly_detectors', 'auto_load_dirs'], []).values())

    pl = PluginLoader()
    mods = pl.find_modules(dirs)

    print mods
    classes = pl.find_subclasses(mods, AnomalyDetector)

    print classes



