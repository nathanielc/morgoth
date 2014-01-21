

ADS = {}

def register_ad(name, ad_class):
    """
    Register an AD by a given name

    @param name: the name to identify the AD
    @param ad_class: the class of the AD
    """
    if name in ADS:
        raise ValueError('AD of name "%s" already exists' % name)
    ADS[name] = ad_class

def get_ad(name):
    """
    """
