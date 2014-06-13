###############
Getting Started
###############

To get started using morgoth follow one of these simple exercises below
or start by reading the configuration, and components section of this documentation
to learn how to start using morgoth for your own project.

.. contents::

Setting anomaly detection of error rates for your app
=====================================================

Install Morgoth
---------------

First lets install morgoth and its dependecies.

.. literalinclude:: dependencies

Next checkout the code and add it to your PYTHONPATH.

.. code-block:: sh

   $ git clone https://github.com/nvcook42/morgoth.git
   $ cd morgoth
   $ source ./pythonpath.sh

Running Morgoth
---------------

Choose either the MongoDB or InfluxDB example configs. Modify them if necessary,
by default they point to localhost for their data store. For InfluxDB you will
have to either create a database named morgoth with user and password morgoth, or
change the config.

.. code-block:: sh

   $ cp <mongodb,influxdb>.example.yaml myproject.yaml


Once you have setup your config start up morgoth

.. code-block:: sh

   $ ./morgoth -c myproject.yaml

The process should remain in the foreground and print out some simple debug messages.


Send graphite data to Morgoth
-----------------------------

There are three ways to get graphite data into morgoth.

* Point your application at morgoth with the graphite fitting running
* Use the pull graphite config fitting periodically query data from graphite
* Configure the carbon-relay to send data to both morgoth and your existing system

For simplicity in this tutorial we will follow the first option.

Using the example config you should see that the graphite fitting is configured and listening on port 2003.

Lets send some data to morgoth:

.. code-block:: sh

   $ echo "graphite.test_data 42 `date +'%s'`" | nc 127.0.0.1 2003

Now query the api to make sure we got the data...

.. code-block:: sh

   $ curl -X GET http://localhost:7001/metrics
   {
      "metrics": [
         "graphite.test_data"
      ]
   }

   $ curl -X GET http://localhost:7001/data/graphite.test_data
   {
     "metric": "graphite.test_data",
     "data": [
       [
         "2014-06-13T04:53:14+00:00",
         42
       ]
     ]
   }



Now lets add more data and examine the dashboard to visualize the data.

.. code-block:: sh

   $ for i in {1..9}; do echo "graphite.test_data 0.$i `date +'%s'`" | nc 127.0.0.1 2003; sleep 1; done

Now navigate to http://localhost:7000 to see the 10 seconds or so of data.
(This dashboard is functional and thats all, if you have any design sinse at all I would gladly accept a pull request)

Thats it, you can add data to morgoth just like you would graphite.

Detecting your first anomaly
----------------------------


Now lets detect and anomaly. The example config has just a single Threshold anomaly detector. As it is configured if 70%
of the values during the data window are above '2' than it will consider the window anomalous.

Lets use the load average data from your current system. Run this script to add the one minute load average every second to morgoth.

.. code-block:: sh

   $ while true; do echo "graphite.load `cat /proc/loadavg | awk '{print $1}'` `date +'%s'`" | nc 127.0.0.1 2003; sleep 1; done &

Now lets create some cpu activity and see if we can't get morgoth to detect it. Run this bash script to spawn
two infinite bash loops. We will kill these later.

.. code-block:: sh

   $ for i in {1..2}; do { while true; do i=0 ; done & }; done


At this point either in the dashboard or the api you should see that the load avg metric is climbing.


.. code-block:: sh

   $ curl -X GET http://localhost:7001/data/graphite.load


Now watch the running morgoth process logs and in a few minutes you should see a critical log about an anomalous window.
In the example configuration the detectors run on a 60s schedule. The duration parameter tells the detectors to
consider a 60s window. The period parameter tells the detectors to run their alogrithm every 60s seconds. The delay option of 0
tells the detectos to consider the most recent 60s window.

The anomaly should also show up as an event along the time line in the dashboard.


Kill the backgrounded jobs:

.. code-block:: sh

   $ kill $(jobs -p)


