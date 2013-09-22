from morgoth.collector import Collector
from datetime import datetime

c = Collector()

dt = datetime.now()
c.insert(dt, "morgoth.test.asdf", 42)
