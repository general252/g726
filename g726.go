package g726

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type G726Rate int

const (
	G726Rate16kbps G726Rate = 0
	G726Rate24kbps G726Rate = 1
	G726Rate32kbps G726Rate = 2
	G726Rate40kbps G726Rate = 3
)

func (r G726Rate) String() string {
	switch r {
	case G726Rate16kbps:
		return "16kbps"
	case G726Rate24kbps:
		return "24kbps"
	case G726Rate32kbps:
		return "32kbps"
	case G726Rate40kbps:
		return "40kbps"
	default:
		return ""
	}
}

func (state_ptr *G726_state) DecodeFrameSize(bytesSize int) (int, error) {
	switch state_ptr.rate {
	case G726Rate16kbps:
		return bytesSize * 8, nil
	case G726Rate24kbps:
		if bytesSize%3 != 0 {
			return -1, fmt.Errorf("input length must be multiple of 3 for 24kbps decoding")
		}
		return (bytesSize * 8 / 3) * 2, nil
	case G726Rate32kbps:
		return bytesSize * 2 * 2, nil
	case G726Rate40kbps:
		if bytesSize%5 != 0 {
			return -1, fmt.Errorf("input length must be multiple of 5 for 40kbps decoding")
		}
		return (bytesSize * 8 / 5) * 2, nil
	default:
		return -1, fmt.Errorf("invalid rate")
	}
}

func (state_ptr *G726_state) EncodeFrameSize(bytesSize int) (int, error) {
	switch state_ptr.rate {
	case G726Rate16kbps:
		if bytesSize%8 != 0 {
			return -1, fmt.Errorf("input length must be a multiple of 4 for 16kbps encoding")
		}

		return bytesSize / 8, nil
	case G726Rate24kbps:
		if bytesSize%16 != 0 {
			return -1, fmt.Errorf("input length must be a multiple of 8 for 24kbps encoding")
		}
		return bytesSize / 16 * 3, nil
	case G726Rate32kbps:
		if bytesSize%4 != 0 {
			return -1, fmt.Errorf("input length must be a multiple of 2 for 32kbps encoding")
		}

		return bytesSize / 4, nil
	case G726Rate40kbps:
		if bytesSize%16 != 0 {
			return -1, fmt.Errorf("input length must be a multiple of 8 for 40kbps encoding")
		}
		return bytesSize / 10 * 8, nil
	}

	return -1, fmt.Errorf("invalid rate")
}

func (state_ptr *G726_state) EncodeToBytes(pcm []byte, pkt []byte) (int, error) {
	length := len(pcm)
	size, err := state_ptr.EncodeFrameSize(length)
	if err != nil {
		return -1, err
	}

	_ = pkt[size-1]
	var n int
	switch state_ptr.rate {
	case G726Rate16kbps:
		state := G726_state{}
		for i := 0; i < length; i += 8 {
			a := state.g726_16_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i:]))))
			b := state.g726_16_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+2:]))))
			c := state.g726_16_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+4:]))))
			d := state.g726_16_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+6:]))))

			// 4b -> 1b
			v := byte((a << 6) | (b << 4) | (c << 2) | d)
			pkt[n] = v
			n++
		}
		break
	case G726Rate24kbps:
		for i := 0; i < length; i += 16 {
			s0 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i:]))))
			s1 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+2:]))))
			s2 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+4:]))))
			s3 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+6:]))))
			s4 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+8:]))))
			s5 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+10:]))))
			s6 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+12:]))))
			s7 := state_ptr.g726_24_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+14:]))))

			// 打包8个3位值到3个字节
			// 8b -> 3b
			b0 := byte(s0<<5) | byte(s1<<2) | byte(s2>>1)
			b1 := byte((s2&1)<<7) | byte(s3<<4) | byte(s4<<1) | byte(s5>>2)
			b2 := byte((s5&3)<<6) | byte(s6<<3) | byte(s7)

			pkt[n] = b0
			pkt[n+1] = b1
			pkt[n+2] = b2
			n += 3
		}
		break
	case G726Rate32kbps:
		for i := 0; i < length; i += 4 {
			a := state_ptr.g726_32_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i:]))))
			b := state_ptr.g726_32_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+2:]))))
			// 2b -> 1b
			pkt[n] = byte((a << 4) | b)
			n++
		}
		break
	case G726Rate40kbps:
		for i := 0; i < length; i += 16 {
			s0 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i:]))))
			s1 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+2:]))))
			s2 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+4:]))))
			s3 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+6:]))))
			s4 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+8:]))))
			s5 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+10:]))))
			s6 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+12:]))))
			s7 := state_ptr.g726_40_encoder(int(int16(binary.LittleEndian.Uint16(pcm[i+14:]))))
			// 将8个5位值打包成5个字节
			b0 := byte((s0 << 3) | (s1 >> 2))
			b1 := byte(((s1 & 0x03) << 6) | (s2 << 1) | (s3 >> 4))
			b2 := byte(((s3 & 0x0F) << 4) | (s4 >> 1))
			b3 := byte(((s4 & 0x01) << 7) | (s5 << 2) | (s6 >> 3))
			b4 := byte(((s6 & 0x07) << 5) | s7)

			pkt[n] = b0
			pkt[n+1] = b1
			pkt[n+2] = b2
			pkt[n+3] = b3
			pkt[n+4] = b4
			n += 5
		}
		break
	}

	return n, nil
}

func (state_ptr *G726_state) Encode(pcm []int16) ([]byte, error) {
	size, err := state_ptr.EncodeFrameSize(len(pcm) * 2)
	if err != nil {
		return nil, err
	}

	pcmBytes := make([]byte, len(pcm)*2)
	for i, v := range pcm {
		binary.LittleEndian.PutUint16(pcmBytes[i*2:], uint16(v))
	}

	pkt := make([]byte, size)
	n, err := state_ptr.EncodeToBytes(pcmBytes, pkt)
	if err != nil {
		return nil, err
	}

	return pkt[:n], nil
}

func (state_ptr *G726_state) DecodeToBytes(bitstream []byte, pcm []byte) (int, error) {
	input_len := len(bitstream)
	pcmSize, err := state_ptr.DecodeFrameSize(input_len)
	if err != nil {
		return -1, err
	}

	_ = pcm[pcmSize-1]

	switch state_ptr.rate {
	case G726Rate16kbps:
		for i := 0; i < input_len; i++ {
			a := (bitstream[i] & byte(192)) >> 6
			b := (bitstream[i] & byte(48)) >> 4
			c := (bitstream[i] & byte(12)) >> 2
			d := (bitstream[i] & byte(3)) >> 0

			binary.LittleEndian.PutUint16(pcm[i*8:], uint16(state_ptr.g726_16_decoder(int(a))))
			binary.LittleEndian.PutUint16(pcm[i*8+2:], uint16(state_ptr.g726_16_decoder(int(b))))
			binary.LittleEndian.PutUint16(pcm[i*8+4:], uint16(state_ptr.g726_16_decoder(int(c))))
			binary.LittleEndian.PutUint16(pcm[i*8+6:], uint16(state_ptr.g726_16_decoder(int(d))))
		}
		return pcmSize, nil
	case G726Rate24kbps:
		for i := 0; i < input_len; i += 3 {
			b0 := bitstream[i]
			b1 := bitstream[i+1]
			b2 := bitstream[i+2]

			s0 := (b0 & 0xE0) >> 5
			s1 := (b0 & 0x1C) >> 2
			s2 := ((b0 & 0x03) << 1) | ((b1 & 0x80) >> 7)
			s3 := (b1 & 0x70) >> 4
			s4 := (b1 & 0x0E) >> 1
			s5 := ((b1 & 0x01) << 2) | ((b2 & 0xC0) >> 6)
			s6 := (b2 & 0x38) >> 3
			s7 := (b2 & 0x07) >> 0

			n := i / 3 * 16
			binary.LittleEndian.PutUint16(pcm[n:], uint16(state_ptr.g726_24_decoder(int(s0))))
			binary.LittleEndian.PutUint16(pcm[n+2:], uint16(state_ptr.g726_24_decoder(int(s1))))
			binary.LittleEndian.PutUint16(pcm[n+4:], uint16(state_ptr.g726_24_decoder(int(s2))))
			binary.LittleEndian.PutUint16(pcm[n+6:], uint16(state_ptr.g726_24_decoder(int(s3))))
			binary.LittleEndian.PutUint16(pcm[n+8:], uint16(state_ptr.g726_24_decoder(int(s4))))
			binary.LittleEndian.PutUint16(pcm[n+10:], uint16(state_ptr.g726_24_decoder(int(s5))))
			binary.LittleEndian.PutUint16(pcm[n+12:], uint16(state_ptr.g726_24_decoder(int(s6))))
			binary.LittleEndian.PutUint16(pcm[n+14:], uint16(state_ptr.g726_24_decoder(int(s7))))
		}
		return pcmSize, nil
	case G726Rate32kbps:
		for i := 0; i < input_len; i++ {
			a := (bitstream[i] & byte(240)) >> 4
			b := (bitstream[i] & byte(15)) >> 0
			binary.LittleEndian.PutUint16(pcm[i*4:], uint16(state_ptr.g726_32_decoder(int(a))))
			binary.LittleEndian.PutUint16(pcm[i*4+2:], uint16(state_ptr.g726_32_decoder(int(b))))
		}

		return pcmSize, nil
	case G726Rate40kbps:
		for i := 0; i < input_len; i += 5 {
			b0 := bitstream[i]
			b1 := bitstream[i+1]
			b2 := bitstream[i+2]
			b3 := bitstream[i+3]
			b4 := bitstream[i+4]

			// 解包5个字节到8个5位值
			s0 := (b0 & 0xF8) >> 3
			s1 := ((b0 & 0x07) << 2) | ((b1 & 0xC0) >> 6)
			s2 := (b1 & 0x3E) >> 1
			s3 := ((b1 & 0x01) << 4) | ((b2 & 0xF0) >> 4)
			s4 := ((b2 & 0x0F) << 1) | ((b3 & 0x80) >> 7)
			s5 := (b3 & 0x7C) >> 2
			s6 := ((b3 & 0x03) << 3) | ((b4 & 0xE0) >> 5)
			s7 := (b4 & 0x1F) >> 0

			n := i / 5 * 16
			binary.LittleEndian.PutUint16(pcm[n:], uint16(state_ptr.g726_40_decoder(int(s0))))
			binary.LittleEndian.PutUint16(pcm[n+2:], uint16(state_ptr.g726_40_decoder(int(s1))))
			binary.LittleEndian.PutUint16(pcm[n+4:], uint16(state_ptr.g726_40_decoder(int(s2))))
			binary.LittleEndian.PutUint16(pcm[n+6:], uint16(state_ptr.g726_40_decoder(int(s3))))
			binary.LittleEndian.PutUint16(pcm[n+8:], uint16(state_ptr.g726_40_decoder(int(s4))))
			binary.LittleEndian.PutUint16(pcm[n+10:], uint16(state_ptr.g726_40_decoder(int(s5))))
			binary.LittleEndian.PutUint16(pcm[n+12:], uint16(state_ptr.g726_40_decoder(int(s6))))
			binary.LittleEndian.PutUint16(pcm[n+14:], uint16(state_ptr.g726_40_decoder(int(s7))))
		}
		return pcmSize, nil
	default:
		return -1, fmt.Errorf("invalid rate")
	}
}

func (state_ptr *G726_state) Decode(bitstream []byte) ([]int16, error) {
	size, err := state_ptr.DecodeFrameSize(len(bitstream))
	if err != nil {
		return nil, err
	} else if size%2 != 0 {
		panic("input length must be even")
	}

	pcmData := make([]byte, size)
	_, err = state_ptr.DecodeToBytes(bitstream, pcmData)
	if err != nil {
		return nil, err
	}

	return (*[1 << 30]int16)(unsafe.Pointer(&pcmData[0]))[: size/2 : size/2], nil
}

func (state_ptr *G726_state) EncodeSimple(pcm []byte) ([]byte, error) {
	if len(pcm)%2 != 0 {
		return nil, fmt.Errorf("pcm length must be even")
	}

	return state_ptr.Encode((*[1 << 30]int16)(unsafe.Pointer(&pcm[0]))[: len(pcm)/2 : len(pcm)/2])
}

func (state_ptr *G726_state) DecodeSimple(bitstream []byte) ([]byte, error) {
	pcm_out, err := state_ptr.Decode(bitstream)
	if err != nil {
		return nil, err
	}

	pcm := make([]byte, len(pcm_out)*2)
	for i := 0; i < len(pcm_out); i++ {
		binary.LittleEndian.PutUint16(pcm[2*i:2*i+2], uint16(pcm_out[i]))
	}
	return pcm, nil
}
