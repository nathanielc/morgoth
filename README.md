
Morgoth [![Build Status](https://travis-ci.org/nathanielc/morgoth.svg?branch=master)](https://travis-ci.org/nathanielc/morgoth)
=======

Morgoth is a flexible anomaly detection application.
It is designed to integrate with many existing technologies
and provide powerful anomaly detection on top of those systems.
Namely morgoth currently integrates with InfluxDB 0.9.x.


Getting Started
---------------

See the tutorial at http://docs.morgoth.io/docs/overview/getting-started/ for an example.


What makes Morgoth different from other similiar projects?
----------------------------------------------------------
Morgoth takes a very modular approach to anomaly detection. There is no single algorithm
or method to detect anomalies. Morgoth comes built in with several algorithms but can
easily be extended to add more.

Morgoth does have one anomaly detection alorithm that is unique and well suited to
detecting anomalies in system level type metrics (cpu, load, etc). The algorithm is
called MGOF (Multinomial Goodness of Fit), and has been adapted from this research paper
(http://www.hpl.hp.com/techreports/2011/HPL-2011-8.html)


What is MGOF?
-------------

In a nutshell MGOF (Multinomial Goodness of Fit, as I am calling it) is a pattern
learning algorithm adept at detecting anomalies in non Gaussian data.

Using standard deviations or 3-sigma algorithms is a very common way of detecting anmalies in metric data.
These techniques assume the data they operate on follows a Gaussian distribution (http://en.wikipedia.org/wiki/Normal_distribution).
Unfortunately much metric data is not Gaussian. Take cpu usage for example. Servers perfoming work from a queue tend to be nearly idle
then spike to 100% cpu usage and then drop back down. Most of the time the server is either near 0% utilization
or 100% following a bimodal distribution, (http://en.wikipedia.org/wiki/Bimodal_distribution). If cpu usage
were gaussian then the cpu would spend most of the time around 50% utilized and rarley 10% or 90% utilized.

The MGOF alorithm assumes no distribution of the data. Rather the way is detects anomalies is to calculate the
distribution for different windows of time. Then compare each of those distributions to the distribution of the window in question
using a simple chi-squared test (http://en.wikipedia.org/wiki/Chi-squared_test).

In summary the MGOF algorithm is well suited for data collected from systems and applications because it doesn't assume a distribution
of the data.


