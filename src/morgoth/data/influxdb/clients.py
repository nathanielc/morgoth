

from influxdb import client


class InfluxDBClients(object):
    Normal = client.InfluxDBClient('localhost', 8099, 'root', 'root', 'morgoth')
