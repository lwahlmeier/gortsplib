package gortsplib

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var casesRequest = []struct {
	name string
	byts []byte
	req  *Request
}{
	{
		"options",
		[]byte("OPTIONS rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
			"CSeq: 1\r\n" +
			"Proxy-Require: gzipped-messages\r\n" +
			"Require: implicit-play\r\n" +
			"\r\n"),
		&Request{
			Method: "OPTIONS",
			Url:    &url.URL{Scheme: "rtsp", Host: "example.com", Path: "/media.mp4"},
			Header: Header{
				"CSeq":          []string{"1"},
				"Require":       []string{"implicit-play"},
				"Proxy-Require": []string{"gzipped-messages"},
			},
		},
	},
	{
		"describe",
		[]byte("DESCRIBE rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
			"CSeq: 2\r\n" +
			"\r\n"),
		&Request{
			Method: "DESCRIBE",
			Url:    &url.URL{Scheme: "rtsp", Host: "example.com", Path: "/media.mp4"},
			Header: Header{
				"CSeq": []string{"2"},
			},
		},
	},
	{
		"announce",
		[]byte("ANNOUNCE rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
			"CSeq: 7\r\n" +
			"Content-Length: 306\r\n" +
			"Content-Type: application/sdp\r\n" +
			"Date: 23 Jan 1997 15:35:06 GMT\r\n" +
			"Session: 12345678\r\n" +
			"\r\n" +
			"v=0\n" +
			"o=mhandley 2890844526 2890845468 IN IP4 126.16.64.4\n" +
			"s=SDP Seminar\n" +
			"i=A Seminar on the session description protocol\n" +
			"u=http://www.cs.ucl.ac.uk/staff/M.Handley/sdp.03.ps\n" +
			"e=mjh@isi.edu (Mark Handley)\n" +
			"c=IN IP4 224.2.17.12/127\n" +
			"t=2873397496 2873404696\n" +
			"a=recvonly\n" +
			"m=audio 3456 RTP/AVP 0\n" +
			"m=video 2232 RTP/AVP 31\n"),
		&Request{
			Method: "ANNOUNCE",
			Url:    &url.URL{Scheme: "rtsp", Host: "example.com", Path: "/media.mp4"},
			Header: Header{
				"CSeq":           []string{"7"},
				"Date":           []string{"23 Jan 1997 15:35:06 GMT"},
				"Session":        []string{"12345678"},
				"Content-Type":   []string{"application/sdp"},
				"Content-Length": []string{"306"},
			},
			Content: []byte("v=0\n" +
				"o=mhandley 2890844526 2890845468 IN IP4 126.16.64.4\n" +
				"s=SDP Seminar\n" +
				"i=A Seminar on the session description protocol\n" +
				"u=http://www.cs.ucl.ac.uk/staff/M.Handley/sdp.03.ps\n" +
				"e=mjh@isi.edu (Mark Handley)\n" +
				"c=IN IP4 224.2.17.12/127\n" +
				"t=2873397496 2873404696\n" +
				"a=recvonly\n" +
				"m=audio 3456 RTP/AVP 0\n" +
				"m=video 2232 RTP/AVP 31\n",
			),
		},
	},
	{
		"get_parameter",
		[]byte("GET_PARAMETER rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
			"CSeq: 9\r\n" +
			"Content-Length: 24\r\n" +
			"Content-Type: text/parameters\r\n" +
			"Session: 12345678\r\n" +
			"\r\n" +
			"packets_received\n" +
			"jitter\n"),
		&Request{
			Method: "GET_PARAMETER",
			Url:    &url.URL{Scheme: "rtsp", Host: "example.com", Path: "/media.mp4"},
			Header: Header{
				"CSeq":           []string{"9"},
				"Content-Type":   []string{"text/parameters"},
				"Session":        []string{"12345678"},
				"Content-Length": []string{"24"},
			},
			Content: []byte("packets_received\n" +
				"jitter\n",
			),
		},
	},
}

func TestRequestRead(t *testing.T) {
	for _, c := range casesRequest {
		t.Run(c.name, func(t *testing.T) {
			req, err := readRequestFromBytes(c.byts)
			require.NoError(t, err)
			require.Equal(t, c.req, req)
		})
	}
}

func TestRequestReadBadEnd(t *testing.T) {
	msg := []byte("GET_PARAMETER rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n")
	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "Could not find end of reqest to parse", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestReadBadH1(t *testing.T) {
	msg := []byte("GET_PARAMETER rtsp://example.com/media.mp4\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n\r\n")
	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "unable to parse Request Header", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestReadBadURL(t *testing.T) {
	msg := []byte("GET_PARAMETER   RTSP/1.0\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n\r\n")
	msg[14] = 0x08

	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "unable to parse url '\b'", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestReadBadScheme(t *testing.T) {
	msg := []byte("GET_PARAMETER https://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n\r\n")
	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "invalid url scheme 'https'", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestReadBadProto(t *testing.T) {
	msg := []byte("GET_PARAMETER rtsp://example.com/media.mp4 RTSP/1.2\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n\r\n")
	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "expected 'RTSP/1.0', got 'RTSP/1.2'", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestReadContentLength(t *testing.T) {
	msg := []byte("GET_PARAMETER rtsp://example.com/media.mp4 RTSP/1.0\r\n" +
		"CSeq: 9\r\n" +
		"Content-Length: 24\r\n" +
		"Content-Type: text/parameters\r\n" +
		"Session: 12345678\r\n\r\n")
	r, err := readRequestFromBytes(msg)
	require.Error(t, err)
	require.Equal(t, "Not enough bytes to get all content", err.Error())
	var r2 *Request = nil
	require.Equal(t, r2, r)
}

func TestRequestWrite(t *testing.T) {
	for _, c := range casesRequest {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.byts, []byte(c.req.String()))
		})
	}
}
