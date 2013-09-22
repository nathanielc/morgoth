

import time
import random
from multiprocessing import Process

def spike():
    while True:
        count = int(random.random() * 10**9)
        print "counting too", count
        while count > 0:
            count -= 1
        s = random.random()* 10**3
        print "sleeping for", s
        time.sleep(s)

for i in range(8):
    p = Process(target=spike)
    p.start()


