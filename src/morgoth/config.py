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
from collections import OrderedDict
import yaml
import os


import logging
logger = logging.getLogger(__name__)


class Config(OrderedDict):
    """
    Static config class that provides easy access to the morgoth config
    """
    @classmethod
    def load(cls, filepath='morgoth.yml'):
        with open(filepath) as f:
            conf_data = yaml.load(f, IncludeLoader)
        return Config(conf_data)

    @classmethod
    def loads(cls, data):
        conf_data = yaml.load(data, OrderedLoader)
        return Config(conf_data)

    def __init__(self, data=None):
        super(Config, self).__init__()
        if data is None:
            data = {}
        self._accessed = {}
        for k, v in data.iteritems():
            self._accessed[k] = False
            value_type = type(v)
            if value_type == DictType or value_type == OrderedDict:
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
        try:
            return self[attr]
        except KeyError:
            return object.__getattribute__(self, attr)

    def __getitem__(self, attr):
        object.__getattribute__(self, '_accessed')[attr] = True
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
        base_iter = super(ConfigList, self).__iter__()
        for k in base_iter:
            yield self[k]


class OrderedLoader(yaml.Loader):
    """
    Preserves the order of dictionaries
    """
    def __init__(self, stream):
        super(OrderedLoader, self).__init__(stream)

        self.add_constructor(u'tag:yaml.org,2002:map', type(self).construct_yaml_map)
        self.add_constructor(u'tag:yaml.org,2002:omap', type(self).construct_yaml_map)

    def construct_yaml_map(self, node):
        data = OrderedDict()
        yield data
        value = self.construct_mapping(node)
        data.update(value)

    def construct_mapping(self, node, deep=False):
        if isinstance(node, yaml.MappingNode):
            self.flatten_mapping(node)
        else:
            raise yaml.constructor.ConstructorError(None, None,
                'expected a mapping node, but found %s' % node.id, node.start_mark)

        mapping = OrderedDict()
        for key_node, value_node in node.value:
            key = self.construct_object(key_node, deep=deep)
            try:
                hash(key)
            except TypeError, exc:
                raise yaml.constructor.ConstructorError('while constructing a mapping',
                    node.start_mark, 'found unacceptable key (%s)' % exc, key_node.start_mark)
            value = self.construct_object(value_node, deep=deep)
            mapping[key] = value
        return mapping

class IncludeLoader(OrderedLoader):
    """
    Allows for relative includes of other yaml files
    """
    extensions = ['.yaml', '.yml', '.conf']
    def __init__(self, stream):
        super(IncludeLoader, self).__init__(stream)
        self._root = os.path.dirname(stream.name)

        self.add_constructor('!include', type(self)._include)
        self.add_constructor('!include_dir', type(self)._include_dir)

    def _include(self, node):
        filename = os.path.join(self._root, self.construct_scalar(node))
        with open(filename, 'r') as f:
            return yaml.load(f, IncludeLoader)

    def _include_dir(self, node):
        directory = os.path.join(self._root, self.construct_scalar(node))
        data = []
        for filename in os.listdir(directory):
            extension = os.path.splitext(filename)[1]
            filepath = os.path.join(directory, filename)
            if extension not in self.extensions:
                logger.warn('skipping %s when including confs from %s', filename, directory)
                continue
            with open(filepath, 'r') as f:
                data.append(yaml.load(f, IncludeLoader))
        return data

