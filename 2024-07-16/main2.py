# Method 2: use multiprocessing pool instead

import multiprocessing
import asyncio
import time

def test_function(x: int):
    time.sleep(1)
    result = x**2
    return result


if __name__ == "__main__":
    start = time.time()
    
    inputs = [5, 3, 9, 7, 8]
    
    with multiprocessing.Pool() as pool:
        results = pool.map(test_function, inputs)
    
    print(results)
    
    end = time.time()
    print(f"Time taken: {end - start} seconds")

