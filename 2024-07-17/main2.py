import concurrent.futures
import time

start = time.perf_counter()

def do_something(amount: int):
    print(f'Sleeping {amount} seconds')
    time.sleep(amount)
    return amount

with concurrent.futures.ThreadPoolExecutor() as executor:
    secs = [5, 4, 3, 2, 1]
    # This returns a future -> integer
    results = [executor.submit(do_something, sec) for sec in secs]
    
    # thi blocks until all futures are resolved
    for f in concurrent.futures.as_completed(results):
        print(f.result())
    
    # =========
    
    # This returns an integer
    results = executor.map(do_something, secs)
    
    # this blocks and does nothing until the entire results array is complete
    for result in results:
        print(result)

end = time.perf_counter()

print(f'{end - start:.2f} seconds')
