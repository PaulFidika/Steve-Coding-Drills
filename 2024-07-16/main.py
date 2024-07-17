# Method 1: use multiprocessing.Process, and pass in the queue as an argument

import multiprocessing
import asyncio
import time

def test_function(x: int, queue: multiprocessing.Queue):
    time.sleep(1)
    result = x**2
    queue.put(result) # place on queue rather than return

if __name__ == "__main__":
    start = time.time()
    
    queue = multiprocessing.Queue()
    processes = []
    inputs = [5, 3, 9, 7, 8]
    
    for x in inputs:
        p = multiprocessing.Process(target=test_function, args=(x,queue))
        processes.append(p)
        p.start()

    # for p in processes:
    #     p.join()
    
    results = [queue.get() for _ in range(len(processes))]
    print(results)
    
    end = time.time()
    print(f"Time taken: {end - start} seconds")

