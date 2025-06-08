# Method one: basic thread creation

import threading
import time

def do_something(amount: int):
    print('starting now')
    time.sleep(amount)
    print('done sleeping')

t1 = threading.Thread(target=do_something, args=(2,))
t2 = threading.Thread(target=do_something, args=(5,))

start = time.perf_counter()

t1.start()
t2.start()

print('about to join')
t1.join()
print('in the middle')
t2.join()

finish = time.perf_counter()

print(f'Finished in {finish - start:.2f} second(s)')
