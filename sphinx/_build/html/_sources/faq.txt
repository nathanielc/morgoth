###
FAQ
###

.. contents::
   :depth: 2

Why anomaly detection?
======================

We are in the era of data. Everyone has data everyone tries to use their data
to make better decisions. Recently new open source projects have really taken off
that allow applications to easily collect massive amount of data about the way their systems
are behaving (see `elasticsearch <http://elasticsearch.org>`_, `logstash <http://logstash.net>`_,
`graphite <http://graphite.wikidot.com>`_, many more...).
But collecting the data is not enough, Morgoth attempts to provide a flexible
system to consume application data and make inteligent decisions about the data. What Morgoth does with
the data is up to you. It can simply notify your team of discovered anomalies or reintegrate
with the application to help it become self healing.

What makes Morgoth different from other similiar projects?
==========================================================

Morgoth takes a very modular approach to anomaly detection. There is no
single algorithm or method to detect anomalies. Morogth comes built in with several
algorithms but can easily be extended to add more.

Morgoth also tries to be very simple. It doesn't have tons of dependencies and most
of its dependencies can be switched out for a system of your choosing. For example it
can be run on top of either MongoDB or InfluxDB at the moment.

Morgoth does have one anomaly detection alorithm that is unique and well suited to
detecting anomalies in system level type metrics (cpu, load, etc). The algorithm is called
MGOF (Multinomial Goodness of Fit), and has been adapted from this research
`paper <http://www.hpl.hp.com/techreports/2011/HPL-2011-8.html>`_.


Can Morgoth integrate with ...?
================================

Morgoth can integrate with almost anything. Each component of Morgoth is plugable. Currently Morgoth
has several :doc:`fittings <fittings>` that integrate with graphite.


Why the name 'Morgoth'?
=======================

The name Morgoth comes from the Lord of the Ring folklore. Morgoth is Sauron's master. Morgoth has been written
to replace an internal system named Sauron. This older system was called Sauron because it 'saw all' metrics in our
systems. Morgoth is an attempt and improving on that design and is a complete rewrite of Sauron.

