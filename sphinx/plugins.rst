################
Plugin Framework
################

Morgoth is implemented in a simple plugin framework. Each of its components 
follow similiar patterns. Plugins are first class citizens as the core functionality 
of morgoth is implemented as plugins.


Writting a plugins
==================

Lets look at the basics of writting a plugin to morgoth. There are three pieces
to creating a plugin. 

* The class to extend
* The configuration for the plugin
* The location to place the plugin


The class to extend
-------------------


Each type of plugin in morgoth will require that you extend a base class.
For example to create a new detector simply write a python class that impletements
the :class:`Detector <morgoth.detectors.detector.Detector>` class.

The plugin configuration
------------------------


Each plugin must implement a 'from_conf' classmethod that will return an instance
of the plugin based on a given configuration object.

For example, the email notifier from method looks like this (check out the source):

.. py:module:: morgoth.notifiers.email_notifier
.. autoclass:: EmailNotifier

   .. automethod:: from_conf






As you can see the method simply consumes a config object and returns and instance of the
EmailNotifier plugin. This allows each plugin to be fully configurable and define sane
defaults.


The location to place the plugin
--------------------------------

If a plugin is defined in the main configuration then it will be autoloaded at startup.
In order for morgoth to find the plugin is need to be in a specific location:

* In the approriate directory within to code base
* In a user configured directory

