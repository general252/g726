

#### 提取测试pcm数据

```bash
ffmpeg -i input.mp3 -ar 8000 -ac 1 -acodec pcm_s16le -f s16le audio-samples.pcm
```

1. `-ar 8000` 设置采样率为 8kHz
2. `-ac 1` 单声道输出
3. `-acodec pcm_s16le`编码器设为 16 位小端（Little-Endian）PCM（通用兼容格式）
4. `-f s16le` 输出格式为原始 PCM 数据


#### 播放`pcm`数据
```bash
ffplay -ar 8000 -ac 1 -f s16le -i audio-samples.pcm
```

#### 播放`g726`数据
```bash
ffplay -f g726le -ar 8000 -ac 1 -code_size 4 -i audio-samples-32kbps.g726
```
1. -f #格式 小端: g726le 大端: g726
2. -ac #音频通道
3. -ar #采样率
4. -code_size #采样宽度 取值范围 2~5 分别代表 16kbps 24kbps 32kbps 40kbps


####
```bash
ffmpeg -f s16le -ar 8000 -ac 1 -i audio-samples.pcm -acodec g726 -b:a 32k -f g726 audio-samples-32kbps-ffmpeg.g726
```