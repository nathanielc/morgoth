
class Input(object):
    def __init__(self):
        pass

    @classmethod
    def from_conf(cls, conf):
        """
        Create a input from the given conf

        @param conf: a conf object
        """
        raise NotImplementedError("%s.from_conf is not implemented" % cls.__name__)

    def start(self):
        """ Start collecting data """
        raise NotImplementedError("%s.start is not implemented" % self.__class__.__name__)

    def stop(self):
        """ Stop the collection of data """
        raise NotImplementedError("%s.stop is not implemented" % self.__class__.__name__)
