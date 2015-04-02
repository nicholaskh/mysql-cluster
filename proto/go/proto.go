//from gopbuf,  https://github.com/akunspy/gopbuf
package proto
import ( 
	"fmt"
	"unsafe"
)

const (
	SIZE_32BIT       = 4
	SIZE_64BIT       = 8
	WIRE_TYPE_VARINT = 0
	WIRE_TYPE_64BIT  = 1
	WIRE_TYPE_LENGTH = 2
	WIRE_TYPE_32_BIT = 5
)

func Encode32(buf []byte, index int, number uint32) int {
	if number < 0x80 {
		buf[index+0] = byte(number)
		return 1
	}
	buf[index+0] = byte(number | 0x80)
	if number < 0x4000 {
		buf[index+1] = byte(number >> 7)
		return 2
	}
	buf[index+1] = byte((number >> 7) | 0x80)
	if number < 0x200000 {
		buf[index+2] = byte(number >> 14)
		return 3
	}
	buf[index+2] = byte((number >> 14) | 0x80)
	if number < 0x10000000 {
		buf[index+3] = byte(number >> 21)
		return 4
	}
	buf[index+3] = byte((number >> 21) | 0x80)
	buf[index+4] = byte(number >> 28)
	return 5
}

func Encode32Size(number uint32) int {
	switch {
	case number < 0x80:
		return 1
	case number < 0x4000:
		return 2
	case number < 0x200000:
		return 3
	case number < 0x10000000:
		return 4
	default:
		return 5
	}
}

func Encode64(buf []byte, index int, number uint64) int {
	if (number & 0xffffffff) == number {
		return Encode32(buf, index, uint32(number))
	}

	i := 0
	for ; number >= 0x80; i++ {
		buf[index+i] = byte(number | 0x80)
		number >>= 7
	}
	buf[index+i] = byte(number)
	return i + 1
}

func Encode64Size(number uint64) int {
	if (number & 0xffffffff) == number {
		return Encode32Size(uint32(number))
	}

	i := 0
	for ; number >= 0x80; i++ {
		number >>= 7
	}
	return i + 1
}

func Decode(buf []byte, index int) (x uint64, n int) {
	buf = buf[index:]

	for shift := uint(0); ; shift += 7 {
		if n >= len(buf) {
			return 0, n
		}

		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			break
		}
	}
	return x, n
}

func Zigzag32(n int) int {
	return (n << 1) ^ (n >> 31)
}

func Zigzag64(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}

func Dezigzag32(n int) int {
	return (n >> 1) ^ -(n & 1)
}

func Dezigzag64(n int64) int64 {
	return (int64(n) >> 1) ^ -(n & 1)
}

func StringToBytes(s string) []byte {
	l := len(s)
	ret := make([]byte, l)
	copy(ret, s)
	return ret
}

func BooleanToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func WriteInt32Size(x int) int {
	if x < 0 {
		return Encode64Size(uint64(x))
	} else {
		return Encode32Size(uint32(x))
	}
}

func WriteInt64Size(x int64) int {
	return Encode64Size(uint64(x))
}

func WriteUInt32Size(x uint32) int {
	return Encode32Size(x)
}

func WriteUInt64Size(x uint64) int {
	return Encode64Size(x)
}

func WriteSInt32Size(x int) int {
	v := Zigzag32(x)
	return Encode32Size(uint32(v))
}

func WriteSInt64Size(x int64) int {
	v := Zigzag64(x)
	return Encode64Size(uint64(v))
}

func WriteStringSize(s string) int {
	return len(s)
}

//ProtoBuffer
type ProtoBuffer struct {
	buf []byte
	pos int
}

func NewProtoBuffer(buf []byte) *ProtoBuffer {
	buffer := new(ProtoBuffer)
	buffer.buf = buf

	return buffer
}

func (p *ProtoBuffer) AddPos(v int) {
	if v <= 0 {
		return
	}

	p.pos += v
	if p.pos > len(p.buf) {
		p.pos = len(p.buf)
	}
}

func (p *ProtoBuffer) ResetPos() {
	p.pos = 0
}

func (p *ProtoBuffer) PrintBuffer() {
	fmt.Println("Size:", p.pos)
	for i := 0; i < p.pos; i++ {
		fmt.Printf("%2.2X ", p.buf[i])
	}
	fmt.Println()
}

//int
func (p *ProtoBuffer) WriteInt32(x int) {
	if x < 0 {
		p.AddPos(Encode64(p.buf, p.pos, uint64(x)))
	} else {
		p.AddPos(Encode32(p.buf, p.pos, uint32(x)))
	}
}

func (p *ProtoBuffer) ReadInt32() int {
	x, n := Decode(p.buf, p.pos)
	p.AddPos(n)
	return int(x)
}

//uint32
func (p *ProtoBuffer) WriteUInt32(x uint32) {
	p.AddPos(Encode32(p.buf, p.pos, uint32(x)))
}

func (p *ProtoBuffer) ReadUInt32() uint32 {
	x, n := Decode(p.buf, p.pos)
	p.AddPos(n)
	return uint32(x)
}

//int64
func (p *ProtoBuffer) WriteInt64(x int64) {
	p.AddPos(Encode64(p.buf, p.pos, uint64(x)))
}

func (p *ProtoBuffer) ReadInt64() int64 {
	x, n := Decode(p.buf, p.pos)
	p.AddPos(n)
	return int64(x)
}

//uint64
func (p *ProtoBuffer) WriteUInt64(x uint64) {
	p.AddPos(Encode64(p.buf, p.pos, uint64(x)))
}

func (p *ProtoBuffer) ReadUInt64() uint64 {
	x, n := Decode(p.buf, p.pos)
	p.AddPos(n)
	return x
}

//sint32
func (p *ProtoBuffer) WriteSInt32(x int) {
	v := Zigzag32(x)
	p.WriteInt32(v)
}

func (p *ProtoBuffer) ReadSInt32() int {
	v := p.ReadInt32()
	return Dezigzag32(v)
}

//sint64
func (p *ProtoBuffer) WriteSInt64(x int64) {
	v := Zigzag64(x)
	p.WriteInt64(v)
}

func (p *ProtoBuffer) ReadSInt64() int64 {
	v := p.ReadInt64()
	return Dezigzag64(v)
}

//sfixed32
func (p *ProtoBuffer) WriteSFixed32(x int) {
	p.WriteFixed32(uint32(x))
}

func (p *ProtoBuffer) ReadSFixed32() int {
	return int(p.ReadFixed32())
}

//sfixed64
func (p *ProtoBuffer) WriteSFixed64(x int64) {
	p.WriteFixed64(uint64(x))
}

func (p *ProtoBuffer) ReadSFixed64() int64 {
	return int64(p.ReadFixed64())
}

//fixed32
func (p *ProtoBuffer) WriteFixed32(x uint32) {
	buf := p.buf
	pos := p.pos

	buf[pos] = uint8(x)
	buf[pos+1] = uint8(x >> 8)
	buf[pos+2] = uint8(x >> 16)
	buf[pos+3] = uint8(x >> 24)
	p.AddPos(SIZE_32BIT)
}

func (p *ProtoBuffer) ReadFixed32() uint32 {
	buf := p.buf
	pos := p.pos
	p.AddPos(SIZE_32BIT)

	if p.pos >= len(buf) {
		return 0
	}

	x := uint32(buf[pos]) |
		(uint32(buf[pos+1]) << 8) |
		(uint32(buf[pos+2]) << 16) |
		(uint32(buf[pos+3]) << 24)
	return x
}

//fixed64
func (p *ProtoBuffer) WriteFixed64(x uint64) {
	buf := p.buf
	pos := p.pos
	p.AddPos(SIZE_64BIT)

	buf[pos] = uint8(x)
	buf[pos+1] = uint8(x >> 8)
	buf[pos+2] = uint8(x >> 16)
	buf[pos+3] = uint8(x >> 24)
	buf[pos+4] = uint8(x >> 32)
	buf[pos+5] = uint8(x >> 40)
	buf[pos+6] = uint8(x >> 48)
	buf[pos+7] = uint8(x >> 56)
}

func (p *ProtoBuffer) ReadFixed64() uint64 {
	buf := p.buf
	pos := p.pos
	p.AddPos(SIZE_64BIT)

	if p.pos >= len(buf) {
		return 0
	}

	ret_low := uint32(buf[pos]) |
		(uint32(buf[pos+1]) << 8) |
		(uint32(buf[pos+2]) << 16) |
		(uint32(buf[pos+3]) << 24)
	ret_high := uint32(buf[pos+4]) |
		(uint32(buf[pos+5]) << 8) |
		(uint32(buf[pos+6]) << 16) |
		(uint32(buf[pos+7]) << 24)
	return (uint64(ret_high) << 32) | uint64(ret_low)
}

//float32
func (p *ProtoBuffer) WriteFloat32(f float32) {
	p.WriteFixed32(*(*uint32)(unsafe.Pointer(&f)))
}

func (p *ProtoBuffer) ReadFloat32() float32 {
	x := p.ReadFixed32()
	return *(*float32)(unsafe.Pointer(&x))
}

//float64
func (p *ProtoBuffer) WriteFloat64(f float64) {
	p.WriteFixed64(*(*uint64)(unsafe.Pointer(&f)))
}

func (p *ProtoBuffer) ReadFloat64() float64 {
	x := p.ReadFixed64()
	return *(*float64)(unsafe.Pointer(&x))
}

//boolean
func (p *ProtoBuffer) WriteBoolean(b bool) {
	if b {
		p.WriteInt32(1)
	} else {
		p.WriteInt32(0)
	}
}

func (p *ProtoBuffer) ReadBoolean() bool {
	x := p.ReadInt32()
	return x == 1
}

//string
func (p *ProtoBuffer) WriteString(s string) {
	l := len(s)
	p.WriteUInt32(uint32(l))
	copy(p.buf[p.pos:], s)
	p.AddPos(l)
}

func (p *ProtoBuffer) ReadString() string {
	l := p.ReadUInt32()
	old_pos := p.pos
	p.AddPos(int(l))
	s := string(p.buf[old_pos:p.pos])
	return s
}

//[]byte
func (p *ProtoBuffer) WriteBytes(b []byte) {
	l := len(b)
	p.WriteUInt32(uint32(l))
	copy(p.buf[p.pos:], b)
	p.AddPos(l)
}

func (p *ProtoBuffer) ReadBytes() []byte {
	l := p.ReadUInt32()
	old_pos := p.pos
	p.AddPos(int(l))

	b := make([]byte, l)
	copy(b, p.buf[old_pos:p.pos])
	return b
}

func (p *ProtoBuffer) GetUnknowFieldValueSize(wire_tag int) {
	wire_type := wire_tag & 0x7

	switch wire_type {
	case WIRE_TYPE_VARINT:
		p.ReadUInt32()
	case WIRE_TYPE_64BIT:
		p.ReadFixed64()
	case WIRE_TYPE_LENGTH:
		size := p.ReadInt32()
		if size > 0 {
			p.AddPos(int(size))
		}
	case WIRE_TYPE_32_BIT:
		p.ReadFixed32()
	}
}

//error
type ProtoError struct {
	What string
}

func (e *ProtoError) Error() string {
	return e.What
}

//proto interface
type Message interface {
	Serialize(buf []byte) (size int, err error)
	Parse(buf []byte, msg_size int) error
	Clear()
	ByteSize() int
	IsInitialized() bool
}
 
type QueryStruct struct {
    cached_byte_size int
    has_flag_0 uint32
    pool string
    sql string
    args []string
}

func NewQueryStruct() *QueryStruct {
    p := new(QueryStruct)
    p.pool = ""
    p.sql = ""
    p.args = make([]string, 0)
    return p
}
//pool
func (p *QueryStruct) Getpool() string {
    return p.pool
}
func (p *QueryStruct) Setpool(v string) {
    p.pool = v
    if v == "" {
        p.has_flag_0 &= 0xfffffffe
    } else {
        p.has_flag_0 |= 0x1
    }
    p.cached_byte_size = 0
}
func (p *QueryStruct) Haspool() bool {
    return (p.has_flag_0 & 0x1) != 0
}
func (p *QueryStruct) Clearpool() {
    p.pool = ""
    p.has_flag_0 &= 0xfffffffe
    p.cached_byte_size = 0
}

//sql
func (p *QueryStruct) Getsql() string {
    return p.sql
}
func (p *QueryStruct) Setsql(v string) {
    p.sql = v
    if v == "" {
        p.has_flag_0 &= 0xfffffffd
    } else {
        p.has_flag_0 |= 0x2
    }
    p.cached_byte_size = 0
}
func (p *QueryStruct) Hassql() bool {
    return (p.has_flag_0 & 0x2) != 0
}
func (p *QueryStruct) Clearsql() {
    p.sql = ""
    p.has_flag_0 &= 0xfffffffd
    p.cached_byte_size = 0
}

//args
func (p *QueryStruct) Sizeargs() int { return len(p.args) }
func (p *QueryStruct) Getargs(index int) string { return p.args[index] }
func (p *QueryStruct) Addargs(v string) {
    p.args = append(p.args, v)
    p.cached_byte_size = 0
}
func (p *QueryStruct) Clearargs() { 
    p.args = make([]string, 0)
    p.cached_byte_size = 0
}

func (p *QueryStruct) Serialize(buf []byte) (size int, err error) {
    b := NewProtoBuffer(buf)
    err = p.do_serialize(b)
    return b.pos, err
}
func (p *QueryStruct) do_serialize(b *ProtoBuffer) error {
    buf_size := len(b.buf) - b.pos
    byte_size := p.ByteSize()
    if byte_size > buf_size {
        return &ProtoError{"Serialize error, byte_size > buf_size"}
    }

    list_count := 0
    if p.Haspool() {
        b.WriteInt32(10)
        b.WriteString(p.pool)
    }
    if p.Hassql() {
        b.WriteInt32(18)
        b.WriteString(p.sql)
    }
    list_count = p.Sizeargs()
    for i := 0; i < list_count; i++ {
        b.WriteInt32(26)
        b.WriteString(p.args[i])
    }
    return nil
}

func (p *QueryStruct) Parse(buf []byte, msg_size int) error {
    b := NewProtoBuffer(buf)
    return p.do_parse(b, msg_size)
}
func (p *QueryStruct) do_parse(b *ProtoBuffer, msg_size int) error {
    msg_end := b.pos + msg_size
    if msg_end > len(b.buf) {
        return &ProtoError{"Parse QueryStruct error, msg size out of buffer"}
    }
    p.Clear()
    for b.pos < msg_end {
        wire_tag := b.ReadInt32()
        switch wire_tag {
        case 10:
            p.Setpool(b.ReadString())
        case 18:
            p.Setsql(b.ReadString())
        case 26:
            p.Addargs(b.ReadString())
        default: b.GetUnknowFieldValueSize(wire_tag)
        }
    }

    if !p.IsInitialized() {
        return &ProtoError{"proto.QueryStruct parse error, miss required field"} 
    }
    return nil
}

func (p *QueryStruct) Clear() {
    p.pool = ""
    p.sql = ""
    p.args = make([]string, 0)
    p.cached_byte_size = 0
    p.has_flag_0 = 0
}

func (p *QueryStruct) ByteSize() int {
    if p.cached_byte_size != 0 {
        return p.cached_byte_size
    }
    list_count := 0
    if p.Haspool() {
        p.cached_byte_size += WriteInt32Size(10)
        size := WriteStringSize(p.pool)
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    if p.Hassql() {
        p.cached_byte_size += WriteInt32Size(18)
        size := WriteStringSize(p.sql)
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    list_count = p.Sizeargs()
    for i := 0; i < list_count; i++ {
        p.cached_byte_size += WriteInt32Size(26)
        size := WriteStringSize(p.args[i])
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    return p.cached_byte_size
}

func (p *QueryStruct) IsInitialized() bool {
    if (p.has_flag_0 & 0x3) != 0x3 {
        return false
    }
    return true
}

type Rows struct {
    cached_byte_size int
    rows []*Rows_Row
}

func NewRows() *Rows {
    p := new(Rows)
    p.rows = make([]*Rows_Row, 0)
    return p
}
//rows
func (p *Rows) Sizerows() int { return len(p.rows) }
func (p *Rows) Getrows(index int) *Rows_Row { return p.rows[index] }
func (p *Rows) Addrows(v *Rows_Row) {
    p.rows = append(p.rows, v)
    p.cached_byte_size = 0
}
func (p *Rows) Clearrows() { 
    p.rows = make([]*Rows_Row, 0)
    p.cached_byte_size = 0
}

func (p *Rows) Serialize(buf []byte) (size int, err error) {
    b := NewProtoBuffer(buf)
    err = p.do_serialize(b)
    return b.pos, err
}
func (p *Rows) do_serialize(b *ProtoBuffer) error {
    buf_size := len(b.buf) - b.pos
    byte_size := p.ByteSize()
    if byte_size > buf_size {
        return &ProtoError{"Serialize error, byte_size > buf_size"}
    }

    list_count := 0
    list_count = p.Sizerows()
    for i := 0; i < list_count; i++ {
        b.WriteInt32(10)
        size := p.rows[i].ByteSize()
        b.WriteInt32(size)
        p.rows[i].do_serialize(b)
    }
    return nil
}

func (p *Rows) Parse(buf []byte, msg_size int) error {
    b := NewProtoBuffer(buf)
    return p.do_parse(b, msg_size)
}
func (p *Rows) do_parse(b *ProtoBuffer, msg_size int) error {
    msg_end := b.pos + msg_size
    if msg_end > len(b.buf) {
        return &ProtoError{"Parse Rows error, msg size out of buffer"}
    }
    p.Clear()
    for b.pos < msg_end {
        wire_tag := b.ReadInt32()
        switch wire_tag {
        case 10:
            rows_size := b.ReadInt32()
            if rows_size > msg_end-b.pos {
                return &ProtoError{"parse Rows_Row error"}
            }
            rows_tmp := NewRows_Row()
            p.Addrows(rows_tmp)
            if rows_size > 0 {
                e := rows_tmp.do_parse(b, int(rows_size))
                if e != nil {
                    return e
                }
            } else {
                return &ProtoError{"parse Rows_Row error"}
            }
        default: b.GetUnknowFieldValueSize(wire_tag)
        }
    }

    if !p.IsInitialized() {
        return &ProtoError{"proto.Rows parse error, miss required field"} 
    }
    return nil
}

func (p *Rows) Clear() {
    p.rows = make([]*Rows_Row, 0)
    p.cached_byte_size = 0
}

func (p *Rows) ByteSize() int {
    if p.cached_byte_size != 0 {
        return p.cached_byte_size
    }
    list_count := 0
    list_count = p.Sizerows()
    for i := 0; i < list_count; i++ {
        p.cached_byte_size += WriteInt32Size(10)
        size := p.rows[i].ByteSize()
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    return p.cached_byte_size
}

func (p *Rows) IsInitialized() bool {
    return true
}

type Rows_Row struct {
    cached_byte_size int
    has_flag_0 uint32
    column string
    value string
}

func NewRows_Row() *Rows_Row {
    p := new(Rows_Row)
    p.column = ""
    p.value = ""
    return p
}
//column
func (p *Rows_Row) Getcolumn() string {
    return p.column
}
func (p *Rows_Row) Setcolumn(v string) {
    p.column = v
    if v == "" {
        p.has_flag_0 &= 0xfffffffe
    } else {
        p.has_flag_0 |= 0x1
    }
    p.cached_byte_size = 0
}
func (p *Rows_Row) Hascolumn() bool {
    return (p.has_flag_0 & 0x1) != 0
}
func (p *Rows_Row) Clearcolumn() {
    p.column = ""
    p.has_flag_0 &= 0xfffffffe
    p.cached_byte_size = 0
}

//value
func (p *Rows_Row) Getvalue() string {
    return p.value
}
func (p *Rows_Row) Setvalue(v string) {
    p.value = v
    if v == "" {
        p.has_flag_0 &= 0xfffffffd
    } else {
        p.has_flag_0 |= 0x2
    }
    p.cached_byte_size = 0
}
func (p *Rows_Row) Hasvalue() bool {
    return (p.has_flag_0 & 0x2) != 0
}
func (p *Rows_Row) Clearvalue() {
    p.value = ""
    p.has_flag_0 &= 0xfffffffd
    p.cached_byte_size = 0
}

func (p *Rows_Row) Serialize(buf []byte) (size int, err error) {
    b := NewProtoBuffer(buf)
    err = p.do_serialize(b)
    return b.pos, err
}
func (p *Rows_Row) do_serialize(b *ProtoBuffer) error {
    buf_size := len(b.buf) - b.pos
    byte_size := p.ByteSize()
    if byte_size > buf_size {
        return &ProtoError{"Serialize error, byte_size > buf_size"}
    }

    if p.Hascolumn() {
        b.WriteInt32(10)
        b.WriteString(p.column)
    }
    if p.Hasvalue() {
        b.WriteInt32(18)
        b.WriteString(p.value)
    }
    return nil
}

func (p *Rows_Row) Parse(buf []byte, msg_size int) error {
    b := NewProtoBuffer(buf)
    return p.do_parse(b, msg_size)
}
func (p *Rows_Row) do_parse(b *ProtoBuffer, msg_size int) error {
    msg_end := b.pos + msg_size
    if msg_end > len(b.buf) {
        return &ProtoError{"Parse Rows_Row error, msg size out of buffer"}
    }
    p.Clear()
    for b.pos < msg_end {
        wire_tag := b.ReadInt32()
        switch wire_tag {
        case 10:
            p.Setcolumn(b.ReadString())
        case 18:
            p.Setvalue(b.ReadString())
        default: b.GetUnknowFieldValueSize(wire_tag)
        }
    }

    if !p.IsInitialized() {
        return &ProtoError{"proto.Rows.Row parse error, miss required field"} 
    }
    return nil
}

func (p *Rows_Row) Clear() {
    p.column = ""
    p.value = ""
    p.cached_byte_size = 0
    p.has_flag_0 = 0
}

func (p *Rows_Row) ByteSize() int {
    if p.cached_byte_size != 0 {
        return p.cached_byte_size
    }
    if p.Hascolumn() {
        p.cached_byte_size += WriteInt32Size(10)
        size := WriteStringSize(p.column)
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    if p.Hasvalue() {
        p.cached_byte_size += WriteInt32Size(18)
        size := WriteStringSize(p.value)
        p.cached_byte_size += WriteInt32Size(size)
        p.cached_byte_size += size
    }
    return p.cached_byte_size
}

func (p *Rows_Row) IsInitialized() bool {
    if (p.has_flag_0 & 0x3) != 0x3 {
        return false
    }
    return true
}

