package g726

type g726_param_16 struct {
	qtab_723_16 [1]int
	_dqlntab    [4]int
	_witab      [4]int
	_fitab      [4]int
}

var p16 g726_param_16

func init() {
	p16.qtab_723_16 = [1]int{261}

	/*
	 * Maps G.723_16 code word to reconstructed scale factor normalized log
	 * magnitude values.  Comes from Table 11/G.726
	 */
	p16._dqlntab = [4]int{116, 365, 365, 116}

	/* Maps G.723_16 code word to log of scale factor multiplier.
	 *
	 * _witab[4] is actually {-22 , 439, 439, -22}, but FILTD wants it
	 * as WI << 5  (multiplied by 32), so we'll do that here
	 */
	p16._witab = [4]int{-704, 14048, 14048, -704}

	/*
	 * Maps G.723_16 code words to a set of values whose long and short
	 * term averages are computed and then compared to give an indication
	 * how stationary (steady state) the signal is.
	 */

	/* Comes from FUNCTF */
	p16._fitab = [4]int{0, 0xE00, 0xE00, 0}
}

func (state_ptr *G726_state) g726_16_encoder(sl int) int {
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

	d = sl - se /* d = estimation diff. */

	/* quantize prediction difference d */
	y = state_ptr.step_size()              /* quantizer step size */
	i = quantize(d, y, p16.qtab_723_16[:]) /* i = ADPCM code */

	/* Since quantize() only produces a three level output
	 * (1, 2, or 3), we must create the fourth one on our own
	 */
	if i == 3 { /* i code for the zero region */
		if (d & 0x8000) == 0 { /* If d > 0, i=3 isn't right... */
			i = 0
		}
	}

	dq = reconstruct(i&2, int(p16._dqlntab[i]), y) /* quantized diff. */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconstructed signal */

	dqsez = sr + sez - se /* pole prediction diff. */

	state_ptr.update(2, y, int(p16._witab[i]), int(p16._fitab[i]), dq, sr, dqsez)

	return i
}

func (state_ptr *G726_state) g726_16_decoder(i int) int {
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

	i &= 0x03 /* mask to get proper bits */
	sezi = state_ptr.predictor_zero()
	sez = sezi >> 1
	sei = sezi + state_ptr.predictor_pole()
	se = sei >> 1 /* se = estimated signal */

	y = state_ptr.step_size()                         /* adaptive quantizer step size */
	dq = reconstruct(i&0x02, int(p16._dqlntab[i]), y) /* unquantize pred diff */

	sr = IfElse[int](dq < 0, se-(dq&0x3FFF), se+dq) /* reconst. signal */

	dqsez = sr - se + sez /* pole prediction diff. */

	state_ptr.update(2, y, int(p16._witab[i]), int(p16._fitab[i]), dq, sr, dqsez)

	return sr << 2 /* sr was of 14-bit dynamic range */
}
