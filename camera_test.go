package framegrabber_test

import (
	"fmt"
	"framegrabber"
	"image/jpeg"
	"os"
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	framegrabber.Detect()
}

func TestInit(t *testing.T) {
	cam, _ := framegrabber.New("foo")
	if cam != nil {
		t.Error("Empty config is rejected")
	}

	cam, _ = framegrabber.New("example.conf")
	if cam == nil {
		t.Error("Valid config is allowed")
	}
	cam.Destroy()
}

func TestGrab(t *testing.T) {
	cam, _ := framegrabber.New("example.conf")
	defer cam.Destroy()

	frame := cam.GrabFrameWithTimeout(5 * time.Second)
	if frame == nil {
		t.Error("Proper config allows frame grabbing within 5s")
	}

	tsStart := time.Now()
	frameCount := 0
	testDuration := int(3)
	for {
		frame := cam.GrabFrameWithTimeout(1 * time.Second)
		if frame != nil {
			frameCount++
			f, _ := os.Create(fmt.Sprintf("capture_%d.jpg", frameCount))
			defer f.Close()
			jpeg.Encode(f, frame, nil)
		}
		if time.Since(tsStart) > time.Duration(testDuration)*time.Second {
			break
		}
	}
	if frameCount < testDuration {
		t.Error("Can grab at least 1 frame per second")
	}
	fmt.Printf("Fetched %d frames within %d s\n", frameCount, testDuration)

	go cam.FrameGrabberStart()
	ts := time.Now()
	tsStart = time.Now()
	frameCount = 0
	for {
		img, imgTs := cam.FrameGrabberGet()
		if img != nil {
			if ts != imgTs {
				ts = imgTs
				frameCount++
			}
		}

		if time.Since(tsStart) > time.Duration(testDuration)*time.Second {
			break
		}
	}
	if frameCount < testDuration {
		t.Error("Can grab at least 1 frame per second using framegrabber")
	}
	fmt.Printf("Fetched %d frames within %d s using framegrabber\n", frameCount, testDuration)

}
