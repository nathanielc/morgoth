# Detection Framework

There are two basic components to the detection framework Morgoth employs.

1. A Lossy Counting strategy for finding anomlous behavior.
2. A fingerprinting mechanism to summarize/fingperprint behaviors.

At a high level Morgoth fingperprints each window of data it sees and keeps track of which fingperprints it has seen before and marks new/infrequent fingperprints as anomalous.
The power of the framework is in which fingerprinting algorithms are used and in the simplicity of configuring how to count frequent fingperprints.


## Lossy Counting

The [Lossy Counting algorithm](http://www.vldb.org/conf/2002/S10P03.pdf) is a way of counting frequent items efficiently.
The algorithm is lossy since it will drop infrequent items but does so in such a way that it can guarantee certain behaviors:

1. There are no false negatives. The frequency of an item cannot be over estimated.
2. False positives are guaranteed to have a frequency of at least `sN-eN`, where `N` is the number of items processed, `s` is the minimum support and `e` is the error tolerance.
3. The frequency of an item can be underestimated by at most `eN`.
4. The space requirements of the algorithm are `1/e log(eN)`.

These constraints allow the user to have intuitive control over what is considered an anomaly.
For example `#3` states that items can be underestimated by at most `eN`.
What this means given an `e=0.10` that items that are less than 10% frequent might get dropped and when they occur again could be marked as anomalies.

