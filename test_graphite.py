#!/usr/bin/python

import os
import socket
import time
from multiprocessing import Process
from threading import Thread

def t_worker(name, thread):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect(('127.0.0.1', 4200))
    for i in xrange(1000):
        s.sendall("test_graphite.metric.%s.t%d %d %d\n" % (name, thread, i * 100, time.time()))
        time.sleep(1)
    s.close()

def worker(name):
    ts = []
    for i in xrange(100):
        t = Thread(target=t_worker, args=(name, i))
        t.start()
        ts.append(t)
    for t in ts:
        t.join()

ps = []
for i in range(2):
    p = Process(target=worker, args=("p%d" % i,))
    p.start()
    ps.append(p)

for p in ps:
    p.join()

