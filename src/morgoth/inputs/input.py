
class Input(object):
    def __init__(self):
        pass

    def start(self):
        """ Start collecting data """
        raise NotImplementedError("%s.start is not implemented" % self.__class__.__name__)

    def stop(self):
        """ Stop the collection of data """
        raise NotImplementedError("%s.stop is not implemented" % self.__class__.__name__)
