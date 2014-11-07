#!/usr/bin/python

"""
Pass -v=3 to morgoth so it prints verbose logs. Feed the logs into this script and it will
create some plots around the histograms the the KS test used to make its decision

NOTE: plotting requires matplotlib
"""

import matplotlib.pyplot as plt
import numpy
from parser import Parser
from datetime import datetime
from matplotlib.backends.backend_pdf import PdfPages

def cdf(data):
    l = len(data)
    xs = numpy.zeros(l)
    ys = numpy.zeros(l)
    i = 0
    y = 0
    inc = 1.0 / float(l)
    for value in data:
        y += inc
        xs[i] = value
        ys[i] = y
        i += 1

    return xs, ys


def plot(date, ident, metric, data, path):
    pdf = PdfPages(path)

    current = data['current']
    current_xs, current_ys = cdf(current)


    i = 0
    total = len(data['fingerprints'])
    for fingerprint in data['fingerprints']:
        i += 1
        xs, ys = cdf(fingerprint['Data'])
        plt.figure(figsize=(6, 6))
        plt.plot(current_xs, current_ys, color='b', label='current')
        plt.plot(xs, ys, color='r', label='fingerprint')
        plt.title('Fingerprint %d of %d: Seen %d times' % (i, total, fingerprint['Count']))
        plt.ylim(0, 1)
        plt.legend()
        pdf.savefig()  # saves the current figure into a pdf page
        plt.close()

    d = pdf.infodict()
    d['Title'] = 'KS TestAnomalous: %s' % data['anomalous']
    d['CreationDate'] = datetime.now()

    pdf.close()

if __name__ == '__main__':
    parser = Parser('kstest', plot)
    parser.parse()
