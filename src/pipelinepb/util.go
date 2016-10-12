package pipelinepb

import (
	"bytes"

	"log"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/otsimo/otsimopb"
)

var (
	marshaler = jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		OrigName:     false,
	}
)

func (f *FlowIn) SetPayload(p []byte) *FlowIn {
	f.Payload = p
	return f
}
func (f *FlowIn) SetData(d *otsimopb.DataSet) *FlowIn {
	f.Data = d
	return f
}
func (f *FlowIn) SetID(key, id string) *FlowIn {
	f.Ids[key] = id
	return f
}

func (f *FlowIn) Clone() *FlowIn {
	return proto.Clone(f).(*FlowIn)
}

func (f *FlowIn) Out() *FlowOut {
	fi := proto.Clone(f).(*FlowIn)
	fo := &FlowOut{}
	fo.Ids = fi.Ids
	fo.Data = fi.Data
	fo.Payload = fi.Payload
	return fo
}
func (f *FlowIn) UnmarshalProtoPayload(pb proto.Message) error {
	return jsonpb.Unmarshal(bytes.NewReader(f.Payload), pb)
}

func (f *FlowOut) SetPayload(p []byte) *FlowOut {
	f.Payload = p
	return f
}

func (f *FlowOut) SetProtoPayload(pb proto.Message) *FlowOut {
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, pb); err != nil {
		log.Printf("failed to set proto payload, %v", err)
	}
	f.Payload = buf.Bytes()
	return f
}

func (f *FlowOut) SetData(d *otsimopb.DataSet) *FlowOut {
	f.Data = d
	return f
}

func (f *FlowOut) SetID(key, id string) *FlowOut {
	if f.Ids == nil {
		f.Ids = make(map[string]string)
	}
	f.Ids[key] = id
	return f
}

func (f *FlowOut) Clone() *FlowOut {
	return proto.Clone(f).(*FlowOut)
}

func (f *FlowOut) In() *FlowIn {
	fo := proto.Clone(f).(*FlowOut)
	fi := &FlowIn{}
	fi.Ids = fo.Ids
	fi.Data = fo.Data
	fi.Payload = fo.Payload
	return fi
}
