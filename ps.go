package mpeg

import (
	"fmt"
	"github.com/lkmio/avformat/utils"
)

const (
	PackHeaderStartCode   = 0x000001BA
	SystemHeaderStartCode = 0x000001BB
	PSMStartCode          = 0x000001BC
	ProgramEndCode        = 0x000001B9

	trickModeControlTypeFastForward = 0x0
	trickModeControlTypeSlowMotion  = 0x1
	trickModeControlTypeFreezeFrame = 0x2
	trickModeControlTypeFastReverse = 0x3
	trickModeControlTypeSlowReverse = 0x4
)

// StreamType 负载流编码器类型, ts流声明在PMT, ps流声明在psm
// Reference from https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/mpeg.h
const (
	StreamTypeVideoMPEG1     = 0x01
	StreamTypeVideoMPEG2     = 0x02
	StreamTypeAudioMPEG1     = 0x03
	StreamTypeAudioMPEG2     = 0x04
	StreamTypePrivateSection = 0x05
	StreamTypePrivateData    = 0x06
	StreamTypeAudioAAC       = 0x0F
	StreamTypeAudioMpeg4AAC  = 0x11
	StreamTypeVideoMpeg4     = 0x10
	StreamTypeVideoH264      = 0x1B
	StreamTypeVideoHEVC      = 0x24
	StreamTypeVideoCAVS      = 0x42
	StreamTypeAudioAC3       = 0x81
	StreamTypeAudioG711A     = 0x90
	StreamTypeAudioG711U     = 0x91
	StreamTypeAudioG722      = 0x92
	StreamTypeAudioG726      = 0x96
	StreamTypeAudioPCM       = 0x9C
	StreamTypeVideoSVAC      = 0x80
)

var (
	streamTypes map[int]int
	// SupportedCodecs PS流支持的编码器
	SupportedCodecs = map[utils.AVCodecID]interface{}{
		utils.AVCodecIdMP3:       StreamTypeAudioMPEG1,
		utils.AVCodecIdAAC:       StreamTypeAudioAAC,
		utils.AVCodecIdPCMALAW:   StreamTypeAudioG711A,
		utils.AVCodecIdPCMMULAW:  StreamTypeAudioG711U,
		utils.AVCodecIdADPCMG722: StreamTypeAudioG722,
		utils.AVCodecIdADPCMG726: StreamTypeAudioG726,
		utils.AVCodecIdPCMS16LE:  StreamTypeAudioPCM,

		utils.AVCodecIdH264:  StreamTypeVideoH264,
		utils.AVCodecIdHEVC:  StreamTypeVideoHEVC,
		utils.AVCodecIdMPEG4: StreamTypeVideoMpeg4,
	}

	// TSSupportedCodecs TS流支持的编码器
	TSSupportedCodecs = map[utils.AVCodecID]interface{}{
		utils.AVCodecIdAAC: StreamTypeAudioAAC,
		utils.AVCodecIdMP3: StreamTypeAudioMPEG1,

		utils.AVCodecIdH264:  StreamTypeVideoH264,
		utils.AVCodecIdHEVC:  StreamTypeVideoHEVC,
		utils.AVCodecIdMPEG4: StreamTypeVideoMpeg4,
	}
)

type StreamType int

func (s StreamType) isAudio() bool {
	return streamTypes[int(s)] == StreamIDAudio
}

func (s StreamType) isVideo() bool {
	return streamTypes[int(s)] == StreamIDVideo || streamTypes[int(s)] == StreamIDH624
}

func init() {
	streamTypes = map[int]int{
		StreamTypeVideoMPEG1:     StreamIDVideo,
		StreamTypeVideoMPEG2:     StreamIDVideo,
		StreamTypeAudioMPEG1:     StreamIDAudio,
		StreamTypeAudioMPEG2:     StreamIDAudio,
		StreamTypePrivateSection: StreamIDPrivateStream1,
		StreamTypePrivateData:    StreamIDPrivateStream1,
		StreamTypeAudioAAC:       StreamIDAudio,
		StreamTypeVideoMpeg4:     StreamIDVideo,
		StreamTypeVideoH264:      StreamIDVideo,
		StreamTypeVideoHEVC:      StreamIDVideo,
		StreamTypeVideoCAVS:      StreamIDVideo,
		StreamTypeAudioAC3:       StreamIDAudio,
	}
}

// AVCodecID2StreamType PS流AVCodecID转StreamType
func AVCodecID2StreamType(id utils.AVCodecID) (int, error) {
	return AVCodecID2StreamTypeFromCodecs(id, SupportedCodecs)
}

// StreamType2AVCodecID PS流StreamType转AVCodecID
func StreamType2AVCodecID(streamType int) (utils.AVCodecID, utils.AVMediaType, error) {
	return StreamType2AVCodecIDFromCodecs(streamType, SupportedCodecs)
}

func TSAVCodecID2StreamType(id utils.AVCodecID) (int, error) {
	return AVCodecID2StreamTypeFromCodecs(id, TSSupportedCodecs)
}

func TSStreamType2AVCodecID(streamType int) (utils.AVCodecID, utils.AVMediaType, error) {
	return StreamType2AVCodecIDFromCodecs(streamType, TSSupportedCodecs)
}

func AVCodecID2StreamTypeFromCodecs(id utils.AVCodecID, codecs map[utils.AVCodecID]interface{}) (int, error) {
	streamType, ok := codecs[id]
	if ok {
		return streamType.(int), nil
	}

	return -1, fmt.Errorf("unsupported codec: %s", id)
}

func StreamType2AVCodecIDFromCodecs(streamType int, codecs map[utils.AVCodecID]interface{}) (utils.AVCodecID, utils.AVMediaType, error) {
	for k, v := range codecs {
		if v != streamType {
			continue
		}

		var mediaType utils.AVMediaType
		if k >= utils.AVCodecIdFIRSTAUDIO {
			mediaType = utils.AVMediaTypeAudio
		} else {
			mediaType = utils.AVMediaTypeVideo
		}

		return k, mediaType, nil
	}

	return utils.AVCodecIdNONE, utils.AVMediaTypeUnknown, fmt.Errorf("the stream type %d is not implemented", streamType)
}
