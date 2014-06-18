
from distutils.core import setup

setup(name='Morgoth',
        version='0.0.1',
        description='Metric Anomaly Detection',
        author='Nathaniel Cook',
        author_email='nvcook42@gmail.com',
        url='http://nvcook42.github.io/morgoth/',
        package_dir={'' : 'src'},
        packages=[
            'morgoth',
            'morgoth.data',
            'morgoth.data.influx',
            'morgoth.data.mongodb',
            'morgoth.detectors',
            'morgoth.detectors.mgof',
            'morgoth.fittings',
            'morgoth.fittings.dashboard',
            'morgoth.notifiers',
        ],
        package_data={
            'morgoth.data.mongodb' : ['*.js'],
            'morgoth.fittings.dashboard' : ['static/*', 'templates/*'],
        },
        scripts=['morgoth']
     )

