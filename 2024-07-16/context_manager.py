import types

class MyContextManager:
    def __enter__(self):
        print("Entering the context")
        return self

    def __exit__(self, 
                 exc_type: type[BaseException] | None, 
                 exc_value: BaseException | None, 
                 traceback: types.TracebackType | None) -> bool:
        print("Exiting the context")
        if exc_type is not None:
            print(f"An exception occurred: {exc_value}")
        return True  # Suppress the exception

with MyContextManager() as manager:
    print("Inside the context")
    raise ValueError("I'm throwing an error here")

# Output:
# Entering the context
# Inside the context
# Exiting the context
# An exception occurred: I'm throwing an error here