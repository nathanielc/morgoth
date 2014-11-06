#!/usr/bin/python

"""
Pass -v=3 to morgoth so it prints verbose logs. Feed the logs into this script and it will
create some plots around the histograms the the MGOF algo used to make its decision

NOTE: plotting requires matplotlib
"""

import matplotlib.pyplot as plt
import sys
import os
import re
import json
import numpy
import calendar
from datetime import datetime
from matplotlib.backends.backend_pdf import PdfPages

mgof_pattern = re.compile(r'^\w(\d{4} \d{2}:\d{2}:\d{2}\.\d{6}).*(rot[\d\.]*mgof[_\d\.]*)\|(.*)# ({.*})$')
date_pattern = '%Y%m%d %H:%M:%S.%f'

if __name__ == '__main__':
    #Get data from logs
    year = datetime.now().year
    for line in sys.stdin:
        matches = mgof_pattern.match(line)
        if matches:

            date, ident, metric, json_data = matches.groups()
            date = datetime.strptime(str(year) + date, date_pattern)

            name = '%s.%s.%d' % (metric, ident, calendar.timegm(date.timetuple()))
            if 'total.idle' not in metric:
                continue


            path = 'mgof/%s.pdf' % name
            if os.path.exists(path):
                continue

            print name

            pdf = PdfPages(path)
            data = json.loads(json_data)

            min = data['current']['Min']
            max = data['current']['Max']
            current_bins = data['current']['Bins']
            print sum(current_bins)
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
