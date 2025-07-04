package g726

type g726_param_40 struct {
	qtab_723_40 [15]int
	_dqlntab    [32]int16
	_witab      [32]int16
	_fitab      [32]int16
}

var p40 g726_param_40

func init() {
	p40.qtab_723_40 = [15]int{-122, -16, 68, 139, 198, 250, 298, 339,
		378, 413, 445, 475, 502, 528, 553}
	/*
	 * Maps G.723_40 code word to ructeconstructed scale factor normalized log
	 * magnitude values.
	 */
	p40._dqlntab = [32]int16{-2048, -66, 28, 104, 169, 224, 274, 318,
		358, 395, 429, 459, 488, 514, 539, 566,
		566, 539, 514, 488, 459, 429, 395, 358,
		318, 274, 224, 169, 104, 28, -66, -2048}

	/* Maps G.723_40 code word to log of scale factor multiplier. */
	p40._witab = [32]int16{448, 448, 768, 1248, 1280, 1312, 1856, 3200,
		4512, 5728, 7008, 8960, 11456, 14080, 16928, 22272,
		22272, 16928, 14080, 11456, 8960, 7008, 5728, 4512,
		3200, 1856, 1312, 1280, 1248, 768, 448, 448}

	/*
	 * Maps G.723_40 code words to a set of values whose long and short
	 * term averages are computed and then compared to give an indication
	 * how stationary (steady state) the signal is.
	 */
	p40._fitab = [32]int16{0, 0, 0, 0, 0, 0x200, 0x200, 0x200,
		0x200, 0x200, 0x400, 0x600, 0x800, 0xA00, 0xC00, 0xC00,
		0xC00, 0xC00, 0xA00, 0x800, 0x600, 0x400, 0x200, 0x200,
		0x200, 0x200, 0x200, 0, 0, 0, 0, 0}

}

func (state_ptr *G726_state) g726_40_encoder(sl int) int {

	var (
		sezi  int
		sez   int
		sei   int
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

	d = sl - se /* d = estimation difference */

	/* quantize prediction difference */
	y = state_ptr.step_size()              /* adaptive quantizer step size */
	i = quantize(d, y, p40.qtab_723_40[:]) /* i = ADPCM code */

	dq = reconstruct(i&0x10, int(p40._dqlntab[i]), y) /* quantized diff */

	sr = IfElse[int](dq < 0, se-(dq&0x7FFF), se+dq) /* reconstructed signal */

	dqsez = sr + sez - se /* dqsez = pole prediction diff. */

	state_ptr.update(5, y, int(p40._witab[i]), int(p40._fitab[i]), dq, sr, dqsez)

	return i
}

func (state_ptr *G726_state) g726_40_decoder(i int) int {
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

	i &= 0x1f /* mask to get proper bits */
	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	sei = sezi + state_ptr.predictor_pole()
	se = sei >> 1 /* se = estimated signal */

	y = state_ptr.step_size()                         /* adaptive quantizer step size */
	dq = reconstruct(i&0x10, int(p40._dqlntab[i]), y) /* estimation diff. */

	sr = IfElse[int](dq < 0, se-(dq&0x7FFF), se+dq) /* reconst. signal */

	dqsez = sr - se + sez /* pole prediction diff. */

	state_ptr.update(5, y, int(p40._witab[i]), int(p40._fitab[i]), dq, sr, dqsez)

	return sr << 2 /* sr was of 14-bit dynamic range */
}
