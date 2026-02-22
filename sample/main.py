import time

def slow_function():
    print("Starting slow function...")
    # Simulate an expensive CPU bound operation
    total = 0
    for i in range(10_000_000):
        total += i
    return total

def main():
    slow_function()
    time.sleep(0.5)
    print("Finished.")

if __name__ == "__main__":
    main()
