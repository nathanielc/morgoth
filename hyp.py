
import numpy
from scipy.stats import chi2
import matplotlib.pyplot as plt

def main():
    #Inputs
    size = 100
    u = numpy.array(
        list(5 * numpy.random.rand(size * 2)) +
        list(10 * numpy.random.rand(size) + 20) +
        list(numpy.random.rand(size * 2)) +
        list(20 * numpy.random.rand(size * 2))
        )

    u = numpy.array(list(u) + list(numpy.random.rand(len(u)) + u) + list(numpy.random.rand(len(u)) + u))

    plt.plot(u)
    u_min = numpy.min(u)
    u_max = numpy.max(u) + 1
    n_bins = 20
    n = len(u)
    w = 50
    T = chi2.ppf(0.95, n_bins - 1)
    c_th = 1

    #intermediate
    m = 0
    w_i = 1
    step_size = (u_max - u_min) / float(n_bins)

    Ps = []
    cs = []
#    ws = [1, 3, 8, 15]
#    for w_i in ws:
    while w_i * w <= n:
        anon = False
        u_current = u[(w_i - 1) * w: w_i * w]
        print "Examining window %d" % w_i,
        b_current = [int((v - u_min) / step_size) for v in u_current]
        P = compute_P(b_current, n_bins)
        if m == 0:
            Ps.append(P)
            m = 1
            cs.append(1)
            print " Null"
        else:
            min_D = None
            c_i = None
            i = 0
            for i in range(len(Ps)):
                d = D(P, Ps[i])
                if 2 * w * d < T:
                    if d < min_D or min_D is None:
                        min_D = d
                        c_i = i


            if c_i is not None:
                cs[c_i] += 1

                if cs[c_i] > c_th:
                    print "NOT ANOMALOUS"
                else:
                    print "ANOMALOUS"
                    anon = True


            else:
                print "ANOMALOUS"
                anon = True
                m += 1
                Ps.append(P)
                cs.append(1)


        color = 'k'
        lsize = 1
        if anon:
            color = 'r'
            lsize = 2
        plt.vlines((w_i - 1) * w, u_min, u_max, color=color, linewidth=lsize, label=w_i)
        w_i += 1
        print

    print cs

    plt.show()

def D(q, p):
    assert len(q) == len(p)
    d = 0
    for i in range(len(q)):
        d += q[i] * numpy.log(q[i] / p[i])
    return d


def compute_P(b_current, n_bins):
    P = numpy.ones(n_bins)
    for b in b_current:
        P[b] += 1

    P /= (len(b_current) + n_bins)
    return P


main()
