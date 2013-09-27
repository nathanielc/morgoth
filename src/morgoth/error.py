
class MorgothError(Exception):
    def __init__(self, *args, **kwargs):
        super(MorgothError, self).__init__(*args, **kwargs)

