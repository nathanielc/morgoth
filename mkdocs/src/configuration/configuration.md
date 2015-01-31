# Configuration


# Sections

The Morgoth configuration contains the following sections:

* [Engine](engine.md)
* [Schedule](schedule.md)
* [Metrics](metrics.md)
* [Fittings](fittings.md)
* [Morgoth](morgoth.md)


# Defaults

Morgoth internally uses yaml marshaling to read the configuration.
What this means is that defaults are stored as golang [tags](https://golang.org/ref/spec#Struct_types) on the
configuration structs. While this doesn't serve as extensive API
documentation it is relatively quick to determine what configuration is
being used where in Morgoth. Later I plan to leverage this feature to
make the morgoth binary capable of generating the default configuration contianed
within itself.


