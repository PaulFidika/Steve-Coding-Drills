# Claude 3.5 Sonnet thinks that ThreadPoolExecutor helps here
# which it doesn't of course, beause time.sleep is a synchronous
# thread-blocking route

from concurrent.futures import ThreadPoolExecutor
import threading
import time
from concurrent.futures import ProcessPoolExecutor

lock = threading.Lock()
shared_data = []
    
def task(n):
    with lock:
        time.sleep(1)
        shared_data.append(n * n)

if __name__ == "__main__":
    start_time = time.time()
    
    with ThreadPoolExecutor(max_workers=4) as executor:
        futures = [executor.submit(task, i) for i in range(10)]
        for future in futures:
            future.result()  # Wait for all tasks to complete
            
    end_time = time.time()
    execution_time = end_time - start_time

    print(f"Execution time: {execution_time:.4f} seconds")

# if __name__ == "__main__":
#     start_time = time.time()

#     with ProcessPoolExecutor(max_workers=4) as executor:
#         futures = [executor.submit(task, i) for i in range(10)]
#         for future in futures:
#             future.result()  # Wait for all tasks to complete

#         end_time = time.time()
#         execution_time = end_time - start_time

#         print(shared_data)

#         print(f"Execution time: {execution_time:.4f} seconds")
