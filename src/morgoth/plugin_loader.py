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
import os
import imp
import sys
import inspect


import logging

logger = logging.getLogger(__name__)

class PluginLoader(object):
    """
    A class to facilitate loading plugins dynamically
    """
    def __init__(self, search_dirs, base_class, depth=1, ignored=set()):
        """

        @param search_dirs: list of directories to search in
        @param depth: the depth of subdirs to search
        @param base_class: the base class that all plugins should extend
        """
        self._search_dirs = search_dirs
        self._base_class = base_class
        self._depth = depth
        self._ignored = ignored
        self._classes = {}

        # Find modules and classes
        modules = self._find_modules(self._search_dirs, depth)
        self._classes = self._find_subclasses(modules, self._base_class, self._ignored)

    def load(self, plugins_conf):
        """
        Parses conf and loads discovered plugins into the `container`

        @param plugins_conf: list of Config objects
            plugins_conf must be dictionary config where each value is either
            another dictionary config or a list of dictionary configs. Each
            of the nested dictionries will be passed as the 'conf' to the
            'from_conf' method on the plugin class.
        """
        plugins = []
        for name, confs in plugins_conf.items():
            plugin_class = self.get_plugin_class(name)
            if not confs.is_list:
                confs = [confs]
            for conf in confs:
                try:
                    plugin = plugin_class.from_conf(conf)
                    plugins.append(plugin)
                except Exception as e:
                    logger.error('Could not create "%s" from conf "%s" Error: "%s"',  name, conf, e)
                    raise e
        return plugins

    def get_plugin_class(self, name):
        """
        Return the class for a given plugin

        @param name: the name of the class of the plugin
        """
        try:
            return self._classes[name]
        except KeyError as e:
            raise LoaderError('No such plugin %s' % name)

    def _find_modules(self, search_dirs, depth=1):
        """
        Find a modules in `search_dirs`
        """
        return self.__find_modules(search_dirs, depth, None)

    def __find_modules(self, search_dirs, depth, package):
        """
        Find modules recursively in search_dirs

        @param search_dirs: list of directories to search in
        @param depth: the depth of subdirs to search
        @param package: the name of the package to import modules into
        """
        mods = []
        for search_dir in search_dirs:
            if not os.path.isdir(search_dir): continue
            for entry in os.listdir(search_dir):
                path = os.path.join(search_dir, entry)
                if os.path.isdir(path) and depth > 0:

                    #Find package
                    found_pkg = None
                    try:
                        found_pkg = imp.find_module(entry, [search_dir])
                    except ImportError:
                        pass

                    if not found_pkg:
                        continue

                    if package:
                        name = '%s.%s' % (package, entry)
                    else:
                        name = entry
                    pkg = imp.load_module(name, *found_pkg)
                    mods.append(pkg)
                    mods.extend(self.__find_modules(
                        [path],
                        depth - 1,
                        name,
                    ))
                elif entry.endswith('.py'):
                    entry = entry[:-3]
                    if package:
                        name = '%s.%s' % (package, entry)
                    else:
                        name = entry

                    # Find modules under package
                    found_mod = imp.find_module(entry, [search_dir])
                    if not found_mod:
                        continue
                    mod = imp.load_module(name, *found_mod)
                    mods.append(mod)
        return mods

    def _find_subclasses(self, mods, parent_class, ignored=set()):
        """
        Return all classes found in the list of modules `mods` that are subclasses
        of `parent_class`

        @param mods: list of module objects
        @param parent_class: a parent class
        @param ignore: set of class names to ignore.
            NOTE: `parent_class` is always ignored
        @return dict of class name to class
        """
        classes = {}
        ignored.add(parent_class.__name__)
        for mod in mods:
            members = inspect.getmembers(
                        mod,
                        lambda x:
                            inspect.isclass(x)
                            and issubclass(x, parent_class)
                            and x.__name__ not in ignored
                        )
            for name, member in members:
                if name in classes:
                    raise LoaderError("Plugin of name '%s' already exists, found another in module '%s'" % (name, mod.__name__))
                classes[name] = member
        return classes


class LoaderError(Exception):
    pass
