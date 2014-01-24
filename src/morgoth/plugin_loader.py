import os
import imp
import sys
import inspect

class PluginLoader(object):
    """
    A class to facilitate loading plugins dynamically
    """
    def __init__(self):
        """
        """

    def find_modules(self, search_dirs):
        """
        Find all modules that exists under packages of the 
        same name in the search_dirs

        Example:
        `search_dir`/
                    foo/
                        __init__.py
                        foo.py
                    bar/
                        README.txt
                    bob.py

        Only module foo.foo will be imported

        @param search_dirs: list of directories to search in
        """
        mods = []
        for search_dir in search_dirs:
            if not os.path.isdir(search_dir): continue
            for pkg_dir in os.listdir(search_dir):
                pkg_path = os.path.join(search_dir, pkg_dir)
                if not os.path.isdir(pkg_path): continue

                #Found package importing
                found_pkg = imp.find_module(pkg_dir, [search_dir])
                if not found_pkg:
                    continue

                pkg = imp.load_module(pkg_dir, *found_pkg)

                # Find modules under package

                name = "%s.%s" % (pkg_dir, pkg_dir)
                found = imp.find_module(pkg_dir, [pkg_path])
                if not found:
                    continue
                mod = imp.load_module(name, *found)
                mods.append(mod)
        return mods
    def find_subclasses(self, mods, parent_class):
        """
        Return all classes found in the list of modules `mods` that are subclasses
        of `parent_class`

        @param mods: list of module objects
        @param parent_class: a parent class
        """
        classes = []
        for mod in mods:
            classes.extend(
                    inspect.getmembers(
                        mod,
                        lambda x:
                            inspect.isclass(x)
                            and issubclass(x, parent_class)
                            and not x == parent_class
                        )
                )
        return classes




