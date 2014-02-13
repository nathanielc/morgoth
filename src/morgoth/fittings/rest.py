#
# Copyright 2014 Nathaniel Cook
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
from flask import Flask, request, jsonify, make_response, current_app
from functools import update_wrapper

from gevent import queue
from gevent.pywsgi import WSGIServer

from datetime import timedelta
import dateutil.parser

from morgoth.fittings.fitting import Fitting
from morgoth.data.reader import Reader
from morgoth.utils import timedelta_from_str

import logging
logger = logging.getLogger(__name__)

app = Flask(__name__)

class Rest(Fitting):
    def __init__(self, host, port):
        super(Rest, self).__init__()
        self._host = host
        self._port = port


    @classmethod
    def from_conf(cls, conf):
        host = ''
        port = conf.get('port', 8080)
        return Rest(host, port)


    def start(self):
        logger.info("Starting REST fitting plugin...")
        self._server = WSGIServer((self._host, self._port), app, log=None)
        self._server.serve_forever()

    def stop(self):
        self._server.stop()

##############################
# Flask methods
##############################

reader = Reader()



def crossdomain(origin=None, methods=None, headers=None,
                max_age=21600, attach_to_all=True,
                automatic_options=True):
    if methods is not None:
        methods = ', '.join(sorted(x.upper() for x in methods))
    if headers is not None and not isinstance(headers, basestring):
        headers = ', '.join(x.upper() for x in headers)
    if not isinstance(origin, basestring):
        origin = ', '.join(origin)
    if isinstance(max_age, timedelta):
        max_age = max_age.total_seconds()

    def get_methods():
        if methods is not None:
            return methods

        options_resp = current_app.make_default_options_response()
        return options_resp.headers['allow']

    def decorator(f):
        def wrapped_function(*args, **kwargs):
            if automatic_options and request.method == 'OPTIONS':
                resp = current_app.make_default_options_response()
            else:
                resp = make_response(f(*args, **kwargs))
            if not attach_to_all and request.method != 'OPTIONS':
                return resp

            h = resp.headers

            h['Access-Control-Allow-Origin'] = origin
            h['Access-Control-Allow-Methods'] = get_methods()
            h['Access-Control-Max-Age'] = str(max_age)
            if headers is not None:
                h['Access-Control-Allow-Headers'] = headers
            return resp

        f.provide_automatic_options = False
        return update_wrapper(wrapped_function, f)
    return decorator


@app.route('/metrics')
@crossdomain(origin='*')
def metrics():
    pattern = None
    if 'pattern' in request.args:
        pattern = request.args['pattern']
    return jsonify({'metrics' : reader.get_metrics(pattern)})

@app.route('/data/<metric>')
@crossdomain(origin='*')
def metric_data(metric):
    start = None
    stop = None
    step = None
    if 'start' in request.args:
        try:
            start = dateutil.parser.parse(request.args['start'])
        except Exception as e:
            return jsonify({
                'error' : 'invalid start date format: %s' % str(e)
                }), 400

    if 'stop' in request.args:
        try:
            stop = dateutil.parser.parse(request.args['stop'])
        except Exception as e:
            return jsonify({
                'error' : 'invalid stop date format: %s' % str(e)
                }), 400

    if 'step' in request.args:
        try:
            step = timedelta_from_str(request.args['step'])
        except Exception as e:
            return jsonify({
                'error' : 'invalid step format: %s' % str(e)
                }), 400

    return jsonify({
        'metric' : metric,
        'data': reader.get_data(metric, start, stop, step)
    })

@app.route('/anomalies/<metric>')
@crossdomain(origin='*')
def metric_anomalies(metric):
    start = None
    stop = None
    if 'start' in request.args:
        try:
            start = dateutil.parser.parse(request.args['start'])
        except Exception as e:
            return jsonify({
                'error' : 'invalid start date format: %s' % str(e)
                }), 400

    if 'stop' in request.args:
        try:
            stop = dateutil.parser.parse(request.args['stop'])
        except Exception as e:
            return jsonify({
                'error' : 'invalid stop date format: %s' % str(e)
                }), 400


    return jsonify({
        'metric' : metric,
        'anomalies': reader.get_anomalies(metric, start, stop)
    })
