// Package slp implements a minimal Minecraft Java Edition Server List Ping
// client: just enough of the handshake + status protocol to confirm a server
// is alive and returning well-formed status JSON.
package slp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

const maxResponseLength = 1 << 20 // guards against a misbehaving peer claiming an absurd length

// Status performs a Server List Ping against host:port and returns an error
// unless a well-formed JSON status response is received within timeout.
func Status(host string, port int, timeout time.Duration) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}

	if err := writeHandshake(conn, host, port); err != nil {
		return fmt.Errorf("write handshake: %w", err)
	}
	if err := writePacket(conn, varInt(0x00)); err != nil { // empty status-request packet
		return fmt.Errorf("write status request: %w", err)
	}

	body, err := readPacket(bufio.NewReader(conn))
	if err != nil {
		return fmt.Errorf("read status response: %w", err)
	}

	br := bytes.NewReader(body)
	packetID, err := readVarInt(br)
	if err != nil {
		return fmt.Errorf("read response packet id: %w", err)
	}
	if packetID != 0x00 {
		return fmt.Errorf("unexpected response packet id: %d", packetID)
	}

	jsonLen, err := readVarInt(br)
	if err != nil {
		return fmt.Errorf("read json length: %w", err)
	}
	jsonBytes := make([]byte, jsonLen)
	if _, err := io.ReadFull(br, jsonBytes); err != nil {
		return fmt.Errorf("read json body: %w", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(jsonBytes, &payload); err != nil {
		return fmt.Errorf("invalid status json: %w", err)
	}

	return nil
}

func writeHandshake(w io.Writer, host string, port int) error {
	var packet bytes.Buffer
	packet.Write(varInt(0x00)) // packet ID
	packet.Write(varInt(765))  // protocol version; unvalidated by the server for a status ping
	packet.Write(packString(host))

	portBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(portBytes, uint16(port))
	packet.Write(portBytes)

	packet.Write(varInt(1)) // next state: 1 = status

	return writePacket(w, packet.Bytes())
}

func writePacket(w io.Writer, data []byte) error {
	if _, err := w.Write(varInt(int32(len(data)))); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func readPacket(r *bufio.Reader) ([]byte, error) {
	length, err := readVarInt(r)
	if err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}
	if length < 0 || length > maxResponseLength {
		return nil, fmt.Errorf("invalid packet length: %d", length)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return data, nil
}

func packString(s string) []byte {
	b := []byte(s)
	return append(varInt(int32(len(b))), b...)
}

func varInt(n int32) []byte {
	buf := make([]byte, 0, 5)
	u := uint32(n)
	for {
		b := byte(u & 0x7F)
		u >>= 7
		if u != 0 {
			b |= 0x80
		}
		buf = append(buf, b)
		if u == 0 {
			break
		}
	}
	return buf
}

type byteReader interface {
	io.Reader
	io.ByteReader
}

func readVarInt(r byteReader) (int32, error) {
	var result int32
	var shift uint
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		result |= int32(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 35 {
			return 0, fmt.Errorf("varint too long")
		}
	}
	return result, nil
}
