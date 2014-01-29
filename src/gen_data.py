#!/usr/bin/python

import subprocess
import numpy
import time
import socket


def get_stats():
    p = subprocess.Popen(['mpstat', '1', '1'], stdout=subprocess.PIPE)
    out, err = p.communicate()
    lines = out.strip().split('\n')
    names = lines[-3].split()[3:]
    values = lines[-1].split()[2:]
    stats = []
    for i in range(len(names)):
        stats.append((names[i][1:], float(values[i])))

    return stats


def main():
    while True:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect(('127.0.0.1', 4200))
        stats = get_stats()
        for metric, value in stats:
            s.sendall("local.cpu.%s %f %d\n" % (metric, value, time.time()))
        s.close()
        time.sleep(10)


if __name__ == '__main__':
    main()
