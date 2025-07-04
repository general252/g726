

#### 提取测试pcm数据

```bash
ffmpeg -i input.mp3 -ar 8000 -ac 1 -acodec pcm_s16le -f s16le audio-samples.pcm
```

1. `-ar 8000` 设置采样率为 8kHz
2. `-ac 1` 单声道输出
3. `-acodec pcm_s16le`编码器设为 16 位小端（Little-Endian）PCM（通用兼容格式）
4. `-f s16le` 输出格式为原始 PCM 数据


#### 播放pcm数据
```bash
ffplay -ar 8000 -ac 1 -f s16le -i audio-samples.pcm
```

