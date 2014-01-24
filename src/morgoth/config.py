
from types import DictType, ListType
import yaml

import logging
logger = logging.getLogger(__name__)


_conf = None

class ConfigType(type):
    """ Metaclass for Config """
    def __getattr__(cls, key):
        return _conf[key]
    def __getitem(cls, key):
        return _conf[key]

class Config(dict):
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
        super(Config, self).__init__()
        self._accessed = {}
        for k, v in data.iteritems():
            self._accessed[k] = False
            value_type = type(v)
            if value_type == DictType:
                self[k] = Config(v)
            elif value_type == ListType:
                dict_v = {}
                for i in range(len(v)):
                    dict_v[i] = v[i]
                self[k] = ConfigList(dict_v)
            else:
                self[k] = v
        # Rebind methods as instance methods
        self.get = self._get
        self.enumerate = self._enumerate
        self.get_ignored_conf= self._get_ignored_conf

    @classmethod
    def get(cls, attr, default):
        """
        Get conf value

        @param attr: string or list. Return the conf value.
            If list will recusively find the desired conf value
        @param default: the value to return if no conf is found
        """
        return _conf.get(attr, default)


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
                a = self[attr]
                if isinstance(a, Config):
                    ignored[attr] = a.get_ignored_conf()
                else:
                    ignored[attr] = True
        return ignored


    def __getattr__(self, attr):
        return self[attr]

    def __getitem__(self, attr):
        self._accessed[attr] = True
        return super(Config, self).__getitem__(attr)

    def _get(self, attr, default):
        """
        Get conf value

        NOTE: If value is not found then the default will be set
        so that future calls can use the same default
        """
        if type(attr) == ListType:
            if len(attr) == 0:
                return self
            if attr[0] not in self:
                if len(attr) > 1:
                    self[attr[0]] = Config()
                else:
                    self[attr[0]] = default
            value = self[attr[0]]
            if isinstance(value, Config):
                return value.get(attr[1:], default)
            else:
                return value
        else:
            if not attr in self:
                logger.info('Setting default conf "%s" on %s' % (default, self))
                self[attr] = default
            return self[attr]


class ConfigList(Config):
    """
    Config object that make a dict behave like a list for iteration
    """
    def __init__(self, *args, **kwargs):
        super(ConfigList, self).__init__(*args, **kwargs)

    def __iter__(self):
        return self.itervalues()


