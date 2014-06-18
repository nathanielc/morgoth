
Morgoth
=======

Morgoth is a flexible anomaly detection application.
It is designed to integrate with many existing technologies
and provide powerful anomaly detection on top of those systems.
Namely morgoth currently integrates with graphite, MongoDB and InfluxDB.
Elasticsearch integration is comming soon.

Getting Started
---------------

See the tutorial at http://nvcook42.github.io/morgoth/docs/started.html for an example
of detecting anomalies in graphite metrics.


What makes Morgoth different from other similiar projects?
----------------------------------------------------------
Morgoth takes a very modular approach to anomaly detection. There is no single algorithm
or method to detect anomalies. Morgoth comes built in with several algorithms but can
easily be extended to add more.

Morgoth also tries to be very simple. It doesn’t have tons of dependencies and most
of its dependencies can be switched out for a system of your choosing. For example
it can be run on top of either MongoDB or InfluxDB at the moment.

Morgoth does have one anomaly detection alorithm that is unique and well suited to
detecting anomalies in system level type metrics (cpu, load, etc). The algorithm is
called MGOF (Multinomial Goodness of Fit), and has been adapted from this research paper
(http://www.hpl.hp.com/techreports/2011/HPL-2011-8.html)


What is MGOF?
-------------

In a nut shell MGOF (Multinomial Goodness of Fit, as I am calling it) is an pattern
learning algorithm adept and detecting anomalies in non Gaussian data.

Using standard deviations or 3-sigma algorithms is a very common way of detecting anmalies in metric data.
These techniques assume the data they operate on follows a Gaussian distribution. Unfortunately much metric
data is not Gaussian. Take cpu usage for example. Servers perfoming work from a queue tend to be nearly idle
then spike to 100% cpu usage and then drop back down. Most of the time the server is either near 0% utilization
or 100% following a bimodal distribution, (http://en.wikipedia.org/wiki/Bimodal_distribution).

The MGOF alorithm assumes no distribution of the data. Rather the way is detects anomalies is to calculate the
distribution for different windows of time. Then compare each of those distributions to the distribution of the window in question
using a simple chi-squared test (http://en.wikipedia.org/wiki/Chi-squared_test).

In summary the MGOF algorithm is well suited for data collected from systems and applications because it doesn't assume a distribution
of the data.


