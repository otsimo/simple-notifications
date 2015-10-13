// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Some tagging information for error messages.
const (
	msgBadDesc            = "Unrecognized descriptor byte"
	msgDecCannotExpandArr = "cannot expand go array from %v to stream length: %v"
)

var (
	onlyMapOrArrayCanDecodeIntoStructErr = errors.New("only encoded map or array can be decoded into a struct")
	cannotDecodeIntoNilErr               = errors.New("cannot decode into nil")
)

// decReader abstracts the reading source, allowing implementations that can
// read from an io.Reader or directly off a byte slice with zero-copying.
type decReader interface {
	unreadn1()

	// readx will use the implementation scratch buffer if possible i.e. n < len(scratchbuf), OR
	// just return a view of the []byte being decoded from.
	// Ensure you call detachZeroCopyBytes later if this needs to be sent outside codec control.
	readx(n int) []byte
	readb([]byte)
	readn1() uint8
	readn1eof() (v uint8, eof bool)
	numread() int // number of bytes read
	track()
	stopTrack() []byte
}

type decReaderByteScanner interface {
	io.Reader
	io.ByteScanner
}

type decDriver interface {
	// this will check if the next token is a break.
	CheckBreak() bool
	TryDecodeAsNil() bool
	// check if a container type: vt is one of: Bytes, String, Nil, Slice or Map.
	// if vt param == valueTypeNil, and nil is seen in stream, consume the nil.
	IsContainerType(vt valueType) bool
	IsBuiltinType(rt uintptr) bool
	DecodeBuiltin(rt uintptr, v interface{})
	//decodeNaked: Numbers are decoded as int64, uint64, float64 only (no smaller sized number types).
	//for extensions, decodeNaked must completely decode them as a *RawExt.
	//extensions should also use readx to decode them, for efficiency.
	//kInterface will extract the detached byte slice if it has to pass it outside its realm.
	DecodeNaked() (v interface{}, vt valueType, decodeFurther bool)
	DecodeInt(bitsize uint8) (i int64)
	DecodeUint(bitsize uint8) (ui uint64)
	DecodeFloat(chkOverflow32 bool) (f float64)
	DecodeBool() (b bool)
	// DecodeString can also decode symbols.
	// It looks redundant as DecodeBytes is available.
	// However, some codecs (e.g. binc) support symbols and can
	// return a pre-stored string value, meaning that it can bypass
	// the cost of []byte->string conversion.
	DecodeString() (s string)

	// DecodeBytes may be called directly, without going through reflection.
	// Consequently, it must be designed to handle possible nil.
	DecodeBytes(bs []byte, isstring, zerocopy bool) (bsOut []byte)

	// decodeExt will decode into a *RawExt or into an extension.
	DecodeExt(v interface{}, xtag uint64, ext Ext) (realxtag uint64)
	// decodeExt(verifyTag bool, tag byte) (xtag byte, xbs []byte)
	ReadMapStart() int
	ReadArrayStart() int
	// ReadEnd registers the end of a map or array.
	ReadEnd()
}

type decNoSeparator struct{}

func (_ decNoSeparator) ReadEnd() {}

type DecodeOptions struct {
	// MapType specifies type to use during schema-less decoding of a map in the stream.
	// If nil, we use map[interface{}]interface{}
	MapType reflect.Type

	// SliceType specifies type to use during schema-less decoding of an array in the stream.
	// If nil, we use []interface{}
	SliceType reflect.Type

	// If ErrorIfNoField, return an error when decoding a map
	// from a codec stream into a struct, and no matching struct field is found.
	ErrorIfNoField bool

	// If ErrorIfNoArrayExpand, return an error when decoding a slice/array that cannot be expanded.
	// For example, the stream contains an array of 8 items, but you are decoding into a [4]T array,
	// or you are decoding into a slice of length 4 which is non-addressable (and so cannot be set).
	ErrorIfNoArrayExpand bool

	// If SignedInteger, use the int64 during schema-less decoding of unsigned values (not uint64).
	SignedInteger bool

	// MaxInitLen defines the initial length that we "make" a collection (slice, chan or map) with.
	// If 0 or negative, we default to a sensible value based on the size of an element in the collection.
	//
	// For example, when decoding, a stream may say that it has MAX_UINT elements.
	// We should not auto-matically provision a slice of that length, to prevent Out-Of-Memory crash.
	// Instead, we provision up to MaxInitLen, fill that up, and start appending after that.
	MaxInitLen int
}

// ------------------------------------

// ioDecByteScanner implements Read(), ReadByte(...), UnreadByte(...) methods
// of io.Reader, io.ByteScanner.
type ioDecByteScanner struct {
	r  io.Reader
	l  byte    // last byte
	ls byte    // last byte status. 0: init-canDoNothing, 1: canRead, 2: canUnread
	b  [1]byte // tiny buffer for reading single bytes
}

func (z *ioDecByteScanner) Read(p []byte) (n int, err error) {
	var firstByte bool
	if z.ls == 1 {
		z.ls = 2
		p[0] = z.l
		if len(p) == 1 {
			n = 1
			return
		}
		firstByte = true
		p = p[1:]
	}
	n, err = z.r.Read(p)
	if n > 0 {
		if err == io.EOF && n == len(p) {
			err = nil // read was successful, so postpone EOF (till next time)
		}
		z.l = p[n-1]
		z.ls = 2
	}
	if firstByte {
		n++
	}
	return
}

func (z *ioDecByteScanner) ReadByte() (c byte, err error) {
	n, err := z.Read(z.b[:])
	if n == 1 {
		c = z.b[0]
		if err == io.EOF {
			err = nil // read was successful, so postpone EOF (till next time)
		}
	}
	return
}

func (z *ioDecByteScanner) UnreadByte() (err error) {
	x := z.ls
	if x == 0 {
		err = errors.New("cannot unread - nothing has been read")
	} else if x == 1 {
		err = errors.New("cannot unread - last byte has not been read")
	} else if x == 2 {
		z.ls = 1
	}
	return
}

// ioDecReader is a decReader that reads off an io.Reader
type ioDecReader struct {
	br decReaderByteScanner
	// temp byte array re-used internally for efficiency during read.
	// shares buffer with Decoder, so we keep size of struct within 8 words.
	x   *[scratchByteArrayLen]byte
	bs  ioDecByteScanner
	n   int    // num read
	tr  []byte // tracking bytes read
	trb bool
}

func (z *ioDecReader) numread() int {
	return z.n
}

func (z *ioDecReader) readx(n int) (bs []byte) {
	if n <= 0 {
		return
	}
	if n < len(z.x) {
		bs = z.x[:n]
	} else {
		bs = make([]byte, n)
	}
	if _, err := io.ReadAtLeast(z.br, bs, n); err != nil {
		panic(err)
	}
	z.n += len(bs)
	if z.trb {
		z.tr = append(z.tr, bs...)
	}
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	n, err := io.ReadAtLeast(z.br, bs, len(bs))
	z.n += n
	if err != nil {
		panic(err)
	}
	if z.trb {
		z.tr = append(z.tr, bs...)
	}
}

func (z *ioDecReader) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	if err != nil {
		panic(err)
	}
	z.n++
	if z.trb {
		z.tr = append(z.tr, b)
	}
	return b
}

func (z *ioDecReader) readn1eof() (b uint8, eof bool) {
	b, err := z.br.ReadByte()
	if err == nil {
		z.n++
		if z.trb {
			z.tr = append(z.tr, b)
		}
	} else if err == io.EOF {
		eof = true
	} else {
		panic(err)
	}
	return
}

func (z *ioDecReader) unreadn1() {
	err := z.br.UnreadByte()
	if err != nil {
		panic(err)
	}
	z.n--
	if z.trb {
		if l := len(z.tr) - 1; l >= 0 {
			z.tr = z.tr[:l]
		}
	}
}

func (z *ioDecReader) track() {
	if z.tr != nil {
		z.tr = z.tr[:0]
	}
	z.trb = true
}

func (z *ioDecReader) stopTrack() (bs []byte) {
	z.trb = false
	return z.tr
}

// ------------------------------------

var bytesDecReaderCannotUnreadErr = errors.New("cannot unread last byte read")

// bytesDecReader is a decReader that reads off a byte slice with zero copying
type bytesDecReader struct {
	b []byte // data
	c int    // cursor
	a int    // available
	t int    // track start
}

func (z *bytesDecReader) numread() int {
	return z.c
}

func (z *bytesDecReader) unreadn1() {
	if z.c == 0 || len(z.b) == 0 {
		panic(bytesDecReaderCannotUnreadErr)
	}
	z.c--
	z.a++
	return
}

func (z *bytesDecReader) readx(n int) (bs []byte) {
	// slicing from a non-constant start position is more expensive,
	// as more computation is required to decipher the pointer start position.
	// However, we do it only once, and it's better than reslicing both z.b and return value.

	if n <= 0 {
	} else if z.a == 0 {
		panic(io.EOF)
	} else if n > z.a {
		panic(io.ErrUnexpectedEOF)
	} else {
		c0 := z.c
		z.c = c0 + n
		z.a = z.a - n
		bs = z.b[c0:z.c]
	}
	return
}

func (z *bytesDecReader) readn1() (v uint8) {
	if z.a == 0 {
		panic(io.EOF)
	}
	v = z.b[z.c]
	z.c++
	z.a--
	return
}

func (z *bytesDecReader) readn1eof() (v uint8, eof bool) {
	if z.a == 0 {
		eof = true
		return
	}
	v = z.b[z.c]
	z.c++
	z.a--
	return
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(len(bs)))
}

func (z *bytesDecReader) track() {
	z.t = z.c
}

func (z *bytesDecReader) stopTrack() (bs []byte) {
	return z.b[z.t:z.c]
}

// ------------------------------------

type decFnInfo struct {
	d     *Decoder
	ti    *typeInfo
	xfFn  Ext
	xfTag uint64
	seq   seqType
}

// ----------------------------------------

type decFn struct {
	i decFnInfo
	f func(*decFnInfo, reflect.Value)
}

func (f *decFnInfo) builtin(rv reflect.Value) {
	f.d.d.DecodeBuiltin(f.ti.rtid, rv.Addr().Interface())
}

func (f *decFnInfo) rawExt(rv reflect.Value) {
	f.d.d.DecodeExt(rv.Addr().Interface(), 0, nil)
}

func (f *decFnInfo) ext(rv reflect.Value) {
	f.d.d.DecodeExt(rv.Addr().Interface(), f.xfTag, f.xfFn)
}

func (f *decFnInfo) getValueForUnmarshalInterface(rv reflect.Value, indir int8) (v interface{}) {
	if indir == -1 {
		v = rv.Addr().Interface()
	} else if indir == 0 {
		v = rv.Interface()
	} else {
		for j := int8(0); j < indir; j++ {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		v = rv.Interface()
	}
	return
}

func (f *decFnInfo) selferUnmarshal(rv reflect.Value) {
	f.getValueForUnmarshalInterface(rv, f.ti.csIndir).(Selfer).CodecDecodeSelf(f.d)
}

func (f *decFnInfo) binaryUnmarshal(rv reflect.Value) {
	bm := f.getValueForUnmarshalInterface(rv, f.ti.bunmIndir).(encoding.BinaryUnmarshaler)
	xbs := f.d.d.DecodeBytes(nil, false, true)
	if fnerr := bm.UnmarshalBinary(xbs); fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) textUnmarshal(rv reflect.Value) {
	tm := f.getValueForUnmarshalInterface(rv, f.ti.tunmIndir).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(f.d.d.DecodeBytes(f.d.b[:], true, true))
	if fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) jsonUnmarshal(rv reflect.Value) {
	tm := f.getValueForUnmarshalInterface(rv, f.ti.junmIndir).(jsonUnmarshaler)
	// bs := f.d.d.DecodeBytes(f.d.b[:], true, true)
	// grab the bytes to be read, as UnmarshalJSON wants the full JSON to unmarshal it itself.
	f.d.r.track()
	f.d.swallow()
	bs := f.d.r.stopTrack()
	// fmt.Printf(">>>>>> REFLECTION JSON: %s\n", bs)
	fnerr := tm.UnmarshalJSON(bs)
	if fnerr != nil {
		panic(fnerr)
	}
}

func (f *decFnInfo) kErr(rv reflect.Value) {
	f.d.errorf("no decoding function defined for kind %v", rv.Kind())
}

func (f *decFnInfo) kString(rv reflect.Value) {
	rv.SetString(f.d.d.DecodeString())
}

func (f *decFnInfo) kBool(rv reflect.Value) {
	rv.SetBool(f.d.d.DecodeBool())
}

func (f *decFnInfo) kInt(rv reflect.Value) {
	rv.SetInt(f.d.d.DecodeInt(intBitsize))
}

func (f *decFnInfo) kInt64(rv reflect.Value) {
	rv.SetInt(f.d.d.DecodeInt(64))
}

func (f *decFnInfo) kInt32(rv reflect.Value) {
	rv.SetInt(f.d.d.DecodeInt(32))
}

func (f *decFnInfo) kInt8(rv reflect.Value) {
	rv.SetInt(f.d.d.DecodeInt(8))
}

func (f *decFnInfo) kInt16(rv reflect.Value) {
	rv.SetInt(f.d.d.DecodeInt(16))
}

func (f *decFnInfo) kFloat32(rv reflect.Value) {
	rv.SetFloat(f.d.d.DecodeFloat(true))
}

func (f *decFnInfo) kFloat64(rv reflect.Value) {
	rv.SetFloat(f.d.d.DecodeFloat(false))
}

func (f *decFnInfo) kUint8(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(8))
}

func (f *decFnInfo) kUint64(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(64))
}

func (f *decFnInfo) kUint(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(uintBitsize))
}

func (f *decFnInfo) kUintptr(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(uintBitsize))
}

func (f *decFnInfo) kUint32(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(32))
}

func (f *decFnInfo) kUint16(rv reflect.Value) {
	rv.SetUint(f.d.d.DecodeUint(16))
}

// func (f *decFnInfo) kPtr(rv reflect.Value) {
// 	debugf(">>>>>>> ??? decode kPtr called - shouldn't get called")
// 	if rv.IsNil() {
// 		rv.Set(reflect.New(rv.Type().Elem()))
// 	}
// 	f.d.decodeValue(rv.Elem())
// }

// var kIntfCtr uint64

func (f *decFnInfo) kInterfaceNaked() (rvn reflect.Value) {
	// nil interface:
	// use some hieristics to decode it appropriately
	// based on the detected next value in the stream.
	v, vt, decodeFurther := f.d.d.DecodeNaked()
	if vt == valueTypeNil {
		return
	}
	// We cannot decode non-nil stream value into nil interface with methods (e.g. io.Reader).
	if num := f.ti.rt.NumMethod(); num > 0 {
		f.d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, num)
		return
	}
	var useRvn bool
	switch vt {
	case valueTypeMap:
		if f.d.h.MapType == nil {
			var m2 map[interface{}]interface{}
			v = &m2
		} else {
			rvn = reflect.New(f.d.h.MapType).Elem()
			useRvn = true
		}
	case valueTypeArray:
		if f.d.h.SliceType == nil {
			var m2 []interface{}
			v = &m2
		} else {
			rvn = reflect.New(f.d.h.SliceType).Elem()
			useRvn = true
		}
	case valueTypeExt:
		re := v.(*RawExt)
		bfn := f.d.h.getExtForTag(re.Tag)
		if bfn == nil {
			re.Data = detachZeroCopyBytes(f.d.bytes, nil, re.Data)
			rvn = reflect.ValueOf(*re)
		} else {
			rvnA := reflect.New(bfn.rt)
			rvn = rvnA.Elem()
			if re.Data != nil {
				bfn.ext.ReadExt(rvnA.Interface(), re.Data)
			} else {
				bfn.ext.UpdateExt(rvnA.Interface(), re.Value)
			}
		}
		return
	}
	if decodeFurther {
		if useRvn {
			f.d.decodeValue(rvn, nil)
		} else if v != nil {
			// this v is a pointer, so we need to dereference it when done
			f.d.decode(v)
			rvn = reflect.ValueOf(v).Elem()
			useRvn = true
		}
	}

	if !useRvn && v != nil {
		rvn = reflect.ValueOf(v)
	}
	return
}

func (f *decFnInfo) kInterface(rv reflect.Value) {
	// debugf("\t===> kInterface")

	// Note:
	// A consequence of how kInterface works, is that
	// if an interface already contains something, we try
	// to decode into what was there before.
	// We do not replace with a generic value (as got from decodeNaked).

	if rv.IsNil() {
		rvn := f.kInterfaceNaked()
		if rvn.IsValid() {
			rv.Set(rvn)
		}
	} else {
		rve := rv.Elem()
		// Note: interface{} is settable, but underlying type may not be.
		// Consequently, we have to set the reflect.Value directly.
		// if underlying type is settable (e.g. ptr or interface),
		// we just decode into it.
		// Else we create a settable value, decode into it, and set on the interface.
		if rve.CanSet() {
			f.d.decodeValue(rve, nil)
		} else {
			rve2 := reflect.New(rve.Type()).Elem()
			rve2.Set(rve)
			f.d.decodeValue(rve2, nil)
			rv.Set(rve2)
		}
	}
}

func (f *decFnInfo) kStruct(rv reflect.Value) {
	fti := f.ti
	d := f.d
	dd := d.d
	if dd.IsContainerType(valueTypeMap) {
		containerLen := dd.ReadMapStart()
		if containerLen == 0 {
			dd.ReadEnd()
			return
		}
		tisfi := fti.sfi
		hasLen := containerLen >= 0
		if hasLen {
			for j := 0; j < containerLen; j++ {
				// rvkencname := dd.DecodeString()
				rvkencname := stringView(dd.DecodeBytes(f.d.b[:], true, true))
				// rvksi := ti.getForEncName(rvkencname)
				if k := fti.indexForEncName(rvkencname); k > -1 {
					si := tisfi[k]
					if dd.TryDecodeAsNil() {
						si.setToZeroValue(rv)
					} else {
						d.decodeValue(si.field(rv, true), nil)
					}
				} else {
					d.structFieldNotFound(-1, rvkencname)
				}
			}
		} else {
			for j := 0; !dd.CheckBreak(); j++ {
				// rvkencname := dd.DecodeString()
				rvkencname := stringView(dd.DecodeBytes(f.d.b[:], true, true))
				// rvksi := ti.getForEncName(rvkencname)
				if k := fti.indexForEncName(rvkencname); k > -1 {
					si := tisfi[k]
					if dd.TryDecodeAsNil() {
						si.setToZeroValue(rv)
					} else {
						d.decodeValue(si.field(rv, true), nil)
					}
				} else {
					d.structFieldNotFound(-1, rvkencname)
				}
			}
			dd.ReadEnd()
		}
	} else if dd.IsContainerType(valueTypeArray) {
		containerLen := dd.ReadArrayStart()
		if containerLen == 0 {
			dd.ReadEnd()
			return
		}
		// Not much gain from doing it two ways for array.
		// Arrays are not used as much for structs.
		hasLen := containerLen >= 0
		for j, si := range fti.sfip {
			if hasLen {
				if j == containerLen {
					break
				}
			} else if dd.CheckBreak() {
				break
			}
			if dd.TryDecodeAsNil() {
				si.setToZeroValue(rv)
			} else {
				d.decodeValue(si.field(rv, true), nil)
			}
		}
		if containerLen > len(fti.sfip) {
			// read remaining values and throw away
			for j := len(fti.sfip); j < containerLen; j++ {
				d.structFieldNotFound(j, "")
			}
		}
		dd.ReadEnd()
	} else {
		f.d.error(onlyMapOrArrayCanDecodeIntoStructErr)
		return
	}
}

func (f *decFnInfo) kSlice(rv reflect.Value) {
	// A slice can be set from a map or array in stream.
	// This way, the order can be kept (as order is lost with map).
	ti := f.ti
	d := f.d
	dd := d.d
	rtelem0 := ti.rt.Elem()

	if dd.IsContainerType(valueTypeBytes) || dd.IsContainerType(valueTypeString) {
		if ti.rtid == uint8SliceTypId || rtelem0.Kind() == reflect.Uint8 {
			if f.seq == seqTypeChan {
				bs2 := dd.DecodeBytes(nil, false, true)
				ch := rv.Interface().(chan<- byte)
				for _, b := range bs2 {
					ch <- b
				}
			} else {
				rvbs := rv.Bytes()
				bs2 := dd.DecodeBytes(rvbs, false, false)
				if rvbs == nil && bs2 != nil || rvbs != nil && bs2 == nil || len(bs2) != len(rvbs) {
					if rv.CanSet() {
						rv.SetBytes(bs2)
					} else {
						copy(rvbs, bs2)
					}
				}
			}
			return
		}
	}

	// array := f.seq == seqTypeChan

	slh, containerLenS := d.decSliceHelperStart()

	var rvlen, numToRead int
	var truncated bool // says that the len of the sequence is not same as the expected number of elements.

	numToRead = containerLenS // if truncated, reset numToRead

	// an array can never return a nil slice. so no need to check f.array here.
	if rv.IsNil() {
		// either chan or slice
		if rvlen, truncated = decInferLen(containerLenS, f.d.h.MaxInitLen, int(rtelem0.Size())); truncated {
			numToRead = rvlen
		}
		if f.seq == seqTypeSlice {
			rv.Set(reflect.MakeSlice(ti.rt, rvlen, rvlen))
		} else if f.seq == seqTypeChan {
			rv.Set(reflect.MakeChan(ti.rt, rvlen))
		}
	} else {
		rvlen = rv.Len()
	}

	if containerLenS == 0 {
		if f.seq == seqTypeSlice && rvlen != 0 {
			rv.SetLen(0)
		}
		// dd.ReadEnd()
		return
	}

	rtelem := rtelem0
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}
	fn := d.getDecFn(rtelem, true, true)

	var rv0, rv9 reflect.Value
	rv0 = rv
	rvChanged := false

	rvcap := rv.Cap()

	// for j := 0; j < containerLenS; j++ {

	if containerLenS >= 0 { // hasLen
		if f.seq == seqTypeChan {
			// handle chan specially:
			for j := 0; j < containerLenS; j++ {
				rv9 = reflect.New(rtelem0).Elem()
				d.decodeValue(rv9, fn)
				rv.Send(rv9)
			}
		} else { // slice or array
			if containerLenS > rvcap {
				if f.seq == seqTypeArray {
					d.arrayCannotExpand(rvlen, containerLenS)
				} else {
					oldRvlenGtZero := rvlen > 0
					rvlen, truncated = decInferLen(containerLenS, f.d.h.MaxInitLen, int(rtelem0.Size()))
					rv = reflect.MakeSlice(ti.rt, rvlen, rvlen)
					if oldRvlenGtZero && !isImmutableKind(rtelem0.Kind()) {
						reflect.Copy(rv, rv0) // only copy up to length NOT cap i.e. rv0.Slice(0, rvcap)
					}
					rvChanged = true
				}
				numToRead = rvlen
			} else if containerLenS != rvlen {
				if f.seq == seqTypeSlice {
					rv.SetLen(containerLenS)
					rvlen = containerLenS
				}
			}
			j := 0
			// we read up to the numToRead
			for ; j < numToRead; j++ {
				d.decodeValue(rv.Index(j), fn)
			}

			// if slice, expand and read up to containerLenS (or EOF) iff truncated
			// if array, swallow all the rest.

			if f.seq == seqTypeArray {
				for ; j < containerLenS; j++ {
					d.swallow()
				}
			} else if truncated { // slice was truncated, as chan NOT in this block
				for ; j < containerLenS; j++ {
					rv = expandSliceValue(rv, 1)
					rv9 = rv.Index(j)
					if resetSliceElemToZeroValue {
						rv9.Set(reflect.Zero(rtelem0))
					}
					d.decodeValue(rv9, fn)
				}
			}
		}
	} else {
		for j := 0; !dd.CheckBreak(); j++ {
			var decodeIntoBlank bool
			// if indefinite, etc, then expand the slice if necessary
			if j >= rvlen {
				if f.seq == seqTypeArray {
					d.arrayCannotExpand(rvlen, j+1)
					decodeIntoBlank = true
				} else if f.seq == seqTypeSlice {
					// rv = reflect.Append(rv, reflect.Zero(rtelem0)) // uses append logic, plus varargs
					rv = expandSliceValue(rv, 1)
					rv9 = rv.Index(j)
					// rv.Index(rv.Len() - 1).Set(reflect.Zero(rtelem0))
					if resetSliceElemToZeroValue {
						rv9.Set(reflect.Zero(rtelem0))
					}
					rvlen++
					rvChanged = true
				}
			} else if f.seq != seqTypeChan { // slice or array
				rv9 = rv.Index(j)
			}
			if f.seq == seqTypeChan {
				rv9 = reflect.New(rtelem0).Elem()
				d.decodeValue(rv9, fn)
				rv.Send(rv9)
			} else if decodeIntoBlank {
				d.swallow()
			} else { // seqTypeSlice
				d.decodeValue(rv9, fn)
			}
		}
		slh.End()
	}

	if rvChanged {
		rv0.Set(rv)
	}
}

func (f *decFnInfo) kArray(rv reflect.Value) {
	// f.d.decodeValue(rv.Slice(0, rv.Len()))
	f.kSlice(rv.Slice(0, rv.Len()))
}

func (f *decFnInfo) kMap(rv reflect.Value) {
	d := f.d
	dd := d.d
	containerLen := dd.ReadMapStart()

	ti := f.ti
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(ti.rt))
	}

	if containerLen == 0 {
		// It is not length-prefix style container. They have no End marker.
		// dd.ReadMapEnd()
		return
	}

	ktype, vtype := ti.rt.Key(), ti.rt.Elem()
	ktypeId := reflect.ValueOf(ktype).Pointer()
	var keyFn, valFn *decFn
	var xtyp reflect.Type
	for xtyp = ktype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() {
	}
	keyFn = d.getDecFn(xtyp, true, true)
	for xtyp = vtype; xtyp.Kind() == reflect.Ptr; xtyp = xtyp.Elem() {
	}
	valFn = d.getDecFn(xtyp, true, true)
	// for j := 0; j < containerLen; j++ {
	if containerLen > 0 {
		for j := 0; j < containerLen; j++ {
			rvk := reflect.New(ktype).Elem()
			d.decodeValue(rvk, keyFn)

			// special case if a byte array.
			if ktypeId == intfTypId {
				rvk = rvk.Elem()
				if rvk.Type() == uint8SliceTyp {
					rvk = reflect.ValueOf(string(rvk.Bytes()))
				}
			}
			rvv := rv.MapIndex(rvk)
			// TODO: is !IsValid check required?
			if !rvv.IsValid() {
				rvv = reflect.New(vtype).Elem()
			}
			d.decodeValue(rvv, valFn)
			rv.SetMapIndex(rvk, rvv)
		}
	} else {
		for j := 0; !dd.CheckBreak(); j++ {
			rvk := reflect.New(ktype).Elem()
			d.decodeValue(rvk, keyFn)

			// special case if a byte array.
			if ktypeId == intfTypId {
				rvk = rvk.Elem()
				if rvk.Type() == uint8SliceTyp {
					rvk = reflect.ValueOf(string(rvk.Bytes()))
				}
			}
			rvv := rv.MapIndex(rvk)
			if !rvv.IsValid() {
				rvv = reflect.New(vtype).Elem()
			}
			d.decodeValue(rvv, valFn)
			rv.SetMapIndex(rvk, rvv)
		}
		dd.ReadEnd()
	}
}

type decRtidFn struct {
	rtid uintptr
	fn   decFn
}

// A Decoder reads and decodes an object from an input stream in the codec format.
type Decoder struct {
	// hopefully, reduce derefencing cost by laying the decReader inside the Decoder.
	// Try to put things that go together to fit within a cache line (8 words).

	d decDriver
	// NOTE: Decoder shouldn't call it's read methods,
	// as the handler MAY need to do some coordination.
	r decReader
	// sa [initCollectionCap]decRtidFn
	s []decRtidFn
	h *BasicHandle

	rb    bytesDecReader
	hh    Handle
	be    bool // is binary encoding
	bytes bool // is bytes reader
	js    bool // is json handle

	ri ioDecReader
	f  map[uintptr]*decFn
	// _  uintptr // for alignment purposes, so next one starts from a cache line

	b [scratchByteArrayLen]byte
}

// NewDecoder returns a Decoder for decoding a stream of bytes from an io.Reader.
//
// For efficiency, Users are encouraged to pass in a memory buffered reader
// (eg bufio.Reader, bytes.Buffer).
func NewDecoder(r io.Reader, h Handle) (d *Decoder) {
	d = &Decoder{hh: h, h: h.getBasicHandle(), be: h.isBinary()}
	// d.s = d.sa[:0]
	d.ri.x = &d.b
	d.ri.bs.r = r
	var ok bool
	d.ri.br, ok = r.(decReaderByteScanner)
	if !ok {
		d.ri.br = &d.ri.bs
	}
	d.r = &d.ri
	_, d.js = h.(*JsonHandle)
	d.d = h.newDecDriver(d)
	return
}

// NewDecoderBytes returns a Decoder which efficiently decodes directly
// from a byte slice with zero copying.
func NewDecoderBytes(in []byte, h Handle) (d *Decoder) {
	d = &Decoder{hh: h, h: h.getBasicHandle(), be: h.isBinary(), bytes: true}
	// d.s = d.sa[:0]
	d.rb.b = in
	d.rb.a = len(in)
	d.r = &d.rb
	_, d.js = h.(*JsonHandle)
	d.d = h.newDecDriver(d)
	// d.d = h.newDecDriver(decReaderT{true, &d.rb, &d.ri})
	return
}

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface{}
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
//
// When decoding into a nil interface{}, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface{} value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
func (d *Decoder) Decode(v interface{}) (err error) {
	defer panicToErr(&err)
	d.decode(v)
	return
}

// this is not a smart swallow, as it allocates objects and does unnecessary work.
func (d *Decoder) swallowViaHammer() {
	var blank interface{}
	d.decodeValue(reflect.ValueOf(&blank).Elem(), nil)
}

func (d *Decoder) swallow() {
	// smarter decode that just swallows the content
	dd := d.d
	switch {
	case dd.TryDecodeAsNil():
	case dd.IsContainerType(valueTypeMap):
		containerLen := dd.ReadMapStart()
		clenGtEqualZero := containerLen >= 0
		for j := 0; ; j++ {
			if clenGtEqualZero {
				if j >= containerLen {
					break
				}
			} else if dd.CheckBreak() {
				break
			}
			d.swallow()
			d.swallow()
		}
		dd.ReadEnd()
	case dd.IsContainerType(valueTypeArray):
		containerLenS := dd.ReadArrayStart()
		clenGtEqualZero := containerLenS >= 0
		for j := 0; ; j++ {
			if clenGtEqualZero {
				if j >= containerLenS {
					break
				}
			} else if dd.CheckBreak() {
				break
			}
			d.swallow()
		}
		dd.ReadEnd()
	case dd.IsContainerType(valueTypeBytes):
		dd.DecodeBytes(d.b[:], false, true)
	case dd.IsContainerType(valueTypeString):
		dd.DecodeBytes(d.b[:], true, true)
		// dd.DecodeStringAsBytes(d.b[:])
	default:
		// these are all primitives, which we can get from decodeNaked
		dd.DecodeNaked()
	}
}

// MustDecode is like Decode, but panics if unable to Decode.
// This provides insight to the code location that triggered the error.
func (d *Decoder) MustDecode(v interface{}) {
	d.decode(v)
}

func (d *Decoder) decode(iv interface{}) {
	// if ics, ok := iv.(Selfer); ok {
	// 	ics.CodecDecodeSelf(d)
	// 	return
	// }

	if d.d.TryDecodeAsNil() {
		switch v := iv.(type) {
		case nil:
		case *string:
			*v = ""
		case *bool:
			*v = false
		case *int:
			*v = 0
		case *int8:
			*v = 0
		case *int16:
			*v = 0
		case *int32:
			*v = 0
		case *int64:
			*v = 0
		case *uint:
			*v = 0
		case *uint8:
			*v = 0
		case *uint16:
			*v = 0
		case *uint32:
			*v = 0
		case *uint64:
			*v = 0
		case *float32:
			*v = 0
		case *float64:
			*v = 0
		case *[]uint8:
			*v = nil
		case reflect.Value:
			d.chkPtrValue(v)
			v = v.Elem()
			if v.IsValid() {
				v.Set(reflect.Zero(v.Type()))
			}
		default:
			rv := reflect.ValueOf(iv)
			d.chkPtrValue(rv)
			rv = rv.Elem()
			if rv.IsValid() {
				rv.Set(reflect.Zero(rv.Type()))
			}
		}
		return
	}

	switch v := iv.(type) {
	case nil:
		d.error(cannotDecodeIntoNilErr)
		return

	case Selfer:
		v.CodecDecodeSelf(d)

	case reflect.Value:
		d.chkPtrValue(v)
		d.decodeValueNotNil(v.Elem(), nil)

	case *string:

		*v = d.d.DecodeString()
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(d.d.DecodeInt(intBitsize))
	case *int8:
		*v = int8(d.d.DecodeInt(8))
	case *int16:
		*v = int16(d.d.DecodeInt(16))
	case *int32:
		*v = int32(d.d.DecodeInt(32))
	case *int64:
		*v = d.d.DecodeInt(64)
	case *uint:
		*v = uint(d.d.DecodeUint(uintBitsize))
	case *uint8:
		*v = uint8(d.d.DecodeUint(8))
	case *uint16:
		*v = uint16(d.d.DecodeUint(16))
	case *uint32:
		*v = uint32(d.d.DecodeUint(32))
	case *uint64:
		*v = d.d.DecodeUint(64)
	case *float32:
		*v = float32(d.d.DecodeFloat(true))
	case *float64:
		*v = d.d.DecodeFloat(false)
	case *[]uint8:
		*v = d.d.DecodeBytes(*v, false, false)

	case *interface{}:
		d.decodeValueNotNil(reflect.ValueOf(iv).Elem(), nil)

	default:
		if !fastpathDecodeTypeSwitch(iv, d) {
			d.decodeI(iv, true, false, false, false)
		}
	}
}

func (d *Decoder) preDecodeValue(rv reflect.Value, tryNil bool) (rv2 reflect.Value, proceed bool) {
	if tryNil && d.d.TryDecodeAsNil() {
		// No need to check if a ptr, recursively, to determine
		// whether to set value to nil.
		// Just always set value to its zero type.
		if rv.IsValid() { // rv.CanSet() // always settable, except it's invalid
			rv.Set(reflect.Zero(rv.Type()))
		}
		return
	}

	// If stream is not containing a nil value, then we can deref to the base
	// non-pointer value, and decode into that.
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return rv, true
}

func (d *Decoder) decodeI(iv interface{}, checkPtr, tryNil, checkFastpath, checkCodecSelfer bool) {
	rv := reflect.ValueOf(iv)
	if checkPtr {
		d.chkPtrValue(rv)
	}
	rv, proceed := d.preDecodeValue(rv, tryNil)
	if proceed {
		fn := d.getDecFn(rv.Type(), checkFastpath, checkCodecSelfer)
		fn.f(&fn.i, rv)
	}
}

func (d *Decoder) decodeValue(rv reflect.Value, fn *decFn) {
	if rv, proceed := d.preDecodeValue(rv, true); proceed {
		if fn == nil {
			fn = d.getDecFn(rv.Type(), true, true)
		}
		fn.f(&fn.i, rv)
	}
}

func (d *Decoder) decodeValueNotNil(rv reflect.Value, fn *decFn) {
	if rv, proceed := d.preDecodeValue(rv, false); proceed {
		if fn == nil {
			fn = d.getDecFn(rv.Type(), true, true)
		}
		fn.f(&fn.i, rv)
	}
}

func (d *Decoder) getDecFn(rt reflect.Type, checkFastpath, checkCodecSelfer bool) (fn *decFn) {
	rtid := reflect.ValueOf(rt).Pointer()

	// retrieve or register a focus'ed function for this type
	// to eliminate need to do the retrieval multiple times

	// if d.f == nil && d.s == nil { debugf("---->Creating new dec f map for type: %v\n", rt) }
	var ok bool
	if useMapForCodecCache {
		fn, ok = d.f[rtid]
	} else {
		for i := range d.s {
			v := &(d.s[i])
			if v.rtid == rtid {
				fn, ok = &(v.fn), true
				break
			}
		}
	}
	if ok {
		return
	}

	if useMapForCodecCache {
		if d.f == nil {
			d.f = make(map[uintptr]*decFn, initCollectionCap)
		}
		fn = new(decFn)
		d.f[rtid] = fn
	} else {
		if d.s == nil {
			d.s = make([]decRtidFn, 0, initCollectionCap)
		}
		d.s = append(d.s, decRtidFn{rtid: rtid})
		fn = &(d.s[len(d.s)-1]).fn
	}

	// debugf("\tCreating new dec fn for type: %v\n", rt)
	ti := d.h.getTypeInfo(rtid, rt)
	fi := &(fn.i)
	fi.d = d
	fi.ti = ti

	// An extension can be registered for any type, regardless of the Kind
	// (e.g. type BitSet int64, type MyStruct { / * unexported fields * / }, type X []int, etc.
	//
	// We can't check if it's an extension byte here first, because the user may have
	// registered a pointer or non-pointer type, meaning we may have to recurse first
	// before matching a mapped type, even though the extension byte is already detected.
	//
	// NOTE: if decoding into a nil interface{}, we return a non-nil
	// value except even if the container registers a length of 0.
	if checkCodecSelfer && ti.cs {
		fn.f = (*decFnInfo).selferUnmarshal
	} else if rtid == rawExtTypId {
		fn.f = (*decFnInfo).rawExt
	} else if d.d.IsBuiltinType(rtid) {
		fn.f = (*decFnInfo).builtin
	} else if xfFn := d.h.getExt(rtid); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.f = (*decFnInfo).ext
	} else if supportMarshalInterfaces && d.be && ti.bunm {
		fn.f = (*decFnInfo).binaryUnmarshal
	} else if supportMarshalInterfaces && !d.be && d.js && ti.junm {
		//If JSON, we should check JSONUnmarshal before textUnmarshal
		fn.f = (*decFnInfo).jsonUnmarshal
	} else if supportMarshalInterfaces && !d.be && ti.tunm {
		fn.f = (*decFnInfo).textUnmarshal
	} else {
		rk := rt.Kind()
		if fastpathEnabled && checkFastpath && (rk == reflect.Map || rk == reflect.Slice) {
			if rt.PkgPath() == "" {
				if idx := fastpathAV.index(rtid); idx != -1 {
					fn.f = fastpathAV[idx].decfn
				}
			} else {
				// use mapping for underlying type if there
				ok = false
				var rtu reflect.Type
				if rk == reflect.Map {
					rtu = reflect.MapOf(rt.Key(), rt.Elem())
				} else {
					rtu = reflect.SliceOf(rt.Elem())
				}
				rtuid := reflect.ValueOf(rtu).Pointer()
				if idx := fastpathAV.index(rtuid); idx != -1 {
					xfnf := fastpathAV[idx].decfn
					xrt := fastpathAV[idx].rt
					fn.f = func(xf *decFnInfo, xrv reflect.Value) {
						// xfnf(xf, xrv.Convert(xrt))
						xfnf(xf, xrv.Addr().Convert(reflect.PtrTo(xrt)).Elem())
					}
				}
			}
		}
		if fn.f == nil {
			switch rk {
			case reflect.String:
				fn.f = (*decFnInfo).kString
			case reflect.Bool:
				fn.f = (*decFnInfo).kBool
			case reflect.Int:
				fn.f = (*decFnInfo).kInt
			case reflect.Int64:
				fn.f = (*decFnInfo).kInt64
			case reflect.Int32:
				fn.f = (*decFnInfo).kInt32
			case reflect.Int8:
				fn.f = (*decFnInfo).kInt8
			case reflect.Int16:
				fn.f = (*decFnInfo).kInt16
			case reflect.Float32:
				fn.f = (*decFnInfo).kFloat32
			case reflect.Float64:
				fn.f = (*decFnInfo).kFloat64
			case reflect.Uint8:
				fn.f = (*decFnInfo).kUint8
			case reflect.Uint64:
				fn.f = (*decFnInfo).kUint64
			case reflect.Uint:
				fn.f = (*decFnInfo).kUint
			case reflect.Uint32:
				fn.f = (*decFnInfo).kUint32
			case reflect.Uint16:
				fn.f = (*decFnInfo).kUint16
				// case reflect.Ptr:
				// 	fn.f = (*decFnInfo).kPtr
			case reflect.Uintptr:
				fn.f = (*decFnInfo).kUintptr
			case reflect.Interface:
				fn.f = (*decFnInfo).kInterface
			case reflect.Struct:
				fn.f = (*decFnInfo).kStruct
			case reflect.Chan:
				fi.seq = seqTypeChan
				fn.f = (*decFnInfo).kSlice
			case reflect.Slice:
				fi.seq = seqTypeSlice
				fn.f = (*decFnInfo).kSlice
			case reflect.Array:
				fi.seq = seqTypeArray
				fn.f = (*decFnInfo).kArray
			case reflect.Map:
				fn.f = (*decFnInfo).kMap
			default:
				fn.f = (*decFnInfo).kErr
			}
		}
	}

	return
}

func (d *Decoder) structFieldNotFound(index int, rvkencname string) {
	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorf("no matching struct field found when decoding stream array at index %v", index)
			return
		} else if rvkencname != "" {
			d.errorf("no matching struct field found when decoding stream map with key %s", rvkencname)
			return
		}
	}
	d.swallow()
}

func (d *Decoder) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", sliceLen, streamLen)
	}
}

func (d *Decoder) chkPtrValue(rv reflect.Value) {
	// We can only decode into a non-nil pointer
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return
	}
	if !rv.IsValid() {
		d.error(cannotDecodeIntoNilErr)
		return
	}
	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
		return
	}
	rvi := rv.Interface()
	d.errorf("cannot decode into non-pointer or nil pointer. Got: %v, %T, %v", rv.Kind(), rvi, rvi)
}

func (d *Decoder) error(err error) {
	panic(err)
}

func (d *Decoder) errorf(format string, params ...interface{}) {
	params2 := make([]interface{}, len(params)+1)
	params2[0] = d.r.numread()
	copy(params2[1:], params)
	err := fmt.Errorf("[pos %d]: "+format, params2...)
	panic(err)
}

// --------------------------------------------------

// decSliceHelper assists when decoding into a slice, from a map or an array in the stream.
// A slice can be set from a map or array in stream. This supports the MapBySlice interface.
type decSliceHelper struct {
	dd decDriver
	ct valueType
}

func (d *Decoder) decSliceHelperStart() (x decSliceHelper, clen int) {
	x.dd = d.d
	if x.dd.IsContainerType(valueTypeArray) {
		x.ct = valueTypeArray
		clen = x.dd.ReadArrayStart()
	} else if x.dd.IsContainerType(valueTypeMap) {
		x.ct = valueTypeMap
		clen = x.dd.ReadMapStart() * 2
	} else {
		d.errorf("only encoded map or array can be decoded into a slice")
	}
	return
}

func (x decSliceHelper) End() {
	x.dd.ReadEnd()
}

func decByteSlice(r decReader, clen int, bs []byte) (bsOut []byte) {
	if clen == 0 {
		return zeroByteSlice
	}
	if len(bs) == clen {
		bsOut = bs
	} else if cap(bs) >= clen {
		bsOut = bs[:clen]
	} else {
		bsOut = make([]byte, clen)
	}
	r.readb(bsOut)
	return
}

func detachZeroCopyBytes(isBytesReader bool, dest []byte, in []byte) (out []byte) {
	if xlen := len(in); xlen > 0 {
		if isBytesReader || xlen <= scratchByteArrayLen {
			if cap(dest) >= xlen {
				out = dest[:xlen]
			} else {
				out = make([]byte, xlen)
			}
			copy(out, in)
			return
		}
	}
	return in
}

// decInferLen will infer a sensible length, given the following:
//    - clen: length wanted.
//    - maxlen: max length to be returned.
//      if <= 0, it is unset, and we infer it based on the unit size
//    - unit: number of bytes for each element of the collection
func decInferLen(clen, maxlen, unit int) (rvlen int, truncated bool) {
	// handle when maxlen is not set i.e. <= 0
	if clen <= 0 {
		return
	}
	if maxlen <= 0 {
		// no maxlen defined. Use maximum of 256K memory, with a floor of 4K items.
		// maxlen = 256 * 1024 / unit
		// if maxlen < (4 * 1024) {
		// 	maxlen = 4 * 1024
		// }
		if unit < (256 / 4) {
			maxlen = 256 * 1024 / unit
		} else {
			maxlen = 4 * 1024
		}
	}
	if clen > maxlen {
		rvlen = maxlen
		truncated = true
	} else {
		rvlen = clen
	}
	return
	// if clen <= 0 {
	// 	rvlen = 0
	// } else if maxlen > 0 && clen > maxlen {
	// 	rvlen = maxlen
	// 	truncated = true
	// } else {
	// 	rvlen = clen
	// }
	// return
}

// // implement overall decReader wrapping both, for possible use inline:
// type decReaderT struct {
// 	bytes bool
// 	rb    *bytesDecReader
// 	ri    *ioDecReader
// }
//
// // implement *Decoder as a decReader.
// // Using decReaderT (defined just above) caused performance degradation
// // possibly because of constant copying the value,
// // and some value->interface conversion causing allocation.
// func (d *Decoder) unreadn1() {
// 	if d.bytes {
// 		d.rb.unreadn1()
// 	} else {
// 		d.ri.unreadn1()
// 	}
// }
// ... for other methods of decReader.
// Testing showed that performance improvement was negligible.
