package g726

type g726_param_24 struct {
	qtab_723_24 [3]int
	_dqlntab    [8]int
	_witab      [8]int
	_fitab      [8]int
}

var p24 g726_param_24

func init() {
	p24.qtab_723_24 = [3]int{8, 218, 331}

	/*
	 * Maps G.723_24 code word to reconstructed scale factor normalized log
	 * magnitude values.
	 */
	p24._dqlntab = [8]int{-2048, 135, 273, 373, 373, 273, 135, -2048}

	/* Maps G.723_24 code word to log of scale factor multiplier. */
	p24._witab = [8]int{-128, 960, 4384, 18624, 18624, 4384, 960, -128}

	/*
	 * Maps G.723_24 code words to a set of values whose long and short
	 * term averages are computed and then compared to give an indication
	 * how stationary (steady state) the signal is.
	 */
	p24._fitab = [8]int{0, 0x200, 0x400, 0xE00, 0xE00, 0x400, 0x200, 0}
}

func (state_ptr *G726_state) g726_24_encoder(sl int) int {
	var (
		sezi  int
		sei   int
		sez   int
		se    int
		d     int
		y     int
		i     int
		dq    int
		sr    int
		dqsez int
	)

	sl >>= 2 /* sl of 14-bit dynamic range */

	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	sei = sezi + state_ptr.predictor_pole()
	se = sei >> 1 /* se = estimated signal */

	d = sl - se /* d = estimation diff. */

	/* quantize prediction difference d */
	y = state_ptr.step_size()                      /* quantizer step size */
	i = quantize(d, y, p24.qtab_723_24[:])         /* i = ADPCM code */
	dq = reconstruct(i&4, int(p24._dqlntab[i]), y) /* quantized diff. */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconstructed signal */

	dqsez = sr + sez - se /* pole prediction diff. */

	state_ptr.update(3, y, int(p24._witab[i]), int(p24._fitab[i]), dq, sr, dqsez)

	return i
}

func (state_ptr *G726_state) g726_24_decoder(i int) int {
	var (
		sezi  int
		sez   int
		sei   int
		se    int
		y     int
		dq    int
		sr    int
		dqsez int
	)

	i &= 0x07 /* mask to get proper bits */
	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	sei = sezi + state_ptr.predictor_pole()
	se = sei >> 1 /* se = estimated signal */

	y = state_ptr.step_size()                         /* adaptive quantizer step size */
	dq = reconstruct(i&0x04, int(p24._dqlntab[i]), y) /* unquantize pred diff */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconst. signal */

	dqsez = sr - se + sez /* pole prediction diff. */

	state_ptr.update(3, y, int(p24._witab[i]), int(p24._fitab[i]), dq, sr, dqsez)

	return sr << 2 /* sr was of 14-bit dynamic range */
}
