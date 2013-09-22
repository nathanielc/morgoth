from morgoth.window import Window
from datetime import datetime


start = datetime(2013,9,20,00)
end = datetime(2013,9,21,14, 30)
m = Window(start, end)
m._update('usr')

