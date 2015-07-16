#!/usr/bin/python

import sys
import os
import re
import json
import calendar
from datetime import datetime

line_pattern = r'^\w(\d{4} \d{2}:\d{2}:\d{2}\.\d{6}).*(rot[\d\.]*%s[-_\d\.]*)\|(.*)# ({.*})$'
date_pattern = '%Y%m%d %H:%M:%S.%f'

class Parser(object):
    def __init__(self, key, callback):
        self.key = key
        self.line_pattern = re.compile(line_pattern % key)
        self.callback = callback

    def parse(self):
        #Get data from logs
        year = datetime.now().year
        uniq = {}
        for line in sys.stdin:
            matches = self.line_pattern.match(line)
            if matches:

                date, ident, metric, json_data = matches.groups()
                date = datetime.strptime(str(year) + date, date_pattern)

                name = '%s.%s.%d' % (metric, ident, calendar.timegm(date.timetuple()))
                if name not in uniq:
                    uniq[name] = 0
                else:
                    uniq[name] += 1
                    name += '.' + str(uniq[name])



                #if 'total.idle' not in metric:
                #    continue

                #if not (date > datetime(2014, 11, 6, 6, 43) and date < datetime(2014, 11, 6, 6, 50)):
                #    continue
                #print date

                path = '%s/%s.pdf' % (self.key, name)
                if os.path.exists(path):
                    continue

                data = json.loads(json_data)

                self.callback(date, ident, metric, data, path)

