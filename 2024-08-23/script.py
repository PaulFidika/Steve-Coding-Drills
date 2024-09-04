import sys

def process_request(request):
    # Example processing: just echo the request back
    return f"Processed: {request}"

def main():
    while True:
        # Read input from stdin
        request = sys.stdin.readline().strip()
        if not request:
            break

        # Process the request
        response = process_request(request)

        # Write the response to stdout
        print(response)
        sys.stdout.flush()

if __name__ == "__main__":
    main()