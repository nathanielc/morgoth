
from gevent import monkey; monkey.patch_all()
from gevent.server import StreamServer
from gevent.baseserver import BaseServer
from gevent.event import Event
import gevent
import signal

def handle(socket, address):
    print socket, address
    socket.close()

e = Event()
e.set()
s = StreamServer(
        listener=('', 4200),
        handle=handle,
        spawn=1000
    )


def signal_handler():
    e.clear()
    print "sig handler"
    s.stop()
    e.set()

gevent.signal(signal.SIGINT, signal_handler)


gs = gevent.spawn(s.serve_forever, 1000000)

gs.join()

e.wait()
