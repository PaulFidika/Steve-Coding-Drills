import concurrent.futures
import time

def test_function(x: int, y: int):
    time.sleep(1)
    result = x**2 + y
    return result

if __name__ == "__main__":
    inputs = [5, 3, 9, 7, 8]
    results = []
    
    start = time.time()
    
    with concurrent.futures.ProcessPoolExecutor() as executor:
        # Submit all tasks to the executor
        future_to_x = {executor.submit(test_function, x, 5): x for x in inputs}
        
        print(f"Submitted tasks: {future_to_x}")
        
        # Collect results as they complete
        for future in concurrent.futures.as_completed(future_to_x):
            results.append(future.result())
    
    print(results)
    
    end = time.time()
    print(f"Time taken: {end - start} seconds")
