package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"runtime"
	"time"

	pi "github.com/f-secure-foundry/tamago/board/raspberrypi"
	"github.com/f-secure-foundry/tamago/board/raspberrypi/pizero"
)

func rng() {
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
}

func timer() {
	log.Println("-- timer -------------------------------------------------------------")

	t := time.NewTimer(time.Second)
	log.Printf("waking up timer after %v", time.Second)

	start := time.Now()

	for now := range t.C {
		log.Printf("woke up at %d (%v)", now.Nanosecond(), now.Sub(start))
		break
	}
}

func ram() {
	log.Println("-- RAM ---------------------------------------------------------------")

	// Check GC is working by forcing more total allocation than available
	allocateAndWipe(400)
	runtime.GC()
	allocateAndWipe(400)
}

func watchdog() {
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
}

/*
func videoCoreInfo() {
	log.Println("-- VideoCore -------------------------------------------------------")

	log.Printf("Firmware Rev: 0x%x", bcm2835.FirmwareRevision())
	log.Printf("Board Model: 0x%x", bcm2835.BoardModel())
	log.Printf("MAC Address: %v", hex.EncodeToString(bcm2835.MACAddress()))
	log.Printf("Serial: 0x%x", bcm2835.Serial())

	start, size := bcm2835.CPUMemory()
	log.Printf("CPU Memory: 0x%x - 0x%x (%d MB)", start, start+size-1, size/(1024*1024))
	start, size = bcm2835.GPUMemory()
	log.Printf("GPU Memory: 0x%x - 0x%x (%d MB)", start, start+size-1, size/(1024*1024))

	log.Printf("DMA Channels: 0x%x", bcm2835.DMAChannels())
}

func display() {
	log.Println("-- Display -------------------------------------------------------")

	data := bcm2835.FrameBuffer.EDID()
	log.Printf("EDID data: %s", hex.EncodeToString(data))

	width, height := bcm2835.FrameBuffer.PhysicalSize()
	log.Printf("Physical Size: %d x %d pixels", width, height)

	log.Printf("Changing to 1024x600")
	bcm2835.FrameBuffer.SetPhysicalSize(1024, 600)
}

func dmaTest() {
	log.Println("-- DMA -------------------------------------------------------------")

	const dmasize = 1024 * 1024

	bufHandle := bcm2835.AllocateMemory(1024*1024, 16, bcm2835.GPU_MEMORY_FLAG_DIRECT|bcm2835.GPU_MEMORY_FLAG_HINT_PERMALOCK)
	dmaAddr := bcm2835.LockMemory(bufHandle)

	log.Printf("Allocated 0x%x for DMA transfers", dmaAddr)

	dma.Init(dmaAddr, dmasize)

	bcm2835.DMA.Init(dma.Default())
	ch, err := bcm2835.DMA.AllocChannel()
	if err != nil {
		log.Fatalf("failed to allocate DMA channel: %v", err)
	}
	log.Printf("Channel1 Debug: %s", ch.DebugInfo().String())
	log.Printf("Channel1 Status: %s", ch.Status().String())

	srcAddr, srcBuf := dma.Reserve(16, 64)
	dstAddr, dstBuf := dma.Reserve(16, 64)
	for i := 0; i < 16; i++ {
		srcBuf[i] = byte(i & 0xff)
		dstBuf[i] = 0xa0
	}

	log.Printf("DMA from 0x%x (%p) to 0x%x (%p)", srcAddr, srcBuf, dstAddr, dstBuf)

	log.Printf("ChannelDebug: %s", ch.DebugInfo().String())
	log.Printf("ChannelStatus: %s", ch.Status().String())
	ch.CopyRAMToRAM(srcAddr, 6, dstAddr)
	log.Printf("ChannelDebug: %s", ch.DebugInfo().String())
	log.Printf("ChannelStatus: %s", ch.Status().String())

	for i := 0; i < 10; i++ {
		log.Printf("%d: 0x%x", i, dstBuf[i])
	}
}
*/
func main() {
	log.Println("Hello World!")

	rng()
	timer()
	ram()
	watchdog()

	/*
		videoCoreInfo()
		display()
		dmaTest()
	*/

	log.Println("-- LED ---------------------------------------------------------------")

	log.Println("Flashing the activity LED")

	board := pizero.Board

	ledOn := false
	for {
		time.Sleep(250 * time.Millisecond)
		ledOn = !ledOn
		board.LED("activity", ledOn)
	}
}

func allocateAndWipe(count int) {
	log.Printf("allocating %dMB", count)

	hold := make([][]byte, 0, count)
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
