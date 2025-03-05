# OS_HW_4
Part 1
he ABA problem is an issue in programming that occurs with lock-free data structures using Compare-And-Swap operations. It happens when a thread reads a value, another thread modifies it, and then it is changed back to the original value before the first thread acts on it. Since CAS only checks if the value is the same as before, it mistakenly assumes nothing has changed, even though modifications did occur.

Links
https://www.baeldung.com/cs/aba-concurrency
https://www.stroustrup.com/isorc2010.pdf
https://lumian2015.github.io/lockFreeProgramming/aba-problem.html

Part 2
How to run: 
1. Place .go file in chooce of IDE
2. in terminal type: go run main.go

The benchmarking effectively measures performance differences by executing concurrent enqueue and dequeue operations across multiple goroutines, stressing synchronization mechanisms under high workloads. QueueMutex uses sync.Mutex, making it suitable for moderate concurrency but prone to contention at high loads. QueueAtomic, leveraging lock-free atomic.Pointer and Compare-And-Swap (CAS), scales best in high-concurrency scenarios by avoiding blocking, though it introduces CPU overhead. QueueLock, similar to QueueMutex, explicitly locks head and tail separately, performing similarly but with slight differences in contention handling. QueueAtomic is ideal for high-throughput environments, while QueueMutex and QueueLock work better for simpler, moderate-concurrency workloads.

AI usage: ChatGPT was used for help with how to create a benchmark segment of the code
