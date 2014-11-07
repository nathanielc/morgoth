#!/usr/bin/python

"""
Pass -v=3 to morgoth so it prints verbose logs. Feed the logs into this script and it will
create some plots around the histograms the the KS test used to make its decision

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

mgof_pattern = re.compile(r'^\w(\d{4} \d{2}:\d{2}:\d{2}\.\d{6}).*(rot[\d\.]*kstest[_\d\.]*)\|(.*)# ({.*})$')
date_pattern = '%Y%m%d %H:%M:%S.%f'

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



if __name__ == '__main__':
    #Get data from logs
    year = datetime.now().year
    for line in sys.stdin:
        matches = mgof_pattern.match(line)
        if matches:

            date, ident, metric, json_data = matches.groups()
            date = datetime.strptime(str(year) + date, date_pattern)

            name = '%s.%s.%d' % (metric, ident, calendar.timegm(date.timetuple()))
            if 'MemFree' not in metric:
                continue

            #if not (date > datetime(2014, 11, 6, 6, 43) and date < datetime(2014, 11, 6, 6, 50)):
            #    continue
            #print date

            path = 'kstest/%s.pdf' % name
            if os.path.exists(path):
                continue

            print name

            pdf = PdfPages(path)
            data = json.loads(json_data)

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
