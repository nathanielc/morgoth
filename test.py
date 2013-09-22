import numpy as np
from scipy.optimize import curve_fit
def func(x, a, b, c):
    return a*np.exp(-b*x) + c
x = np.linspace(0,4,50)
y = func(x, 2.5, 1.3, 0.5)
yn = y + 0.2*np.random.normal(size=len(x))
popt, pcov = curve_fit(func, x, yn)


#print popt, pcov

def rss(x,y, model, params):
    s = 0
    for i in range(len(x)):
        s += (y[i] - model(x[i], *params))**2

    return s

print np.matrix.trace(pcov)
print rss(x, yn, func, popt)

