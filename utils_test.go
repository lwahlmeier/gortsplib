package gortsplib

import (
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReuseBuffer(t *testing.T) {
	rb := GetReuseBuffer()
	rb.SetMaxSpare(10)
	rb.SetDefaultAlloc(64 * 1024)
	sba := make([]*SimpleBuffer, 20)
	for i := 0; i < 20; i++ {
		sba = append(sba, rb.GetBuffer())
	}
	require.Equal(t, 0, len(rb.buffers))
	sba = sba[:0]
	runtime.GC()
	runtime.GC()
	runtime.GC()
	runtime.GC()
	require.Equal(t, 10, len(rb.buffers))
}

func TestReuseBufferMulti(t *testing.T) {
	rb := GetReuseBuffer()
	rb.SetMaxSpare(100)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			for q := 0; q < 20; q++ {
				rb.GetBuffer()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// require.Equal(t, 0, len(rb.buffers))
	// sba = nil
	runtime.GC()
	runtime.GC()
	runtime.GC()
	runtime.GC()
	require.Equal(t, 100, len(rb.buffers))
}

func TestNoContentLength(t *testing.T) {
	msg := "OPTIONS rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 1\r\n" +
		"Proxy-Require: gzipped-messages\r\n" +
		"Require: implicit-play\r\n" +
		"\r\n"
	i, err := getContentLength(msg)
	require.NoError(t, err)
	require.Equal(t, 0, i)
}

func TestContentLength(t *testing.T) {
	msg := "OPTIONS rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 1\r\n" +
		"Content-Length: 12\r\n" +
		"Proxy-Require: gzipped-messages\r\n" +
		"Require: implicit-play\r\n" +
		"\r\n"
	i, err := getContentLength(msg)
	require.NoError(t, err)
	require.Equal(t, 12, i)
}

func TestContentLengthBadValue(t *testing.T) {
	msg := "OPTIONS rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 1\r\n" +
		"Content-Length: bad\r\n" +
		"Proxy-Require: gzipped-messages\r\n" +
		"Require: implicit-play\r\n" +
		"\r\n"
	i, err := getContentLength(msg)
	require.Error(t, err)
	require.Equal(t, -1, i)
}
