
from types import DictType, ListType
import yaml

import logging
logger = logging.getLogger(__name__)


_conf = None

class ConfigType(type):
    """ Metaclass for Config """
    def __getattr__(cls, key):
        return getattr(_conf, key)

class Config(object):
    """
    Static config class that provides easy access to the morogth config
    """
    __metaclass__ = ConfigType
    @classmethod
    def load(cls, filepath='morgoth.yml'):
        with open(filepath) as f:
            conf_data = yaml.safe_load(f)
        global _conf
        _conf = Config(conf_data)

    def __init__(self, data={}):
        self._accessed = {}
        self._attrs = data.keys()
        # Rebind methods as instance methods
        self.get = self._get
        self.enumerate = self._enumerate
        self.get_ignored_conf= self._get_ignored_conf
        for k, v in data.iteritems():
            self._accessed[k] = False
            if type(v) != DictType:
                setattr(self, k, v)
            else:
                setattr(self, k, Config(v))

    @classmethod
    def get(cls, attr, default):
        """
        Get conf value

        @param attr: string or list. Return the conf value.
            If list will recusively find the desired conf value
        @param default: the value to return if no conf is found
        """
        return _conf.get(attr, default)

    def _get(self, attr, default):
        """
        Get conf value

        NOTE: If value is not found then the default will be set
        so that future calls can use the same default
        """
        if type(attr) == ListType:
            conf = self
            while len(attr) > 0:
                try:
                    conf = getattr(conf, attr[0])
                    attr = attr[1:]
                except AttributeError:
                    if len(attr) > 1:
                        setattr(conf, attr[0], Config())
                    else:
                        setattr(conf, attr[0], default)
            return conf
        else:
            try:
                return getattr(self, attr)
            except AttributeError:
                setattr(self, attr, default)
                return default

    @classmethod
    def enumerate(self):
        """ Enumerate conf value """
        return _conf._attrs

    def _enumerate(self):
        """ Enumerate conf value """
        return self._attrs

    @classmethod
    def get_ignored_conf(cls):
        """ Return dict of confs that were not accessed """
        return _conf.get_ignored_conf()

    def _get_ignored_conf(self):
        ignored = {}
        for attr, accessed in self._accessed.iteritems():
            if not accessed:
                a = getattr(self, attr)
                if type(a) == Config:
                    ignored[attr] = a.get_ignored_conf()
                else:
                    ignored[attr] = True
        return ignored

    def __getattribute__(self, attr):
        object.__getattribute__(self, '_accessed')[attr] = True
        return object.__getattribute__(self, attr)


    ### Iterable methods ###

    def __iter__(self):
        for attr in self._attrs:
            yield attr

    def itervalues(self):
        return self.values()

    def values(self):
        for attr in self._attrs:
            yield self.get(attr, None)

    def iterkeys(self):
        return self.keys()

    def keys(self):
        for attr in self._attrs:
            yield attr

    def iteritems(self):
        return self.items()

    def items(self):
        for attr in self._attrs:
            yield attr, self.get(attr, None)



