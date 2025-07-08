### G726 编解码

参考  
> [spandsp](https://github.com/freeswitch/spandsp/blob/master/src/g726.c)  
> [soft-switch](https://www.soft-switch.org/downloads/)  
> [canbus/pcm_adpcm_g726](https://github.com/canbus/pcm_adpcm_g726)    
> [ffmpeg G726](https://github.com/FFmpeg/FFmpeg/blob/master/libavcodec/g726.c)    

- [x] 40kbps
- [x] 32kbps
- [x] 24kbps
- [x] 16kbps



### 示例
```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/general252/g726"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	pcmIn, err := os.ReadFile("audio-samples.pcm")
	if err != nil {
		log.Println(err)
		return
	}

	var packing = g726.PackingRight
	for i := 0; i <= 3; i++ {
		var rate = g726.Rate(i)

		pcmOut, g726Data, err := encodeAndDecode(rate, packing, pcmIn)
		if err != nil {
			log.Println(err)
			return
		} else {
			filenamePCM := fmt.Sprintf("audio-samples-re-%vkbps.pcm", (i+2)*8)
			filenameG726 := fmt.Sprintf("audio-samples-%vkbps.g726", (i+2)*8)
			log.Println(rate, filenamePCM, filenameG726)

			_ = os.WriteFile(filenamePCM, pcmOut, 0644)
			_ = os.WriteFile(filenameG726, g726Data, 0644)

		}
	}
}

func encodeAndDecode(rate g726.Rate, packing g726.PackingType, pcm []byte) (pcmOut, g726Data []byte, err error) {
	encoder := g726.G726_init_state(rate, packing)
	g726Data = encoder.EncodeV2(encoder.Pcm8ToPcm16(pcm))

	decoder := g726.G726_init_state(rate, packing)
	out := decoder.DecodeV2(g726Data)

	return decoder.Pcm16ToPcm8(out), g726Data, nil
}

```


### 测试

原始音频数据
![source](img/audio-samples.jpg)

32kbps 压缩后的音频数据
![32kbps](img/audio-samples-re-32kbps.jpg)

16kbps 压缩后的音频数据
![16kbps](img/audio-samples-re-16kbps.jpg)

  
