import multiprocessing
import time

def worker(shared_list, index):
    time.sleep(index % 2)  # variability in execution time
    shared_list.append(index)
    print(f"Process {index} appended {index}")

if __name__ == "__main__":
    manager = multiprocessing.Manager()
    shared_list = manager.list()
    
    processes = []
    for i in range(10):
        p = multiprocessing.Process(target=worker, args=(shared_list, i))
        p.start()
        processes.append(p)
    
    for p in processes:
        p.join()
    
    print("Final shared list:", list(shared_list))
