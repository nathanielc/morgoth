
import socket
import time
import random
import math

def send(conn, name, value):
    s = "%s %f %d\n" % (name, value, time.time())
    conn.send(s)

class offset(object):
    def __init__(self):
        self.offset = random.random() * 2
        print "New offest", self.offset

    def __call__(self, value):
        return self.offset + value


class scale(object):
    def __init__(self):
        self.scale = random.random() * 2
        print "New scale", self.scale

    def __call__(self, value):
        return self.scale * value


effects = [
        offset,
        scale
]

steps = random.randint(3, 10)
bases = []
for i in range(steps):
    bases.append((random.random()* 100, random.random() * 60 + 20))


def base(t):
    total = sum([p[1] for p in bases])
    t = t % total
    current = 0
    for v, stop in bases:
        current += stop
        if t < current:
            return v + random.random() * 3



def main():
    conn = socket.socket()
    conn.connect(('localhost', 2003))


    expire = 0


    while True:
        r = base(time.time())
        if time.time() > expire:
            expire = time.time() + random.random() * 120 + 10
            if random.random() < 0.2:
                effect = random.choice(effects)()
                print time.time(), " New effect till", expire
            else:
                effect = None

        if effect:
            r = effect(r)

        if r > 100:
            r = 100
        elif r < 0:
            r = 0
        send(conn, 'cpu.c1', r)
        time.sleep(1)

if __name__ == '__main__':
    while True:
        try:
            main()
        except:
            time.sleep(1)
            pass

