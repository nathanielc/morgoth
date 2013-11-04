
import logging

def init():
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.DEBUG)
    ch = logging.StreamHandler()
    formatter = logging.Formatter('[%(asctime)s|%(name)s:%(lineno)d][%(levelname)s] %(message)s')
    ch.setFormatter(formatter)
    root_logger.addHandler(ch)

    print "Initialized Logging"


