package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/f-secure-foundry/tamago/board/pi-foundation/pizero"
)

/*
func videoCoreInfo() {
	log.Println("-- VideoCore -------------------------------------------------------")

	log.Printf("Firmware Rev: 0x%x", videocore.FirmwareRevision())
	log.Printf("Board Model: 0x%x", videocore.BoardModel())
	log.Printf("MAC Address: %v", hex.EncodeToString(videocore.MACAddress()))
	log.Printf("Serial: 0x%x", videocore.Serial())

	start, size := videocore.CPUMemory()
	log.Printf("CPU Memory: 0x%x - 0x%x (%d MB)", start, start+size-1, size/(1024*1024))
	start, size = videocore.GPUMemory()
	log.Printf("GPU Memory: 0x%x - 0x%x (%d MB)", start, start+size-1, size/(1024*1024))

	log.Printf("DMA Channels: 0x%x", videocore.DMAChannels())
}

func display() {
	log.Println("-- Display -------------------------------------------------------")

	data := videocore.FrameBuffer.EDID()
	log.Printf("EDID data: %s", hex.EncodeToString(data))

	width, height := videocore.FrameBuffer.PhysicalSize()
	log.Printf("Physical Size: %d x %d pixels", width, height)

	log.Printf("Changing to 1024x600")
	videocore.FrameBuffer.SetPhysicalSize(1024, 600)
}

func dmaTest() {
	log.Println("-- DMA -------------------------------------------------------------")

	const dmasize = 1024 * 1024

	bufHandle := videocore.AllocateMemory(1024*1024, 16, videocore.GPU_MEMORY_FLAG_DIRECT|videocore.GPU_MEMORY_FLAG_HINT_PERMALOCK)
	dmaAddr := videocore.LockMemory(bufHandle)

	log.Printf("Allocated 0x%x for DMA transfers", dmaAddr)

	dma.Init(dmaAddr, dmasize)

	pi.DMA.Init(dma.Default())
	ch, err := pi.DMA.AllocChannel()
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
	/*
		videoCoreInfo()
		display()
		dmaTest()
	*/
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
	/*
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
