#!/usr/bin/python

"""
Pass -v=3 to morgoth so it prints verbose logs. Feed the logs into this script and it will
create some plots around the histograms the the MGOF algo used to make its decision

NOTE: plotting requires matplotlib and numpy
"""

#Headless config
import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import numpy
from datetime import datetime
from parser import Parser
from matplotlib.backends.backend_pdf import PdfPages


def plot(date, ident, metric, data, path):
    print path
    pdf = PdfPages(path)

    min = data['current']['Min']
    max = data['current']['Max']
    current_bins = data['current']['Bins']
    nbins = len(current_bins)
    ind = numpy.linspace(min, max, nbins)
    width = (max - min) / float(nbins) * 0.4


    i = 0
    total = len(data['fingerprints'])
    for fingerprint in data['fingerprints']:
        i += 1
        plt.figure(figsize=(6, 6))
        plt.bar(ind, current_bins, width, color='b', label='current')
        plt.bar(ind+width, fingerprint['Hist']['Bins'], width, color='r', label='fingerprint')
        plt.title('Fingerprint %d of %d: Seen %d times' % (i, total, fingerprint['Count']))
        plt.xlim(min, max)
        plt.legend()
        pdf.savefig()  # saves the current figure into a pdf page
        plt.close()

    d = pdf.infodict()
    d['Title'] = 'MGOF Anomalous: %s' % data['anomalous']
    d['CreationDate'] = datetime.now()

    pdf.close()


if __name__ == '__main__':
    parser = Parser('mgof', plot)
    parser.parse()

