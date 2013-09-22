#!/usr/bin/python

import numpy as np
import matplotlib.pyplot as plt
from scipy.optimize import curve_fit


top = 30
t = np.linspace(0, top, int(top * 10 * np.pi))
y = np.sin(t)

print len(t)

#plt.scatter(t,y)
#plt.show()

p = 10
stop = 0.00

def f0(t, a):
    return a

def f1(t, a, b):
    return a*t + b

def f2(t, a, b, c):
    return a*t*t + b*t + c

def f3(t, a, b, c, d):
    return a*t*t*t + b*t*t + c*t + d


basis = [f0, f1, f2, f3]


def detect_change_points(t, y, basis):
    new_cp = find_candidate(t, y, 0, basis)
    change_points = []
    candidates = []

    t1, t2 = get_new_time_range(t, y, change_points, new_cp)
    prev_lh = -1
    while True:
        c1 = find_candidate(*t1, basis=basis)
        c2 = find_candidate(*t2, basis=basis)
        candidates.append(c1)
        candidates.append(c2)
        print c1
        print c2
        new_cp = sorted(candidates, key=lambda c: c[1])[0]
        lh = new_cp[1]
        print "Likelihood: %f" % lh
        if ((prev_lh - lh) / prev_lh) < stop:
            break
        prev_lh = lh

        change_points.append(new_cp)
        t1 ,t2 = get_new_time_range(t, y, change_points, new_cp)
    return change_points

def get_new_time_range(t, y, change_points, new_cp):
    cps =  sorted(change_points + [new_cp], key=lambda cp: cp[2])[:2]
    ts = []

    if len(cps) == 1:
        base = cps[0][0]
        ts.append((t[:base], y[:base], 0))
        ts.append((t[base:], y[base:], base))

    else:
        for cp in cps:
            ti = t[cp[0]:cp[1]]
            yi = y[cp[0]:cp[1]]
            base = cp[0]
            ts.append((ti, yi, base))


    print ts[0][2]
    print ts[1][2]
    return ts



def find_candidate(t, y, base, basis):
    optimal_likelihood = None
    split = None
    print p, len(t) - p
    for i in range(p, len(t) - p):
        likelihood = find_likelihood(t[:i], y[:i], basis) + \
            find_likelihood(t[i + 1:], y[i + 1:], basis)

        if likelihood < optimal_likelihood or optimal_likelihood is None:
            split = base + i, likelihood, t[i] - t[0]
            optimal_likelihood = likelihood
    return split


def find_likelihood(t, y, basis):
    min_likelihood = None
    likelihood = None
    for model in basis:
        likelihood = fit(t, y, model)
        if likelihood < min_likelihood or min_likelihood is None:
            min_likelihood = likelihood

    return likelihood

def fit(t, y, model):
    popt, pcov = curve_fit(model, t, y)
    return rss(t, y, model, popt)

def rss(t, y, model, params):
    s = 0
    for i in range(len(t)):
        s += (y[i] - model(t[i], *params)) ** 2
    return s


print detect_change_points(t, y, basis)

