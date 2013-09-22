

from morgoth.reader import Reader
from datetime import datetime

r = Reader()
print r.get_data('usr', datetime(2013, 9, 20, 12))

