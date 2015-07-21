# Mappings

Morgoth receives data from the scheduled queries.
Once it receives this data it needs to decide which detection instance should process the data.
This is done via simple `mappings`.

The process for finding a detector instance is:

1. Check all active detector instances for an exact match of name and all tags on the window. If found pass window to detector instance.
2. Check configured `mappings` for regex match to a new detector instance. If found create new detector instance with name and all tags on the window, and pass it the window.

A `mapping` is a name regex pattern and a set of tag regex patterns.
The measurement name must match the name regex pattern.
Each tag value must match its corresponding regex pattern.
If a tag is not specified in the mapping but exists in the window it is ignored.
Once a match is found, a new detector instance is created witht the configured detector settings.
So if a new window arrives with different data but the same name and tag set it will be routed to the same detector.

