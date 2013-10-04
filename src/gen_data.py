#!/usr/bin/python

import subprocess
import numpy
import time
from morgoth.collector import Collector
from datetime import datetime


def get_stats():
    p = subprocess.Popen(['mpstat', '1', '1'], stdout=subprocess.PIPE)
    out, err = p.communicate()
    lines = out.strip().split('\n')
    names = lines[-3].split()[3:]
    values = lines[-1].split()[2:]
    print values
    stats = []
    for i in range(len(names)):
        stats.append((names[i][1:], values[i]))

    return stats


def main():
    c = Collector()
    while True:
        stats = get_stats()
        for metric, value in stats:
            c.insert(datetime.utcnow(), metric, value)
        time.sleep(1)

if __name__ == '__main__':
    main()
