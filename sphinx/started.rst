###############
Getting Started
###############

To get started using Morgoth follow one of these simple exercises below
or start by reading the configuration, and components section of this documentation
to learn how to start using Morgoth for your own project.

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

There are three ways to get graphite data into Morgoth.

* Point your application at Morgoth with the graphite fitting running
* Use the pull graphite config fitting periodically query data from graphite
* Configure the carbon-relay to send data to both Morgoth and your existing system

For simplicity in this tutorial we will follow the first option.

Using the example config you should see that the graphite fitting is configured and listening on port 2003.

Lets send some data to Morgoth:

.. code-block:: sh

   $ echo "graphite.test_data 42 `date +'%s'`" | nc localhost 2003

Now query the api to make sure we got the data...

.. code-block:: sh

   $ curl -X GET http://localhost:8001/metrics
   {
      "metrics": [
         "graphite.test_data"
      ]
   }

   $ curl -X GET http://localhost:8001/data/graphite.test_data
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

   $ for i in {1..10}; do echo "graphite.test_data $i `date +'%s'`" | nc localhost 2003; sleep 1; done

Now navigate to http://localhost:4000 to see the 10 seconds or so of data.
(This dashboard is functional and thats all, if you have any design sinse at all I would gladly accept a pull request)

Thats it, you can add data to Morgoth just like you would graphite.




