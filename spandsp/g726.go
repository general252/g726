package spandsp

import (
	"errors"
	"fmt"
)

/*
 * Maps G.726_16 code word to reconstructed scale factor normalized log
 * magnitude values.
 */
var g726_16_dqlntab = [4]int_t{
	116, 365, 365, 116,
}

/* Maps G.726_16 code word to log of scale factor multiplier. */
var g726_16_witab = [4]int_t{
	-704, 14048, 14048, -704,
}

/*
 * Maps G.726_16 code words to a set of values whose long and short
 * term averages are computed and then compared to give an indication
 * how stationary (steady state) the signal is.
 */
var g726_16_fitab = [4]int_t{
	0x000, 0xE00, 0xE00, 0x000,
}

var qtab_726_16 = [1]int_t{
	261,
}

/*
 * Maps G.726_24 code word to reconstructed scale factor normalized log
 * magnitude values.
 */
var g726_24_dqlntab = [8]int_t{
	-2048, 135, 273, 373, 373, 273, 135, -2048,
}

/* Maps G.726_24 code word to log of scale factor multiplier. */
var g726_24_witab = [8]int_t{-128, 960, 4384, 18624, 18624, 4384, 960, -128}

/*
 * Maps G.726_24 code words to a set of values whose long and short
 * term averages are computed and then compared to give an indication
 * how stationary (steady state) the signal is.
 */
var g726_24_fitab = [8]int_t{
	0x000, 0x200, 0x400, 0xE00, 0xE00, 0x400, 0x200, 0x000,
}

var qtab_726_24 = [3]int_t{
	8, 218, 331,
}

/*
 * Maps G.726_32 code word to reconstructed scale factor normalized log
 * magnitude values.
 */
var g726_32_dqlntab = [16]int_t{
	-2048, 4, 135, 213, 273, 323, 373, 425,
	425, 373, 323, 273, 213, 135, 4, -2048,
}

/* Maps G.726_32 code word to log of scale factor multiplier. */
var g726_32_witab = [16]int_t{
	-384, 576, 1312, 2048, 3584, 6336, 11360, 35904,
	35904, 11360, 6336, 3584, 2048, 1312, 576, -384,
}

/*
 * Maps G.726_32 code words to a set of values whose long and short
 * term averages are computed and then compared to give an indication
 * how stationary (steady state) the signal is.
 */
var g726_32_fitab = [16]int_t{
	0x000, 0x000, 0x000, 0x200, 0x200, 0x200, 0x600, 0xE00,
	0xE00, 0x600, 0x200, 0x200, 0x200, 0x000, 0x000, 0x000,
}

var qtab_726_32 = [7]int_t{
	-124, 80, 178, 246, 300, 349, 400,
}

/*
 * Maps G.726_40 code word to ructeconstructed scale factor normalized log
 * magnitude values.
 */
var g726_40_dqlntab = [32]int_t{
	-2048, -66, 28, 104, 169, 224, 274, 318,
	358, 395, 429, 459, 488, 514, 539, 566,
	566, 539, 514, 488, 459, 429, 395, 358,
	318, 274, 224, 169, 104, 28, -66, -2048,
}

/* Maps G.726_40 code word to log of scale factor multiplier. */
var g726_40_witab = [32]int_t{
	448, 448, 768, 1248, 1280, 1312, 1856, 3200,
	4512, 5728, 7008, 8960, 11456, 14080, 16928, 22272,
	22272, 16928, 14080, 11456, 8960, 7008, 5728, 4512,
	3200, 1856, 1312, 1280, 1248, 768, 448, 448,
}

/*
 * Maps G.726_40 code words to a set of values whose long and short
 * term averages are computed and then compared to give an indication
 * how stationary (steady state) the signal is.
 */
var g726_40_fitab = [32]int_t{
	0x000, 0x000, 0x000, 0x000, 0x000, 0x200, 0x200, 0x200,
	0x200, 0x200, 0x400, 0x600, 0x800, 0xA00, 0xC00, 0xC00,
	0xC00, 0xC00, 0xA00, 0x800, 0x600, 0x400, 0x200, 0x200,
	0x200, 0x200, 0x200, 0x000, 0x000, 0x000, 0x000, 0x000,
}

var qtab_726_40 = [15]int_t{
	-122, -16, 68, 139, 198, 250, 298, 339,
	378, 413, 445, 475, 502, 528, 553,
}

/*
 * returns the integer product of the 14-bit integer "an" and
 * "floating point" representation (4-bit exponent, 6-bit mantessa) "srn".
 */
func fmult(an int_t, srn int_t) int_t {
	var (
		anmag   int_t
		anexp   int_t
		anmant  int_t
		wanexp  int_t
		wanmant int_t
		retval  int_t
	)

	// anmag = (an > 0)  ?  an  :  ((-an) & 0x1FFF);
	if an > 0 {
		anmag = an
	} else {
		anmag = -an & 0x1FFF
	}

	anexp = top_bit(uint32_t(anmag)) - 5
	// anmant = (anmag == 0)  ?  32  :  (anexp >= 0)  ?  (anmag >> anexp)  :  (anmag << -anexp)
	if anmag == 0 {
		anmant = 32
	} else if anexp >= 0 {
		anmant = anmag >> anexp
	} else {
		anmant = anmag << -anexp
	}
	wanexp = anexp + ((srn >> 6) & 0xF) - 13

	wanmant = (anmant*(srn&0x3F) + 0x30) >> 4
	// retval = (wanexp >= 0)  ?  ((wanmant << wanexp) & 0x7FFF)  :  (wanmant >> -wanexp)
	if wanexp >= 0 {
		retval = (wanmant << wanexp) & 0x7FFF
	} else {
		retval = wanmant >> -wanexp
	}

	//return ((an ^ srn) < 0)  ?  -retval  :  retval)
	if (an ^ srn) < 0 {
		return -retval
	} else {
		return retval
	}
}

/*
 * Compute the estimated signal from the 6-zero predictor.
 */
func (s *g726_state_t) predictor_zero() int_t {
	var i int_t
	var sezi int_t

	sezi = fmult(s.b[0]>>2, s.dq[0])
	/* ACCUM */
	for i = 1; i < 6; i++ {
		sezi += fmult(s.b[i]>>2, s.dq[i])
	}

	return sezi
}

/*
 * Computes the estimated signal from the 2-pole predictor.
 */
func (s *g726_state_t) predictor_pole() int_t {
	return fmult(s.a[1]>>2, s.sr[1]) + fmult(s.a[0]>>2, s.sr[0])
}

/*
 * Computes the quantization step size of the adaptive quantizer.
 */
func (s *g726_state_t) step_size() int_t {
	var y int_t
	var dif int_t
	var al int_t

	if s.ap >= 256 {
		return s.yu
	}

	y = s.yl >> 6
	dif = s.yu - y
	al = s.ap >> 2
	if dif > 0 {
		y += (dif * al) >> 6
	} else if dif < 0 {
		y += (dif*al + 0x3F) >> 6
	}

	return y
}

func quantize(d int_t, y int_t, table []int_t, quantizer_states int_t) int_t {
	var dqm int_t  /* Magnitude of 'd' */
	var exp int_t  /* Integer part of base 2 log of 'd' */
	var mant int_t /* Fractional part of base 2 log */
	var dl int_t   /* Log of magnitude of 'd' */
	var dln int_t  /* Step size scale factor normalized log */
	var i int_t
	var size int_t

	/*
	 * LOG
	 *
	 * Compute base 2 log of 'd', and store in 'dl'.
	 */
	// dqm = (int16_t) abs(d);
	dqm = d
	if d < 0 {
		dqm = -d
	}

	exp = top_bit(uint32_t(dqm>>1)) + 1
	/* Fractional portion. */
	mant = ((dqm << 7) >> exp) & 0x7F
	dl = (exp << 7) + mant

	/*
	 * SUBTB
	 *
	 * "Divide" by step size multiplier.
	 */
	dln = dl - y>>2

	/*
	 * QUAN
	 *
	 * Search for codword i for 'dln'.
	 */
	size = (quantizer_states - 1) >> 1
	for i = 0; i < size; i++ {
		if dln < table[i] {
			break
		}
	}

	if d < 0 {
		/* Take 1's complement of i */
		return (size << 1) + 1 - i
	}

	if i == 0 && (quantizer_states&1) != 0 {
		/* Zero is only valid if there are an even number of states, so
		   take the 1's complement if the code is zero. */
		return quantizer_states
	}
	return i
}

/*
 * Returns reconstructed difference signal 'dq' obtained from
 * codeword 'i' and quantization step size scale factor 'y'.
 * Multiplication is performed in log base 2 domain as addition.
 */
func reconstruct(
	sign int_t, /* 0 for non-negative value */
	dqln int_t, /* G.72x codeword */
	y int_t,    /* Step size multiplier */
) int_t {
	var dql int_t /* Log of 'dq' magnitude */
	var dex int_t /* Integer part of log */
	var dqt int_t
	var dq int_t /* Reconstructed difference signal sample */

	dql = dqln + (y >> 2) /* ADDA */

	if dql < 0 {
		if sign != 0 {
			return -0x8000
		} else {
			return 0
		}
	}

	/* ANTILOG */
	dex = (dql >> 7) & 15
	dqt = 128 + (dql & 127)
	dq = (dqt << 7) >> (14 - dex)

	if sign != 0 {
		var tmp = 0x8000
		return dq - int_t(tmp)
	} else {
		return dq
	}
}

func (s *g726_state_t) update(y, wi, fi, dq, sr, dqsez int_t) {
	var mag int_t
	var exp int_t
	var a2p int_t  /* LIMC */
	var a1ul int_t /* UPA1 */
	var pks1 int_t /* UPA2 */
	var fa1 int_t
	var ylint int_t
	var dqthr int_t
	var ylfrac int_t
	var thr int_t
	var pk0 int_t
	var i int_t
	var tr bool

	_ = thr

	a2p = 0
	/* Needed in updating predictor poles */
	// pk0 = (dqsez < 0)  ?  1  :  0;
	if dqsez < 0 {
		pk0 = 1
	} else {
		pk0 = 0
	}

	/* prediction difference magnitude */
	mag = (int_t)(dq & 0x7FFF)
	/* TRANS */
	ylint = int_t(s.yl >> 15)           /* exponent part of yl */
	ylfrac = int_t((s.yl >> 10) & 0x1F) /* fractional part of yl */
	/* Limit threshold to 31 << 10 */
	// thr = (ylint > 9)  ?  (31 << 10)  :  ((32 + ylfrac) << ylint);
	if ylint > 9 {
		thr = 31 << 10
	} else {
		thr = (32 + ylfrac) << ylint
	}
	dqthr = (thr + (thr >> 1)) >> 1

	if !s.td { /* signal supposed voice */
		tr = false
	} else if mag <= dqthr { /* supposed data, but small mag */
		tr = false /* treated as voice */
	} else { /* signal is data (modem) */
		tr = true
	}

	/*
	 * Quantizer scale factor adaptation.
	 */

	/* FUNCTW & FILTD & DELAY */
	/* update non-steady state step size multiplier */
	s.yu = overflow(int16_t(y + ((wi - y) >> 5)))

	/* LIMB */
	if s.yu < 544 {
		s.yu = 544
	} else if s.yu > 5120 {
		s.yu = 5120
	}

	/* FILTE & DELAY */
	/* update steady state step size multiplier */
	s.yl += s.yu + (-s.yl)>>6

	/*
	 * Adaptive predictor coefficients.
	 */
	if tr {
		/* Reset the a's and b's for a modem signal */
		s.a[0] = 0
		s.a[1] = 0
		s.b[0] = 0
		s.b[1] = 0
		s.b[2] = 0
		s.b[3] = 0
		s.b[4] = 0
		s.b[5] = 0
	} else {
		/* Update the a's and b's */
		/* UPA2 */
		pks1 = pk0 ^ s.pk[0]
		/* Update predictor pole a[1] */
		a2p = s.a[1] - s.a[1]>>7
		if dqsez != 0 {
			// fa1 = (pks1)  ?  s->a[0]  :  -s->a[0];
			if pks1 != 0 {
				fa1 = s.a[0]
			} else {
				fa1 = -s.a[0]
			}

			/* a2p = function of fa1 */
			if fa1 < -8191 {
				a2p -= 0x100
			} else if fa1 > 8191 {
				a2p += 0xFF
			} else {
				a2p += fa1 >> 5
			}

			if (pk0 ^ s.pk[1]) != 0 {
				/* LIMC */
				if a2p <= -12160 {
					a2p = -12288
				} else if a2p >= 12416 {
					a2p = 12288
				} else {
					a2p -= 0x80
				}
			} else {
				if a2p <= -12416 {
					a2p = -12288
				} else if a2p >= 12160 {
					a2p = 12288
				} else {
					a2p += 0x80
				}
			}
		}

		/* TRIGB & DELAY */
		s.a[1] = a2p

		/* UPA1 */
		/* Update predictor pole a[0] */
		s.a[0] -= s.a[0] >> 8
		if dqsez != 0 {
			if pks1 == 0 {
				s.a[0] += 192
			} else {
				s.a[0] -= 192
			}
		}

		/* LIMD */
		a1ul = 15360 - a2p
		if s.a[0] < -a1ul {
			s.a[0] = -a1ul
		} else if s.a[0] > a1ul {
			s.a[0] = a1ul
		}

		/* UPB : update predictor zeros b[6] */
		for i = 0; i < 6; i++ {
			/* Distinguish 40Kbps mode from the others */
			// s->b[i] -= s->b[i] >> ((s->bits_per_sample == 5)  ?  9  :  8);
			if s.bits_per_sample == 5 {
				s.b[i] -= s.b[i] >> 9
			} else {
				s.b[i] -= s.b[i] >> 8
			}

			if (dq & 0x7FFF) != 0 {
				/* XOR */
				if (dq ^ int_t(s.dq[i])) >= 0 {
					s.b[i] += 128
				} else {
					s.b[i] -= 128
				}
			}
		}
	}

	for i = 5; i > 0; i-- {
		s.dq[i] = s.dq[i-1]
	}

	/* FLOAT A : convert dq[0] to 4-bit exp, 6-bit mantissa f.p. */
	if mag == 0 {
		//s->dq[0] = (dq >= 0)  ?  0x20  :  0xFC20;
		if dq >= 0 {
			s.dq[0] = 0x20
		} else {
			var tmp = 0xFC20
			s.dq[0] = overflow(int16_t(tmp))
		}
	} else {
		exp = top_bit(uint32_t(mag)) + 1
		//s->dq[0] = (dq >= 0)
		//?  ((exp << 6) + ((mag << 6) >> exp))
		//:  ((exp << 6) + ((mag << 6) >> exp) - 0x400);
		if dq >= 0 {
			s.dq[0] = overflow(int16_t((exp << 6) + ((mag << 6) >> exp)))
		} else {
			s.dq[0] = overflow(int16_t((exp << 6) + ((mag << 6) >> exp) - 0x400))
		}
	}

	s.sr[1] = s.sr[0]
	/* FLOAT B : convert sr to 4-bit exp., 6-bit mantissa f.p. */
	if sr == 0 {
		s.sr[0] = 0x20
	} else if sr > 0 {
		exp = top_bit(uint32_t(sr)) + 1
		s.sr[0] = (exp << 6) + ((sr << 6) >> exp)
	} else if sr > -32768 {
		mag = -sr
		exp = top_bit(uint32_t(mag)) + 1

		s.sr[0] = (exp << 6) + ((mag << 6) >> exp) - 0x400
	} else {
		var tmp = 0xFC20
		s.sr[0] = int_t(tmp)
	}

	/* DELAY A */
	s.pk[1] = s.pk[0]
	s.pk[0] = pk0

	/* TONE */
	if tr { /* this sample has been treated as data */
		s.td = false /* next one will be treated as voice */
	} else if a2p < -11776 { /* small sample-to-sample correlation */
		s.td = true /* signal may be data */
	} else { /* signal is voice */
		s.td = false
	}

	/* Adaptation speed control. */
	/* FILTA */
	s.dms += (fi - s.dms) >> 5
	/* FILTB */
	s.dml += (fi<<2 - s.dml) >> 7

	tmp := (s.dms << 2) - s.dml
	if tmp < 0 {
		tmp = -tmp
	}

	if tr {
		s.ap = 256
	} else if y < 1536 { /* SUBTC */
		s.ap += (0x200 - s.ap) >> 4
	} else if s.td {
		s.ap += (0x200 - s.ap) >> 4
	} else if tmp >= (s.dml >> 3) {
		s.ap += (0x200 - s.ap) >> 4
	} else {
		s.ap += (-s.ap) >> 4
	}

	s.packets += 1
}

func tandem_adjust_alaw(
	sr int16_t, /* decoder output linear PCM sample */
	se int_t,   /* predictor estimate sample */
	y int_t,    /* quantizer step size */
	i int_t,    /* decoder input code */
	sign int_t,
	qtab []int_t,
	quantizer_states int_t,
) int16_t {
	var sp int_t /* A-law compressed 8-bit code */
	var dx int_t /* prediction error */
	var id int_t /* quantized prediction error */
	var sd int_t /* adjusted A-law decoded sample value */

	if sr <= -32768 {
		sr = -1
	}

	sp = int_t(linear_to_alaw((int32_t(sr) >> 1) << 3))
	/* 16-bit prediction error */
	dx = int_t(alaw_to_linear(uint8_t(sp)))>>2 - se
	id = quantize(dx, y, qtab, quantizer_states)
	if id == i {
		/* No adjustment of sp required */
		return int16_t(sp)
	}

	/* sp adjustment needed */
	/* ADPCM codes : 8, 9, ... F, 0, 1, ... , 6, 7 */
	/* 2's complement to biased unsigned */
	if (id ^ sign) > (i ^ sign) {
		/* sp adjusted to next lower value */
		if (sp & 0x80) != 0 {
			// sd = (sp == 0xD5)  ?  0x55  :  (((sp ^ 0x55) - 1) ^ 0x55);
			if sp == 0xD5 {
				sd = 0x55
			} else {
				sd = ((sp ^ 0x55) - 1) ^ 0x55
			}
		} else {
			// sd = (sp == 0x2A)  ?  0x2A  :  (((sp ^ 0x55) + 1) ^ 0x55);
			if sp == 0x2A {
				sd = 0x2A
			} else {
				sd = ((sp ^ 0x55) + 1) ^ 0x55
			}
		}
	} else {
		/* sp adjusted to next higher value */
		if (sp & 0x80) != 0 {
			// sd = (sp == 0xAA)  ?  0xAA  :  (((sp ^ 0x55) + 1) ^ 0x55);
			if sp == 0xAA {
				sd = 0xAA
			} else {
				sd = ((sp ^ 0x55) + 1) ^ 0x55
			}
		} else {
			// sd = (sp == 0x55)  ?  0xD5  :  (((sp ^ 0x55) - 1) ^ 0x55);
			if sp == 0x55 {
				sd = 0xD5
			} else {
				sd = ((sp ^ 0x55) - 1) ^ 0x55
			}
		}
	}

	return int16_t(sd)
}

func tandem_adjust_ulaw(
	sr int16_t, /* decoder output linear PCM sample */
	se int_t,   /* predictor estimate sample */
	y int_t,    /* quantizer step size */
	i int_t,    /* decoder input code */
	sign int_t,
	qtab []int_t,
	quantizer_states int_t,
) int16_t {
	var sp int_t /* u-law compressed 8-bit code */
	var dx int_t /* prediction error */
	var id int_t /* quantized prediction error */
	var sd int_t /* adjusted u-law decoded sample value */

	if sr <= -32768 {
		sr = 0
	}

	sp = int_t(linear_to_ulaw(int32_t(sr << 2)))
	/* 16-bit prediction error */
	dx = int_t(ulaw_to_linear(uint8_t(sp)))>>2 - se
	id = int_t(quantize(dx, y, qtab, quantizer_states))
	if id == i {
		/* No adjustment of sp required. */
		return (int16_t)(sp)
	}
	/*endif*/
	/* ADPCM codes : 8, 9, ... F, 0, 1, ... , 6, 7 */
	/* 2's complement to biased unsigned */
	if (id ^ sign) > (i ^ sign) {
		/* sp adjusted to next lower value */
		if (sp & 0x80) != 0 {
			// sd = (sp == 0xFF)  ?  0x7E  :  (sp + 1);
			if sp == 0xFF {
				sd = 0x7E
			} else {
				sd = sp + 1
			}
		} else {
			// sd = (sp == 0x00)  ?  0x00  :  (sp - 1);
			if sp == 0x00 {
				sd = 0x00
			} else {
				sd = sp - 1
			}
		}
	} else {
		/* sp adjusted to next higher value */
		if (sp & 0x80) != 0 {
			// sd = (sp == 0x80)  ?  0x80  :  (sp - 1);
			if sp == 0x80 {
				sd = 0x80
			} else {
				sd = sp - 1
			}
		} else {
			// sd = (sp == 0x7F)  ?  0xFE  :  (sp + 1);
			if sp == 0x7F {
				sd = 0xFE
			} else {
				sd = sp + 1
			}
		}
	}
	/*endif*/
	return (int16_t)(sd)
}

// Encodes a linear PCM, A-law or u-law input sample and returns its 3-bit code.
func (s *g726_state_t) g726_16_encoder(amp int16_t) uint8_t {
	var y int_t
	var sei int_t
	var sezi int_t
	var se int_t
	var d int_t
	var sr int_t
	var dqsez int_t
	var dq int_t
	var i int_t

	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()
	se = sei >> 1
	d = int_t(amp) - se

	/* Quantize prediction difference */
	y = s.step_size()
	i = quantize(d, y, qtab_726_16[:], 4)
	dq = reconstruct(i&2, g726_16_dqlntab[i], y)

	/* Reconstruct the signal */
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_16_witab[i], g726_16_fitab[i], dq, sr, dqsez)
	return (uint8_t)(i)
}

// Decodes a 2-bit CCITT G.726_16 ADPCM code and returns
// the resulting 16-bit linear PCM, A-law or u-law sample value.
func (s *g726_state_t) g726_16_decoder(code uint8_t) int16_t {
	var sezi int_t
	var sei int_t
	var se int_t
	var sr int_t
	var dq int_t
	var dqsez int_t
	var y int_t

	/* Mask to get proper bits */
	code &= 0x03
	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()

	y = s.step_size()
	dq = reconstruct(int_t(code&2), g726_16_dqlntab[code], y)

	/* Reconstruct the signal */
	se = sei >> 1
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_16_witab[code], g726_16_fitab[code], dq, sr, dqsez)

	switch s.ext_coding {
	case G726_ENCODING_ALAW:
		return tandem_adjust_alaw(int16_t(sr), se, y, int_t(code), 2, qtab_726_16[:], 4)
	case G726_ENCODING_ULAW:
		return tandem_adjust_ulaw(int16_t(sr), se, y, int_t(code), 2, qtab_726_16[:], 4)
	}

	return int16_t(sr << 2)
}

// g726_24_encoder
// Encodes a linear PCM, A-law or u-law input sample and returns its 3-bit code.
func (s *g726_state_t) g726_24_encoder(amp int16_t) uint8_t {
	var sei int_t
	var sezi int_t
	var se int_t
	var d int_t
	var sr int_t
	var dqsez int_t
	var dq int_t
	var i int_t
	var y int_t

	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()
	se = sei >> 1
	d = int_t(amp) - se

	/* Quantize prediction difference */
	y = s.step_size()
	i = quantize(d, y, qtab_726_24[:], 7)
	dq = reconstruct(i&4, g726_24_dqlntab[i], y)

	/* Reconstruct the signal */
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_24_witab[i], g726_24_fitab[i], dq, sr, dqsez)
	return (uint8_t)(i)
}

/*
 * Decodes a 3-bit CCITT G.726_24 ADPCM code and returns
 * the resulting 16-bit linear PCM, A-law or u-law sample value.
 */
func (s *g726_state_t) g726_24_decoder(code uint8_t) int16_t {
	var sezi int_t
	var sei int_t
	var se int_t
	var sr int_t
	var dq int_t
	var dqsez int_t
	var y int_t

	/* Mask to get proper bits */
	code &= 0x07
	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()

	y = s.step_size()
	dq = reconstruct(int_t(code&4), g726_24_dqlntab[code], y)

	/* Reconstruct the signal */
	se = sei >> 1
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_24_witab[code], g726_24_fitab[code], dq, sr, dqsez)

	switch s.ext_coding {
	case G726_ENCODING_ALAW:
		return tandem_adjust_alaw(int16_t(sr), se, y, int_t(code), 4, qtab_726_24[:], 7)
	case G726_ENCODING_ULAW:
		return tandem_adjust_ulaw(int16_t(sr), se, y, int_t(code), 4, qtab_726_24[:], 7)
	}

	return int16_t(sr << 2)
}

// g726_32_encoder
// Encodes a linear input sample and returns its 4-bit code.
func (s *g726_state_t) g726_32_encoder(amp int16_t) uint8_t {
	var sei int_t
	var sezi int_t
	var se int_t
	var d int_t
	var sr int_t
	var dqsez int_t
	var dq int_t
	var i int_t
	var y int_t

	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()
	se = sei >> 1
	d = int_t(amp) - se

	/* Quantize the prediction difference */
	y = s.step_size()
	i = quantize(d, y, qtab_726_32[:], 15)
	dq = reconstruct(i&8, g726_32_dqlntab[i], y)

	/* Reconstruct the signal */
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_32_witab[i], g726_32_fitab[i], dq, sr, dqsez)
	return (uint8_t)(i)
}

func (s *g726_state_t) g726_32_decoder(code uint8_t) int16_t {
	var sezi int_t
	var sei int_t
	var se int_t
	var sr int_t
	var dq int_t
	var dqsez int_t
	var y int_t

	/* Mask to get proper bits */
	code &= 0x0F
	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()

	y = s.step_size()
	dq = reconstruct(int_t(code&8), g726_32_dqlntab[code], y)

	/* Reconstruct the signal */
	se = sei >> 1
	// sr = (dq < 0)  ?  (se - (dq & 0x3FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x3FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_32_witab[code], g726_32_fitab[code], dq, sr, dqsez)

	switch s.ext_coding {
	case G726_ENCODING_ALAW:
		return tandem_adjust_alaw(int16_t(sr), se, y, int_t(code), 8, qtab_726_32[:], 15)
	case G726_ENCODING_ULAW:
		return tandem_adjust_ulaw(int16_t(sr), se, y, int_t(code), 8, qtab_726_32[:], 15)
	}

	return int16_t(sr << 2)
}

// g726_40_encoder
// Encodes a 16-bit linear PCM, A-law or u-law input sample and retuens
// the resulting 5-bit CCITT G.726 40Kbps code.
func (s *g726_state_t) g726_40_encoder(amp int16_t) uint8_t {
	var sei int_t
	var sezi int_t
	var se int_t
	var d int_t
	var sr int_t
	var dqsez int_t
	var dq int_t
	var i int_t
	var y int_t

	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()
	se = sei >> 1
	d = int_t(amp) - se

	/* Quantize prediction difference */
	y = s.step_size()
	i = quantize(d, y, qtab_726_40[:], 31)
	dq = reconstruct(i&0x10, g726_40_dqlntab[i], y)

	/* Reconstruct the signal */
	// sr = (dq < 0)  ?  (se - (dq & 0x7FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x7FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_40_witab[i], g726_40_fitab[i], dq, sr, dqsez)
	return (uint8_t)(i)
}

// g726_40_decoder
// Decodes a 5-bit CCITT G.726 40Kbps code and returns
// the resulting 16-bit linear PCM, A-law or u-law sample value.
func (s *g726_state_t) g726_40_decoder(code uint8_t) int16_t {
	var sezi int_t
	var sei int_t
	var se int_t
	var sr int_t
	var dq int_t
	var dqsez int_t
	var y int_t

	/* Mask to get proper bits */
	code &= 0x1F
	sezi = s.predictor_zero()
	sei = sezi + s.predictor_pole()

	y = s.step_size()
	dq = reconstruct(int_t(code&0x10), g726_40_dqlntab[code], y)

	/* Reconstruct the signal */
	se = sei >> 1
	// sr = (dq < 0)  ?  (se - (dq & 0x7FFF))  :  (se + dq);
	if dq < 0 {
		sr = se - (dq & 0x7FFF)
	} else {
		sr = se + dq
	}

	/* Pole prediction difference */
	dqsez = sr + (sezi >> 1) - se

	s.update(y, g726_40_witab[code], g726_40_fitab[code], dq, sr, dqsez)

	switch s.ext_coding {
	case G726_ENCODING_ALAW:
		return tandem_adjust_alaw(int16_t(sr), se, y, int_t(code), 0x10, qtab_726_40[:], 31)
	case G726_ENCODING_ULAW:
		return tandem_adjust_ulaw(int16_t(sr), se, y, int_t(code), 0x10, qtab_726_40[:], 31)
	}

	return int16_t(sr << 2)
}

func (s *g726_state_t) Decode(g726_data []uint8_t) (amp []int16_t) {
	var g726_bytes = len(g726_data)
	amp = make([]int16_t, 0, g726_bytes)
	var i int
	var samples int
	var code uint8_t
	var sl int16_t

	for {
		if s.packing != G726_PACKING_NONE {
			/* Unpack the code bits */
			if s.packing != G726_PACKING_LEFT {
				if s.bs.residue < s.bits_per_sample {
					if i >= g726_bytes {
						break
					}

					s.bs.bitstream |= uint32_t(g726_data[i]) << uint32_t(s.bs.residue)
					i += 1

					s.bs.residue += 8
				}
				code = (uint8_t)(s.bs.bitstream & ((1 << s.bits_per_sample) - 1))
				s.bs.bitstream >>= s.bits_per_sample
			} else {
				if s.bs.residue < s.bits_per_sample {
					if i >= g726_bytes {
						break
					}
					s.bs.bitstream = (s.bs.bitstream << 8) | uint32_t(g726_data[i])
					i += 1
					s.bs.residue += 8
				}

				code = (uint8_t)((s.bs.bitstream >> (s.bs.residue - s.bits_per_sample)) & ((1 << s.bits_per_sample) - 1))
			}
			s.bs.residue -= s.bits_per_sample
		} else {
			if i >= g726_bytes {
				break
			}
			code = g726_data[i]
			i += 1
		}

		sl = s.dec_func(code)
		if s.ext_coding != G726_ENCODING_LINEAR {
			// TODO:
			// amp[samples] = (uint8_t) sl;
			// samples+=1
		} else {
			// amp[samples] = sl
			amp = append(amp, sl)
			samples += 1
		}
	}

	return amp
}

func (s *g726_state_t) Encode(amp []int16_t) (g726_data []uint8_t) {
	g726_data = make([]uint8_t, 0, len(amp))

	var i int
	var g726_bytes int
	var sl int16_t
	var code uint8_t

	for ; i < len(amp); i++ {
		switch s.ext_coding {
		case G726_ENCODING_ALAW:
			// TODO:
			// sl = alaw_to_linear(((const uint8_t *) amp)[i]) >> 2;
		case G726_ENCODING_ULAW:
			// TODO:
			// sl = ulaw_to_linear(((const uint8_t *) amp)[i]) >> 2;
		default:
			sl = amp[i] >> 2
		}

		code = s.enc_func(sl)
		if false {
			lsb := 0
			if s.bs.lsb_first {
				lsb = 1
			}
			fmt.Printf("%d %4d %4d yl:%v yu:%v dms:%v dml:%v ap:%v a:%v %v b:%v %v %v %v %v %v pk:%v %v dq:%v %v %v %v %v %v sr:%v %v bitstream:%v residue:%v lsb:%v\n",
				s.packets,
				i, code, s.yl, s.yu, s.dms, s.dml, s.ap, s.a[0], s.a[1],
				s.b[0], s.b[1], s.b[2], s.b[3], s.b[4], s.b[5],
				s.pk[0], s.pk[1],
				s.dq[0], s.dq[1], s.dq[2], s.dq[3], s.dq[4], s.dq[5],
				s.sr[0], s.sr[1],
				s.bs.bitstream, s.bs.residue, lsb,
			)
		}

		if s.packing != G726_PACKING_NONE {
			/* Pack the code bits */
			if s.packing != G726_PACKING_LEFT {
				s.bs.bitstream |= uint32_t(code) << uint32_t(s.bs.residue)
				s.bs.residue += s.bits_per_sample
				if s.bs.residue >= 8 {
					// g726_data[g726_bytes] = (uint8_t)(s.bs.bitstream & 0xFF)
					g726_data = append(g726_data, (uint8_t)(s.bs.bitstream&0xFF))
					g726_bytes += 1
					s.bs.bitstream >>= 8
					s.bs.residue -= 8
				}
			} else {
				s.bs.bitstream = (s.bs.bitstream << uint32_t(s.bits_per_sample)) | uint32_t(code)
				s.bs.residue += s.bits_per_sample
				if s.bs.residue >= 8 {
					// g726_data[g726_bytes] = (uint8_t)((s.bs.bitstream >> (s.bs.residue - 8)) & 0xFF)
					g726_data = append(g726_data, (uint8_t)((s.bs.bitstream>>(s.bs.residue-8))&0xFF))
					g726_bytes += 1
					s.bs.residue -= 8
				}
			}
		} else {
			// g726_data[g726_bytes] = code
			g726_data = append(g726_data, code)
			g726_bytes += 1
		}
	}

	return g726_data
}

func G726_init(bit_rate, ext_coding, packing int32_t) (*g726_state_t, error) {
	if bit_rate != 16000 && bit_rate != 24000 && bit_rate != 32000 && bit_rate != 40000 {
		return nil, errors.New("invalid bit rate")
	}

	var i int

	s := &g726_state_t{}

	s.yl = 34816
	s.yu = 544
	s.dms = 0
	s.dml = 0
	s.ap = 0
	s.rate = bit_rate
	s.ext_coding = ext_coding
	s.packing = packing
	for i = 0; i < 2; i++ {
		s.a[i] = 0
		s.pk[i] = 0
		s.sr[i] = 32
	}
	for i = 0; i < 6; i++ {
		s.b[i] = 0
		s.dq[i] = 32
	}
	s.td = false
	switch bit_rate {
	case 16000:
		s.enc_func = s.g726_16_encoder
		s.dec_func = s.g726_16_decoder
		s.bits_per_sample = 2
	case 24000:
		s.enc_func = s.g726_24_encoder
		s.dec_func = s.g726_24_decoder
		s.bits_per_sample = 3
	case 40000:
		s.enc_func = s.g726_40_encoder
		s.dec_func = s.g726_40_decoder
		s.bits_per_sample = 5
		break
	default:
		// 32000
		s.enc_func = s.g726_32_encoder
		s.dec_func = s.g726_32_decoder
		s.bits_per_sample = 4
		break
	}

	s.bs = bitstream_state_s{
		bitstream: 0,
		residue:   0,
		lsb_first: s.packing != G726_PACKING_LEFT,
	}

	return s, nil
}

const (
	G726_ENCODING_LINEAR = 0 /* Interworking with 16 bit signed linear */
	G726_ENCODING_ULAW   = 1 /* Interworking with u-law */
	G726_ENCODING_ALAW   = 2 /* Interworking with A-law */
)

const (
	G726_PACKING_NONE  = 0
	G726_PACKING_LEFT  = 1
	G726_PACKING_RIGHT = 2
)

// G.726 state
type g726_state_t = g726_state_s
type bitstream_state_t = bitstream_state_s

type g726_decoder_func_t func(code uint8_t) int16_t

type g726_encoder_func_t func(amp int16_t) uint8_t

/*!
 * The following is the definition of the state structure
 * used by the G.726 encoder and decoder to preserve their internal
 * state between successive calls.  The meanings of the majority
 * of the state structure fields are explained in detail in the
 * ITU Recommendation G.726.  The field names are essentially indentical
 * to variable names in the bit level description of the coding algorithm
 * included in this recommendation.
 */
type g726_state_s struct {
	/*! The bit rate */
	rate int32_t
	/*! The external coding, for tandem operation */
	ext_coding int32_t
	/*! The number of bits per sample */
	bits_per_sample int32_t
	/*! One of the G.726_PACKING_xxx options */
	packing int32_t

	/*! Locked or steady state step size multiplier. */
	yl int_t // int32_t
	/*! Unlocked or non-steady state step size multiplier. */
	yu int_t // int16_t
	/*! int16_t term energy estimate. */
	dms int_t // int16_t
	/*! Long term energy estimate. */
	dml int_t // int16_t
	/*! Linear weighting coefficient of 'yl' and 'yu'. */
	ap int_t // int16_t

	/*! Coefficients of pole portion of prediction filter. */
	a [2]int_t // int16_t
	/*! Coefficients of zero portion of prediction filter. */
	b [6]int_t // int16_t
	/*! Signs of previous two samples of a partially reconstructed signal. */
	pk [2]int_t // int16_t
	/*! Previous 6 samples of the quantized difference signal represented in
	  an internal floating point format. */
	dq [6]int_t // int16_t
	/*! Previous 2 samples of the quantized difference signal represented in an
	  internal floating point format. */
	sr [2]int_t // int16_t
	/*! Delayed tone detect */
	td bool // td int

	/*! \brief The bit stream processing context. */
	bs bitstream_state_t

	/*! \brief The current encoder function. */
	enc_func g726_encoder_func_t
	/*! \brief The current decoder function. */
	dec_func g726_decoder_func_t

	packets uint64
}

// Bitstream handler state
type bitstream_state_s struct {
	/*! The bit stream. */
	bitstream uint32_t
	/*! The residual bits in bitstream. */
	residue int32_t
	/*! True if the stream is LSB first, else MSB first */
	lsb_first bool
}

type (
	uint32_t = uint32
	int32_t  = int32
	int16_t  = int16
	uint8_t  = uint8
	int_t    = int64
)

type IntValue interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func overflow[T IntValue](v T) int_t {
	return int_t(v)
}
