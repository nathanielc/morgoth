# Configuration


# Sections

The Morgoth configuration contains the following sections:

* [Engine](engine.md)
* [Schedules](schedules.md)
* [Mappings](mappings.md)
* [Alerts](alerts.md)


# Defaults

Morgoth internally uses yaml marshaling to read the configuration.
What this means is that defaults are stored as golang [tags](https://golang.org/ref/spec#Struct_types) on the configuration structs.
While this doesn't serve as extensive API documentation it is relatively quick to determine what configuration is being used where in Morgoth.
More verbose logging will print out each time a default value is used.
Later I plan to leverage this feature to make the morgoth binary capable of generating the default configuration contained within itself.

