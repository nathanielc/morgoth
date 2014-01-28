
class Output(object):
    def __init__(self):
        pass

    def start(self):
        """ Start the output plugin """
        raise NotImplementedError("%s.start is not implemented" % self.__class__.__name__)

    def stop(self):
        """ Stop the output plugin """
        raise NotImplementedError("%s.stop is not implemented" % self.__class__.__name__)
