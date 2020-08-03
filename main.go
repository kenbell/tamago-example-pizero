package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"runtime"
	"time"

	"github.com/f-secure-foundry/tamago/pi"
	"github.com/f-secure-foundry/tamago/pi/pizero"
)

func main() {
	log.Println("Hello World!")

	log.Println("-- rng -------------------------------------------------------------")

	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	log.Printf("random bytes %s", hex.EncodeToString(b))

	size := 32

	for i := 0; i < 10; i++ {
		rng := make([]byte, size)
		rand.Read(rng)
		log.Printf("%x", rng)
	}

	count := 1000
	start := time.Now()

	for i := 0; i < count; i++ {
		rng := make([]byte, size)
		rand.Read(rng)
	}

	log.Printf("retrieved %d random bytes in %s", size*count, time.Since(start))

	log.Println("-- timer -------------------------------------------------------------")

	t := time.NewTimer(time.Second)
	log.Printf("waking up timer after %v", time.Second)

	start = time.Now()

	for now := range t.C {
		log.Printf("woke up at %d (%v)", now.Nanosecond(), now.Sub(start))
		break
	}

	log.Println("-- RAM ---------------------------------------------------------------")

	// Check GC is working by forcing more total allocation than available
	allocateAndWipe(400)
	runtime.GC()
	allocateAndWipe(400)

	log.Println("-- watchdog ----------------------------------------------------------")

	log.Println("Starting watchdog at 1s")

	// Auto-reset after 1 sec
	pi.Watchdog.Start(time.Second)
	time.Sleep(600 * time.Millisecond)
	log.Printf("Watchdog Remaining after 600ms: %v, resetting", pi.Watchdog.Remaining())

	pi.Watchdog.Reset()
	time.Sleep(600 * time.Millisecond)
	log.Printf("Watchdog Remaining after 600ms: %v", pi.Watchdog.Remaining())

	pi.Watchdog.Stop()
	log.Print("Watchdog stopped, waiting for 2 sec")
	time.Sleep(2 * time.Second)

	log.Println("-- LED ---------------------------------------------------------------")

	log.Println("Flashing the activity LED")

	board := pizero.Board

	l := board.LEDs()[0]

	ledOn := false
	l.Init()
	for i := 0; i < 40; i++ {
		time.Sleep(250 * time.Millisecond)
		ledOn = !ledOn
		l.SetState(ledOn)
	}

	log.Println("-- DONE --------------------------------------------------------------")
}

func allocateAndWipe(count int) {
	log.Printf("allocating %dMB", count)

	hold := make([][]byte, 0, 400)
	for i := 0; i < cap(hold); i++ {
		mem := make([]byte, 1024*1024)
		if len(mem) == 0 {
			break
		}
		hold = append(hold, mem)
	}

	log.Println("wiping allocation with 0xff")

	for i := 0; i < len(hold); i++ {
		for j := range hold[i] {
			hold[i][j] = 0xff
		}
	}
}
