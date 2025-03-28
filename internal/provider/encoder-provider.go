package provider

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	StrategyGob     = "gob"
	StrategyJSON    = "json"
	StrategyMsgPack = "msgpack"
)

type IEncoder interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

// Encoder factory
func EncoderFactory(strategy string) (IEncoder, error) {
	switch strategy {
	case StrategyGob:
		return &GobProvider{}, nil
	case StrategyJSON:
		return &JsonProvider{}, nil
	case StrategyMsgPack:
		return &MsgPackProvider{}, nil
	default:
		return nil, errors.New("bsp: unknown encoding strategy")
	}
}

// Gob provider
type GobProvider struct{}

func (g *GobProvider) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GobProvider) Decode(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}

// Json provider
type JsonProvider struct{}

func (j *JsonProvider) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonProvider) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// MsgPack provider
type MsgPackProvider struct{}

func (m *MsgPackProvider) Encode(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *MsgPackProvider) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
