import time

def count_to_100_million():
    start = time.time()
    res = 0

    for i in range(1, 100_000_001):
        res += 1

    duration = time.time() - start
    print(f"Counting to 100,000,000 took {duration:.6f} seconds")

def main():
    count_to_100_million()

if __name__ == "__main__":
    main()
