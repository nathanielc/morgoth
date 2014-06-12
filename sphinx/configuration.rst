#############
Configuration
#############

.. contents::

Morgoth configuration is all in yaml. Morgoth provides a handy
:class:`config <morgoth.config.Config>` class that allows the config to
be easily accessed within the code base.


Features
========

* Include other yaml files relative to the current file via the !include directive
* Include other yaml files in a directory relative to the current file via the !include_dir directive
* Preserve order of yaml dictionaries
* Full support for alias nodes so large repeated configuration can be consolidated
  NOTE: Aliases do not work across files.
* Tracks which configuration is used so unused configuration can be removed

Sections
========

There are four main sections to the Morgoth configuration.


data_engine
-----------

This section of the configuration determines the data engine to be used.

.. code-block:: yaml

   data_engine:
      MongoEngine:
         <specific conf for mongo>


At the moment only one data engine may be specified.


metrics
-------

In this section you define what Morgoth will do with each of the metrics it receives.


Lets look an incomplete starting metrics config:


.. code-block:: yaml

   metrics:
      app.*:
         <specific config for just the app.* metrics>
      node.*:
         <specific config for just the node.* metrics>
      network.*:
         <specific config for just the network.* metrics>


First we can see that the metrics config is a dictionary where each key is a regex
that matches metric names. This way you can define your own metric naming scheme.
Since the matching is done via regex there is no assumed structure to the metrics names.

Since the dictionaries are ordered the regex patterns need not be mutually exlusive.
The first pattern to be matched will be used.


Within each metrics you can define several components:


.. code-block:: yaml

   metrics:
      app.*:
         detectors:
            Threshold:
               <specific threshold config>
            Tukey:
               <specific tukey config>
         notifiers:
            EmailNotifier:
               <specific email config>
         schedule:
            duration: 5m
            period: 5m
            delay: 1m
      node.*:
         detectors:
            Threshold:
               <specific threshold config>
         notifiers:
            LogNotifier:
               <specific log config>
         schedule:
            duration: 5m
            period: 5m
            delay: 1m
      network.*:
         detectors:
            Threshold:
               <specific threshold config>
            MGOF:
               - <specific mgof config>
               - <another mgof config>
         notifiers:
            EmailNotifier:
               <specific email config>
         schedule:
            duration: 5m
            period: 5m
            delay: 1m

fittings
--------

The fittings that you wish to be installed are configured here. The fittings config object
is a dictionary where each key is the name of a fitting class and the values are passed
to the fitting for instantiation.


.. code-block:: yaml

   fittings:
      Graphite:
         ...
      Rest:
         ...
      Dashboard:
         ...

plugin_dirs
-----------

Morgoth will autoload plugins from these directories. Specify a list of directories
for each type of plugin

.. code-block:: yaml

   plugin_dirs:
      detectors:
         - /etc/morgoth/detectors
      notifiers:
         - /etc/morgoth/notifiers
      fittings:
         - /etc/morgoth/fittings
      data_engines:
         - /etc/morgoth/data_engines

NOTE: The respective directories in the Morgoth code base will always be searched for plugins.

