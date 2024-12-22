package zwf

import (
	"encoding/binary"
	"io"

	"github.com/go-audio/audio"
	"github.com/namgo/GameWaveFans/pkg/common"
)

func convertData(d []int) []byte {
	data := make([]byte, len(d)*2)
	for i, sample := range d {
		unpackedSample := make([]byte, 2)
		binary.BigEndian.PutUint16(unpackedSample, uint16(sample))
		data[i*2] = unpackedSample[0]
		data[(i*2)+1] = unpackedSample[1]
	}

	return data
}

func Encode(w io.Writer, buf *audio.IntBuffer) error {
	sampleCount := make([]byte, 4)
	binary.LittleEndian.PutUint32(sampleCount, uint32(len(buf.Data)))

	unpackedSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(unpackedSize, uint32(len(buf.Data)*2))

	data := convertData(buf.Data)
	packedData, err := common.WriteZlibToBuffer(data)
	if err != nil {
		return err
	}

	packedSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(packedSize, uint32(len(packedData)))

	w.Write([]byte("\x02\xee\x90\x7c"))
	w.Write(sampleCount)
	w.Write([]byte("\x01\x00\x00\x00"))
	w.Write(packedSize)
	w.Write(unpackedSize)
	w.Write(packedData)
	return nil
}
