# Method 3: use a multiprocessing Pipe instead

import multiprocessing
from multiprocessing.connection import PipeConnection
import asyncio
import time

def test_function(x: int, conn: PipeConnection):
    time.sleep(1)
    result = x**2
    conn.send(result)
    conn.close()


if __name__ == "__main__":
    inputs = [5, 3, 9, 7, 8]
    results = []
    processes: list[multiprocessing.Process] = []
    parent_conns: list[PipeConnection] = []
    
    start = time.time()
    
    for x in inputs:
        parent_conn, child_conn = multiprocessing.Pipe()
        p = multiprocessing.Process(target=test_function, args=(10, child_conn))
        p.start()
        processes.append(p)
        parent_conns.append(parent_conn)
    
    for parent_conn in parent_conns:
        results.append(parent_conn.recv())
    
    for p in processes:
        p.join()
        p.close()
    
    print(results)
    
    end = time.time()
    print(f"Time taken: {end - start} seconds")
