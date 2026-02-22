package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func slowFunc() {
	// simulate CPU work
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

var leakyData [][]byte

func leakyFunc() {
	// simulate memory allocation (1MB chunks)
	leakyData = append(leakyData, make([]byte, 1024*1024))
}

func main() {
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuProfile.Close()

	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// Spin intentionally for CPU & Memory
	for i := 0; i < 100; i++ {
		slowFunc()
		leakyFunc()
	}
	time.Sleep(100 * time.Millisecond)

	// Dump memory profile
	memProfile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memProfile.Close()
	runtime.GC() // get up-to-date heap statistics
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatal(err)
	}
}
