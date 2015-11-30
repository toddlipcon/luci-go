// Code generated by protoc-gen-go.
// source: metrics.proto
// DO NOT EDIT!

package ts_mon_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MetricsField_FieldType int32

const (
	MetricsField_STRING MetricsField_FieldType = 1
	MetricsField_INT    MetricsField_FieldType = 2
	MetricsField_BOOL   MetricsField_FieldType = 3
)

var MetricsField_FieldType_name = map[int32]string{
	1: "STRING",
	2: "INT",
	3: "BOOL",
}
var MetricsField_FieldType_value = map[string]int32{
	"STRING": 1,
	"INT":    2,
	"BOOL":   3,
}

func (x MetricsField_FieldType) Enum() *MetricsField_FieldType {
	p := new(MetricsField_FieldType)
	*p = x
	return p
}
func (x MetricsField_FieldType) String() string {
	return proto.EnumName(MetricsField_FieldType_name, int32(x))
}
func (x *MetricsField_FieldType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MetricsField_FieldType_value, data, "MetricsField_FieldType")
	if err != nil {
		return err
	}
	*x = MetricsField_FieldType(value)
	return nil
}

type PrecomputedDistribution_SpecType int32

const (
	PrecomputedDistribution_CANONICAL_POWERS_OF_2        PrecomputedDistribution_SpecType = 1
	PrecomputedDistribution_CANONICAL_POWERS_OF_10_P_0_2 PrecomputedDistribution_SpecType = 2
	PrecomputedDistribution_CANONICAL_POWERS_OF_10       PrecomputedDistribution_SpecType = 3
	PrecomputedDistribution_CUSTOM_PARAMETERIZED         PrecomputedDistribution_SpecType = 20
	PrecomputedDistribution_CUSTOM_BOUNDED               PrecomputedDistribution_SpecType = 21
)

var PrecomputedDistribution_SpecType_name = map[int32]string{
	1:  "CANONICAL_POWERS_OF_2",
	2:  "CANONICAL_POWERS_OF_10_P_0_2",
	3:  "CANONICAL_POWERS_OF_10",
	20: "CUSTOM_PARAMETERIZED",
	21: "CUSTOM_BOUNDED",
}
var PrecomputedDistribution_SpecType_value = map[string]int32{
	"CANONICAL_POWERS_OF_2":        1,
	"CANONICAL_POWERS_OF_10_P_0_2": 2,
	"CANONICAL_POWERS_OF_10":       3,
	"CUSTOM_PARAMETERIZED":         20,
	"CUSTOM_BOUNDED":               21,
}

func (x PrecomputedDistribution_SpecType) Enum() *PrecomputedDistribution_SpecType {
	p := new(PrecomputedDistribution_SpecType)
	*p = x
	return p
}
func (x PrecomputedDistribution_SpecType) String() string {
	return proto.EnumName(PrecomputedDistribution_SpecType_name, int32(x))
}
func (x *PrecomputedDistribution_SpecType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(PrecomputedDistribution_SpecType_value, data, "PrecomputedDistribution_SpecType")
	if err != nil {
		return err
	}
	*x = PrecomputedDistribution_SpecType(value)
	return nil
}

type MetricsData_Units int32

const (
	MetricsData_UNKNOWN_UNITS MetricsData_Units = 0
	MetricsData_SECONDS       MetricsData_Units = 1
	MetricsData_MILLISECONDS  MetricsData_Units = 2
	MetricsData_MICROSECONDS  MetricsData_Units = 3
	MetricsData_NANOSECONDS   MetricsData_Units = 4
	MetricsData_BITS          MetricsData_Units = 21
	MetricsData_BYTES         MetricsData_Units = 22
	// * 1000 bytes (not 1024).
	MetricsData_KILOBYTES MetricsData_Units = 31
	// * 1e6 (1,000,000) bytes.
	MetricsData_MEGABYTES MetricsData_Units = 32
	// * 1e9 (1,000,000,000) bytes.
	MetricsData_GIGABYTES MetricsData_Units = 33
	// * 1024 bytes.
	MetricsData_KIBIBYTES MetricsData_Units = 41
	// * 1024^2 (1,048,576) bytes.
	MetricsData_MEBIBYTES MetricsData_Units = 42
	// * 1024^3 (1,073,741,824) bytes.
	MetricsData_GIBIBYTES MetricsData_Units = 43
	// * Extended Units
	MetricsData_AMPS            MetricsData_Units = 60
	MetricsData_MILLIAMPS       MetricsData_Units = 61
	MetricsData_DEGREES_CELSIUS MetricsData_Units = 62
)

var MetricsData_Units_name = map[int32]string{
	0:  "UNKNOWN_UNITS",
	1:  "SECONDS",
	2:  "MILLISECONDS",
	3:  "MICROSECONDS",
	4:  "NANOSECONDS",
	21: "BITS",
	22: "BYTES",
	31: "KILOBYTES",
	32: "MEGABYTES",
	33: "GIGABYTES",
	41: "KIBIBYTES",
	42: "MEBIBYTES",
	43: "GIBIBYTES",
	60: "AMPS",
	61: "MILLIAMPS",
	62: "DEGREES_CELSIUS",
}
var MetricsData_Units_value = map[string]int32{
	"UNKNOWN_UNITS":   0,
	"SECONDS":         1,
	"MILLISECONDS":    2,
	"MICROSECONDS":    3,
	"NANOSECONDS":     4,
	"BITS":            21,
	"BYTES":           22,
	"KILOBYTES":       31,
	"MEGABYTES":       32,
	"GIGABYTES":       33,
	"KIBIBYTES":       41,
	"MEBIBYTES":       42,
	"GIBIBYTES":       43,
	"AMPS":            60,
	"MILLIAMPS":       61,
	"DEGREES_CELSIUS": 62,
}

func (x MetricsData_Units) Enum() *MetricsData_Units {
	p := new(MetricsData_Units)
	*p = x
	return p
}
func (x MetricsData_Units) String() string {
	return proto.EnumName(MetricsData_Units_name, int32(x))
}
func (x *MetricsData_Units) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MetricsData_Units_value, data, "MetricsData_Units")
	if err != nil {
		return err
	}
	*x = MetricsData_Units(value)
	return nil
}

type MetricsCollection struct {
	Data              []*MetricsData `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
	StartTimestampUs  *uint64        `protobuf:"varint,2,opt,name=start_timestamp_us" json:"start_timestamp_us,omitempty"`
	CollectionPointId *string        `protobuf:"bytes,3,opt,name=collection_point_id" json:"collection_point_id,omitempty"`
	XXX_unrecognized  []byte         `json:"-"`
}

func (m *MetricsCollection) Reset()         { *m = MetricsCollection{} }
func (m *MetricsCollection) String() string { return proto.CompactTextString(m) }
func (*MetricsCollection) ProtoMessage()    {}

func (m *MetricsCollection) GetData() []*MetricsData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *MetricsCollection) GetStartTimestampUs() uint64 {
	if m != nil && m.StartTimestampUs != nil {
		return *m.StartTimestampUs
	}
	return 0
}

func (m *MetricsCollection) GetCollectionPointId() string {
	if m != nil && m.CollectionPointId != nil {
		return *m.CollectionPointId
	}
	return ""
}

type MetricsField struct {
	Name             *string                 `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Type             *MetricsField_FieldType `protobuf:"varint,3,opt,name=type,enum=ts_mon.proto.MetricsField_FieldType,def=1" json:"type,omitempty"`
	StringValue      *string                 `protobuf:"bytes,4,opt,name=string_value" json:"string_value,omitempty"`
	IntValue         *int64                  `protobuf:"varint,5,opt,name=int_value" json:"int_value,omitempty"`
	BoolValue        *bool                   `protobuf:"varint,6,opt,name=bool_value" json:"bool_value,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *MetricsField) Reset()         { *m = MetricsField{} }
func (m *MetricsField) String() string { return proto.CompactTextString(m) }
func (*MetricsField) ProtoMessage()    {}

const Default_MetricsField_Type MetricsField_FieldType = MetricsField_STRING

func (m *MetricsField) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MetricsField) GetType() MetricsField_FieldType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_MetricsField_Type
}

func (m *MetricsField) GetStringValue() string {
	if m != nil && m.StringValue != nil {
		return *m.StringValue
	}
	return ""
}

func (m *MetricsField) GetIntValue() int64 {
	if m != nil && m.IntValue != nil {
		return *m.IntValue
	}
	return 0
}

func (m *MetricsField) GetBoolValue() bool {
	if m != nil && m.BoolValue != nil {
		return *m.BoolValue
	}
	return false
}

type PrecomputedDistribution struct {
	SpecType              *PrecomputedDistribution_SpecType `protobuf:"varint,1,opt,name=spec_type,enum=ts_mon.proto.PrecomputedDistribution_SpecType" json:"spec_type,omitempty"`
	Width                 *float64                          `protobuf:"fixed64,2,opt,name=width,def=10" json:"width,omitempty"`
	GrowthFactor          *float64                          `protobuf:"fixed64,3,opt,name=growth_factor,def=0" json:"growth_factor,omitempty"`
	NumBuckets            *int32                            `protobuf:"varint,4,opt,name=num_buckets,def=10" json:"num_buckets,omitempty"`
	LowerBounds           []float64                         `protobuf:"fixed64,5,rep,name=lower_bounds" json:"lower_bounds,omitempty"`
	IsCumulative          *bool                             `protobuf:"varint,6,opt,name=is_cumulative,def=0" json:"is_cumulative,omitempty"`
	Bucket                []int64                           `protobuf:"zigzag64,7,rep,name=bucket" json:"bucket,omitempty"`
	Underflow             *int64                            `protobuf:"zigzag64,8,opt,name=underflow" json:"underflow,omitempty"`
	Overflow              *int64                            `protobuf:"zigzag64,9,opt,name=overflow" json:"overflow,omitempty"`
	Mean                  *float64                          `protobuf:"fixed64,10,opt,name=mean" json:"mean,omitempty"`
	SumOfSquaredDeviation *float64                          `protobuf:"fixed64,11,opt,name=sum_of_squared_deviation" json:"sum_of_squared_deviation,omitempty"`
	XXX_unrecognized      []byte                            `json:"-"`
}

func (m *PrecomputedDistribution) Reset()         { *m = PrecomputedDistribution{} }
func (m *PrecomputedDistribution) String() string { return proto.CompactTextString(m) }
func (*PrecomputedDistribution) ProtoMessage()    {}

const Default_PrecomputedDistribution_Width float64 = 10
const Default_PrecomputedDistribution_GrowthFactor float64 = 0
const Default_PrecomputedDistribution_NumBuckets int32 = 10
const Default_PrecomputedDistribution_IsCumulative bool = false

func (m *PrecomputedDistribution) GetSpecType() PrecomputedDistribution_SpecType {
	if m != nil && m.SpecType != nil {
		return *m.SpecType
	}
	return PrecomputedDistribution_CANONICAL_POWERS_OF_2
}

func (m *PrecomputedDistribution) GetWidth() float64 {
	if m != nil && m.Width != nil {
		return *m.Width
	}
	return Default_PrecomputedDistribution_Width
}

func (m *PrecomputedDistribution) GetGrowthFactor() float64 {
	if m != nil && m.GrowthFactor != nil {
		return *m.GrowthFactor
	}
	return Default_PrecomputedDistribution_GrowthFactor
}

func (m *PrecomputedDistribution) GetNumBuckets() int32 {
	if m != nil && m.NumBuckets != nil {
		return *m.NumBuckets
	}
	return Default_PrecomputedDistribution_NumBuckets
}

func (m *PrecomputedDistribution) GetLowerBounds() []float64 {
	if m != nil {
		return m.LowerBounds
	}
	return nil
}

func (m *PrecomputedDistribution) GetIsCumulative() bool {
	if m != nil && m.IsCumulative != nil {
		return *m.IsCumulative
	}
	return Default_PrecomputedDistribution_IsCumulative
}

func (m *PrecomputedDistribution) GetBucket() []int64 {
	if m != nil {
		return m.Bucket
	}
	return nil
}

func (m *PrecomputedDistribution) GetUnderflow() int64 {
	if m != nil && m.Underflow != nil {
		return *m.Underflow
	}
	return 0
}

func (m *PrecomputedDistribution) GetOverflow() int64 {
	if m != nil && m.Overflow != nil {
		return *m.Overflow
	}
	return 0
}

func (m *PrecomputedDistribution) GetMean() float64 {
	if m != nil && m.Mean != nil {
		return *m.Mean
	}
	return 0
}

func (m *PrecomputedDistribution) GetSumOfSquaredDeviation() float64 {
	if m != nil && m.SumOfSquaredDeviation != nil {
		return *m.SumOfSquaredDeviation
	}
	return 0
}

type MetricsData struct {
	Name                     *string                  `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	MetricNamePrefix         *string                  `protobuf:"bytes,2,opt,name=metric_name_prefix" json:"metric_name_prefix,omitempty"`
	NetworkDevice            *NetworkDevice           `protobuf:"bytes,11,opt,name=network_device" json:"network_device,omitempty"`
	Task                     *Task                    `protobuf:"bytes,12,opt,name=task" json:"task,omitempty"`
	Fields                   []*MetricsField          `protobuf:"bytes,20,rep,name=fields" json:"fields,omitempty"`
	Counter                  *int64                   `protobuf:"varint,30,opt,name=counter" json:"counter,omitempty"`
	Gauge                    *int64                   `protobuf:"varint,32,opt,name=gauge" json:"gauge,omitempty"`
	NoncumulativeDoubleValue *float64                 `protobuf:"fixed64,34,opt,name=noncumulative_double_value" json:"noncumulative_double_value,omitempty"`
	Distribution             *PrecomputedDistribution `protobuf:"bytes,35,opt,name=distribution" json:"distribution,omitempty"`
	StringValue              *string                  `protobuf:"bytes,36,opt,name=string_value" json:"string_value,omitempty"`
	BooleanValue             *bool                    `protobuf:"varint,37,opt,name=boolean_value" json:"boolean_value,omitempty"`
	CumulativeDoubleValue    *float64                 `protobuf:"fixed64,38,opt,name=cumulative_double_value" json:"cumulative_double_value,omitempty"`
	StartTimestampUs         *uint64                  `protobuf:"varint,40,opt,name=start_timestamp_us" json:"start_timestamp_us,omitempty"`
	Units                    *MetricsData_Units       `protobuf:"varint,41,opt,name=units,enum=ts_mon.proto.MetricsData_Units" json:"units,omitempty"`
	XXX_unrecognized         []byte                   `json:"-"`
}

func (m *MetricsData) Reset()         { *m = MetricsData{} }
func (m *MetricsData) String() string { return proto.CompactTextString(m) }
func (*MetricsData) ProtoMessage()    {}

func (m *MetricsData) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MetricsData) GetMetricNamePrefix() string {
	if m != nil && m.MetricNamePrefix != nil {
		return *m.MetricNamePrefix
	}
	return ""
}

func (m *MetricsData) GetNetworkDevice() *NetworkDevice {
	if m != nil {
		return m.NetworkDevice
	}
	return nil
}

func (m *MetricsData) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

func (m *MetricsData) GetFields() []*MetricsField {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *MetricsData) GetCounter() int64 {
	if m != nil && m.Counter != nil {
		return *m.Counter
	}
	return 0
}

func (m *MetricsData) GetGauge() int64 {
	if m != nil && m.Gauge != nil {
		return *m.Gauge
	}
	return 0
}

func (m *MetricsData) GetNoncumulativeDoubleValue() float64 {
	if m != nil && m.NoncumulativeDoubleValue != nil {
		return *m.NoncumulativeDoubleValue
	}
	return 0
}

func (m *MetricsData) GetDistribution() *PrecomputedDistribution {
	if m != nil {
		return m.Distribution
	}
	return nil
}

func (m *MetricsData) GetStringValue() string {
	if m != nil && m.StringValue != nil {
		return *m.StringValue
	}
	return ""
}

func (m *MetricsData) GetBooleanValue() bool {
	if m != nil && m.BooleanValue != nil {
		return *m.BooleanValue
	}
	return false
}

func (m *MetricsData) GetCumulativeDoubleValue() float64 {
	if m != nil && m.CumulativeDoubleValue != nil {
		return *m.CumulativeDoubleValue
	}
	return 0
}

func (m *MetricsData) GetStartTimestampUs() uint64 {
	if m != nil && m.StartTimestampUs != nil {
		return *m.StartTimestampUs
	}
	return 0
}

func (m *MetricsData) GetUnits() MetricsData_Units {
	if m != nil && m.Units != nil {
		return *m.Units
	}
	return MetricsData_UNKNOWN_UNITS
}

func init() {
	proto.RegisterEnum("ts_mon.proto.MetricsField_FieldType", MetricsField_FieldType_name, MetricsField_FieldType_value)
	proto.RegisterEnum("ts_mon.proto.PrecomputedDistribution_SpecType", PrecomputedDistribution_SpecType_name, PrecomputedDistribution_SpecType_value)
	proto.RegisterEnum("ts_mon.proto.MetricsData_Units", MetricsData_Units_name, MetricsData_Units_value)
}