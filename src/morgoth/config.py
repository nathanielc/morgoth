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

from types import DictType, ListType
import yaml
import os

import logging
logger = logging.getLogger(__name__)


class Config(dict):
    """
    Static config class that provides easy access to the morogth config
    """
    @classmethod
    def load(cls, filepath='morgoth.yml'):
        with open(filepath) as f:
            conf_data = yaml.load(f, Loader)
        return Config(conf_data)

    @classmethod
    def loads(cls, data):
        conf_data = yaml.safe_load(data)
        return Config(conf_data)

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

        self._is_list = False

    @property
    def is_list(self):
        return self._is_list

    def enumerate(self):
        """ Enumerate conf value """
        return self._attrs


    def get_ignored_conf(self):
        """
        Return list of config options that were never accessed
        """
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

    def get(self, attr, default):
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
                self[attr] = default
            return self[attr]


class ConfigList(Config):
    """
    Config object that make a dict behave like a list for iteration
    """
    def __init__(self, *args, **kwargs):
        super(ConfigList, self).__init__(*args, **kwargs)
        self._is_list = True

    def __iter__(self):
        return self.itervalues()


class Loader(yaml.Loader):
    """
    yaml.Loader that allows for relative includes of other yaml files
    """
    def __init__(self, stream):
        super(Loader, self).__init__(stream)
        self._root = os.path.dirname(stream.name)

    def _include(self, node):
        filename = os.path.join(self._root, self.construct_scalar(node))
        with open(filename, 'r') as f:
            return yaml.load(f, Loader)

Loader.add_constructor('!include', Loader._include)

