class FileManager:
    def __init__(self, filename, mode):
        self.filename = filename
        self.mode = mode
        self.file = None

    def __enter__(self):
        self.file = open(self.filename, self.mode)
        return self.file  # Return the file object

    def __exit__(self, exc_type, exc_value, traceback):
        if self.file:
            self.file.close()
            
        if exc_type is not None:
            print(f"An exception occurred: {exc_value}")
            
        return False  # Do not suppress exceptions

# Using the custom FileManager context manager
with FileManager('example.txt', 'w') as file:
    file.write('Hello, World!')
    raise ValueError("I'm throwing an error here")
    
# The file is automatically closed here