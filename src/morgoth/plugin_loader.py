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
    def __init__(self):
        """
        """

    def find_modules(self, search_dirs, depth=1):
        """
        Find a modules in `search_dirs`
        @param search_dirs: list of directories to search in
        @param depth: the depth of subdirs to search
        """
        return self._find_modules(search_dirs, depth, None)
    def _find_modules(self, search_dirs, depth, package):
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
                    found_pkg = imp.find_module(entry, [search_dir])
                    if not found_pkg:
                        continue

                    if package:
                        name = '%s.%s' % (package, entry)
                    else:
                        name = entry
                    pkg = imp.load_module(name, *found_pkg)
                    mods.append(pkg)
                    mods.extend(self._find_modules(
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

    def find_subclasses(self, mods, parent_class, ignored=set()):
        """
        Return all classes found in the list of modules `mods` that are subclasses
        of `parent_class`

        @param mods: list of module objects
        @param parent_class: a parent class
        @param ignore: set of class names to ignore.
            NOTE: `parent_class` is always ignored
        """
        classes = []
        ignored.add(parent_class.__name__)
        for mod in mods:
            classes.extend(
                    inspect.getmembers(
                        mod,
                        lambda x:
                            inspect.isclass(x)
                            and issubclass(x, parent_class)
                            and x.__name__ not in ignored
                        )
                )
        return classes




