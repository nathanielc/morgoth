

import logging

logger = logging.getLogger(__name__)

class Notifier(object):
    def __init___(self):
        pass

    def notify(self, window):
        """
        Notify that the window is anomalous
        """
        logger.info("Window: %s is anomalous" % window)

