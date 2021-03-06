// This file a slightly adapted version of the redis packges `scan.go` with the
// same name.

package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"

	"hotpu.cn/xkefu/common/log"
)

var v interface{}

var (
	TypeOfTime = reflect.TypeOf(time.Time{})
	// TypeOfLocation            = reflect.TypeOf(Location{})
	// TypeOfLocationPtr         = reflect.TypeOf(&Location{})
	// TypeOfGender              = reflect.TypeOf(Gender("foo"))
	// TypeOfRole                = reflect.TypeOf(Role("foo"))
	// TypeOfMessageTextBody     = reflect.TypeOf(TextMessage{})
	// TypeOfMessageTextBodyPtr  = reflect.TypeOf(&TextMessage{})
	// TypeOfMessageImageBody    = reflect.TypeOf(ImageMessage{})
	// TypeOfMessageImageBodyPtr = reflect.TypeOf(&ImageMessage{})
	TypeOfInterface = reflect.TypeOf(v)
	// TypeOfPlatformSlice       = reflect.TypeOf([]Platform{})
	TypeOfStringSlice = reflect.TypeOf([]string{})
	// TypeOfRawMessage          = reflect.TypeOf(json.RawMessage{})
	// TypeOfRawMessagePtr       = reflect.TypeOf(&json.RawMessage{})
)

func ensureLen(d reflect.Value, n int) {
	if n > d.Cap() {
		d.Set(reflect.MakeSlice(d.Type(), n, n))
	} else {
		d.SetLen(n)
	}
}

func cannotConvert(d reflect.Value, s interface{}) error {
	return fmt.Errorf("redigo: Scan cannot convert from %s to %s",
		reflect.TypeOf(s), d.Type())
}

func convertAssignBytes(d reflect.Value, s []byte) (err error) {
	switch d.Type().Kind() {
	case reflect.Float32, reflect.Float64:
		var x float64
		x, err = strconv.ParseFloat(string(s), d.Type().Bits())
		d.SetFloat(x)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		x, err = strconv.ParseInt(string(s), 10, d.Type().Bits())
		d.SetInt(x)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x uint64
		x, err = strconv.ParseUint(string(s), 10, d.Type().Bits())
		d.SetUint(x)

	case reflect.Bool:
		var x bool
		x, err = strconv.ParseBool(string(s))
		d.SetBool(x)

	case reflect.String:
		d.SetString(string(s))

	case reflect.Slice:
		if d.Type().Elem().Kind() != reflect.Uint8 {

			// Convert comma separeted string to Platform slice
			// if d.Type() == TypeOfPlatformSlice {
			// 	ss := strings.Split(string(s), ",")
			// 	platforms := make([]Platform, 0)

			// 	for _, s := range ss {
			// 		if trimed := strings.TrimSpace(s); trimed == "" {
			// 			continue
			// 		}
			// 		platforms = append(platforms, Platform(s))
			// 	}
			// 	d.Set(reflect.ValueOf(platforms))
			// 	return
			// } else if d.Type() == TypeOfStringSlice {
			// 	log.Debug("Met string slice: %v", d)
			// 	return
			// }

			err = cannotConvert(d, s)
		} else {
			d.SetBytes(s)
		}

	// -- ------------------------------------------------------------------
	// -- PATCH (from Redis to Go objects)
	// -- ------------------------------------------------------------------
	case reflect.Interface:

		// handle the message body, if it is a message body, then it should
		// have a non-empty field 'redis_message_type'
		var temp map[string]interface{}

		//invalid character '&' looking for beginning of value,so take out '&'
		ss := []byte(strings.TrimPrefix(string(s), "&"))

		err = json.Unmarshal(ss, &temp)
		if err != nil {
			log.Error("common.store.redis", "convertAssignBytes", "Error when unmarhaling redis value %s to a map[string]interface{}", string(ss))
		}

		// if _, ok := temp["redis_message_type"]; ok {
		// 	log.Debug("Detected message body: %s, trying to convert to the proper type.", ss)
		// 	var body RedisMessageBody
		// 	if err = json.Unmarshal(ss, &body); err == nil {
		// 		d.Set(reflect.ValueOf(body.Body))
		// 		return
		// 	} else {
		// 		log.Error("Error when unmarshaling to RedisMessageBody object(%s).", err.Error())
		// 	}
		// }

	default:
		switch d.Type() {
		case TypeOfTime:
			// convert a RFC3339 date time string to time.Time object
			var t time.Time
			if t, err = time.Parse(time.RFC3339, string(s)); err == nil {
				d.Set(reflect.ValueOf(t))
				return
			}

		// case TypeOfLocation:
		// 	// convert a json marshalled json string to a Location object
		// 	var location Location
		// 	if err = json.Unmarshal(s, &location); err == nil {
		// 		d.Set(reflect.ValueOf(location))
		// 		return
		// 	} else {
		// 		log.Error("Error: %s", err.Error())
		// 	}

		// case TypeOfGender:
		// 	// do nothing

		// case TypeOfLocationPtr:
		// 	// convert a json marshaled josn string to a location object
		// 	var location Location
		// 	if err = json.Unmarshal(s, &location); err == nil {
		// 		d.Set(reflect.ValueOf(&location))
		// 		return
		// 	} else {
		// 		log.Error("Error when unmarshalling: %s", err.Error())
		// 	}

		// case TypeOfRawMessage:
		// 	// convert a byte[] to json.RawMessage
		// 	val := json.RawMessage([]byte(s))
		// 	d.Set(reflect.ValueOf(val))
		// 	return

		// case TypeOfRawMessagePtr:
		// 	// convert a byte[] to *json.RawMessage
		// 	val := json.RawMessage([]byte(s))
		// 	d.Set(reflect.ValueOf(&val))
		// 	return

		default:
			if d.Kind() == reflect.Interface {

			}
		} // end of switch

		err = cannotConvert(d, s)
	}
	return
}

func convertAssignInt(d reflect.Value, s int64) (err error) {
	switch d.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.SetInt(s)
		if d.Int() != s {
			err = strconv.ErrRange
			d.SetInt(0)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if s < 0 {
			err = strconv.ErrRange
		} else {
			x := uint64(s)
			d.SetUint(x)
			if d.Uint() != x {
				err = strconv.ErrRange
				d.SetUint(0)
			}
		}
	case reflect.Bool:
		d.SetBool(s != 0)
	default:
		err = cannotConvert(d, s)
	}
	return
}

func convertAssignValue(d reflect.Value, s interface{}) (err error) {
	switch s := s.(type) {
	case []byte:
		err = convertAssignBytes(d, s)
	case int64:
		err = convertAssignInt(d, s)
	default:
		err = cannotConvert(d, s)
	}
	return err
}

func convertAssignValues(d reflect.Value, s []interface{}) error {
	if d.Type().Kind() != reflect.Slice {
		return cannotConvert(d, s)
	}
	ensureLen(d, len(s))
	for i := 0; i < len(s); i++ {
		if err := convertAssignValue(d.Index(i), s[i]); err != nil {
			return err
		}
	}
	return nil
}

func convertAssign(d interface{}, s interface{}) (err error) {
	// Handle the most common destination types using type switches and
	// fall back to reflection for all other types.
	switch s := s.(type) {
	case nil:
		// ingore
	case []byte:
		switch d := d.(type) {
		case *string:
			*d = string(s)
		case *int:
			*d, err = strconv.Atoi(string(s))
		case *bool:
			*d, err = strconv.ParseBool(string(s))
		case *[]byte:
			*d = s
		case *interface{}:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
				err = cannotConvert(d, s)
			} else {
				err = convertAssignBytes(d.Elem(), s)
			}
		}
	case int64:
		switch d := d.(type) {
		case *int:
			x := int(s)
			if int64(x) != s {
				err = strconv.ErrRange
				x = 0
			}
			*d = x
		case *bool:
			*d = s != 0
		case *interface{}:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
				err = cannotConvert(d, s)
			} else {
				err = convertAssignInt(d.Elem(), s)
			}
		}
	case []interface{}:
		switch d := d.(type) {
		case *[]interface{}:
			*d = s
		case *interface{}:
			*d = s
		case nil:
			// skip value
		default:
			if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
				err = cannotConvert(d, s)
			} else {
				err = convertAssignValues(d.Elem(), s)
			}
		}
	case redis.Error:
		err = s
	default:
		err = cannotConvert(reflect.ValueOf(d), s)
	}
	return
}

// Scan copies from src to the values pointed at by dest.
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface{} or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
func Scan(src []interface{}, dest ...interface{}) ([]interface{}, error) {
	if len(src) < len(dest) {
		return nil, errors.New("redigo: Scan array short")
	}
	var err error
	for i, d := range dest {
		err = convertAssign(d, src[i])
		if err != nil {
			break
		}
	}
	return src[len(dest):],
		err
}

type fieldSpec struct {
	name  string
	index []int
	//omitEmpty bool
}

type structSpec struct {
	m map[string]*fieldSpec
	l []*fieldSpec
}

func (ss *structSpec) fieldSpec(name []byte) *fieldSpec {
	return ss.m[string(name)]
}

func compileStructSpec(t reflect.Type, depth map[string]int, index []int, ss *structSpec) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch {
		case f.PkgPath != "":
			// Ignore unexported fields.
		case f.Anonymous:
			// TODO: Handle pointers. Requires change to decoder and
			// protection against infinite recursion.
			if f.Type.Kind() == reflect.Struct {
				compileStructSpec(f.Type, depth, append(index, i), ss)
			}
		default:
			fs := &fieldSpec{name: f.Name}
			tag := f.Tag.Get("redis")
			p := strings.Split(tag, ",")
			if len(p) > 0 {
				if p[0] == "-" {
					continue
				}
				if len(p[0]) > 0 {
					fs.name = p[0]
				}
				for _, s := range p[1:] {
					switch s {
					//case "omitempty":
					//  fs.omitempty = true
					default:
						panic(errors.New("redigo: unknown field flag " + s + " for type " + t.Name()))
					}
				}
			}
			d, found := depth[fs.name]
			if !found {
				d = 1 << 30
			}
			switch {
			case len(index) == d:
				// At same depth, remove from result.
				delete(ss.m, fs.name)
				j := 0
				for i := 0; i < len(ss.l); i++ {
					if fs.name != ss.l[i].name {
						ss.l[j] = ss.l[i]
						j += 1
					}
				}
				ss.l = ss.l[:j]
			case len(index) < d:
				fs.index = make([]int, len(index)+1)
				copy(fs.index, index)
				fs.index[len(index)] = i
				depth[fs.name] = len(index)
				ss.m[fs.name] = fs
				ss.l = append(ss.l, fs)
			}
		}
	}
}

var (
	structSpecMutex  sync.RWMutex
	structSpecCache  = make(map[reflect.Type]*structSpec)
	defaultFieldSpec = &fieldSpec{}
)

func structSpecForType(t reflect.Type) *structSpec {

	structSpecMutex.RLock()
	ss, found := structSpecCache[t]
	structSpecMutex.RUnlock()
	if found {
		return ss
	}

	structSpecMutex.Lock()
	defer structSpecMutex.Unlock()
	ss, found = structSpecCache[t]
	if found {
		return ss
	}

	ss = &structSpec{m: make(map[string]*fieldSpec)}
	compileStructSpec(t, make(map[string]int), nil, ss)
	structSpecCache[t] = ss
	return ss
}

var errScanStructValue = errors.New("redigo: ScanStruct value must be non-nil pointer to a struct")

// ScanStruct scans alternating names and values from src to a struct. The
// HGETALL and CONFIG GET commands return replies in this format.
//
// ScanStruct uses exported field names to match values in the response. Use
// 'redis' field tag to override the name:
//
//      Field int `redis:"myName"`
//
// Fields with the tag redis:"-" are ignored.
//
// Integer, float, boolean, string and []byte fields are supported. Scan uses the
// standard strconv package to convert bulk string values to numeric and
// boolean types.
//
// If a src element is nil, then the corresponding field is not modified.
func ScanStruct(src []interface{}, dest interface{}) error {
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() {
		return errScanStructValue
	}
	d = d.Elem()
	if d.Kind() != reflect.Struct {
		return errScanStructValue
	}
	ss := structSpecForType(d.Type())

	if len(src)%2 != 0 {
		return errors.New("redigo: ScanStruct expects even number of values in values")
	}

	for i := 0; i < len(src); i += 2 {
		s := src[i+1]
		if s == nil {
			continue
		}
		name, ok := src[i].([]byte)
		if !ok {
			return errors.New("redigo: ScanStruct key not a bulk string value")
		}
		fs := ss.fieldSpec(name)
		if fs == nil {
			continue
		}
		if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil {
			return err
		}
	}
	return nil
}

var (
	errScanSliceValue = errors.New("redigo: ScanSlice dest must be non-nil pointer to a struct")
)

// ScanSlice scans src to the slice pointed to by dest. The elements the dest
// slice must be integer, float, boolean, string, struct or pointer to struct
// values.
//
// Struct fields must be integer, float, boolean or string values. All struct
// fields are used unless a subset is specified using fieldNames.
func ScanSlice(src []interface{}, dest interface{}, fieldNames ...string) error {
	d := reflect.ValueOf(dest)
	if d.Kind() != reflect.Ptr || d.IsNil() {
		return errScanSliceValue
	}
	d = d.Elem()
	if d.Kind() != reflect.Slice {
		return errScanSliceValue
	}

	isPtr := false
	t := d.Type().Elem()
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		isPtr = true
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		ensureLen(d, len(src))
		for i, s := range src {
			if s == nil {
				continue
			}
			if err := convertAssignValue(d.Index(i), s); err != nil {
				return err
			}
		}
		return nil
	}

	ss := structSpecForType(t)
	fss := ss.l
	if len(fieldNames) > 0 {
		fss = make([]*fieldSpec, len(fieldNames))
		for i, name := range fieldNames {
			fss[i] = ss.m[name]
			if fss[i] == nil {
				return errors.New("redigo: ScanSlice bad field name " + name)
			}
		}
	}

	if len(fss) == 0 {
		return errors.New("redigo: ScanSlice no struct fields")
	}

	n := len(src) / len(fss)
	if n*len(fss) != len(src) {
		return errors.New("redigo: ScanSlice length not a multiple of struct field count")
	}

	ensureLen(d, n)
	for i := 0; i < n; i++ {
		d := d.Index(i)
		if isPtr {
			if d.IsNil() {
				d.Set(reflect.New(t))
			}
			d = d.Elem()
		}
		for j, fs := range fss {
			s := src[i*len(fss)+j]
			if s == nil {
				continue
			}
			if err := convertAssignValue(d.FieldByIndex(fs.index), s); err != nil {
				return err
			}
		}
	}
	return nil
}

// MyArgs is a helper for constructing command arguments from structured values.
type MyArgs []interface{}

// Add returns the result of appending value to args.
func (args MyArgs) Add(value ...interface{}) MyArgs {
	for _, v := range value {
		rv := reflect.ValueOf(v)
		switch rv.Type() {
		// // convert json.RawMessage to byte[]
		// case TypeOfRawMessage:
		// 	p, ok := v.(json.RawMessage)
		// 	if ok {
		// 		args = append(args, []byte(p))
		// 		continue
		// 	}

		// // convert *json.RawMessage to byte[]
		// case TypeOfRawMessagePtr:
		// 	p, ok := v.(*json.RawMessage)
		// 	if ok {
		// 		args = append(args, []byte(*p))
		// 		continue
		// 	}

		default:
			args = append(args, v)
			continue
		}
	}

	return args
}

// AddFlat returns the result of appending the flattened value of v to args.
//
// Maps are flattened by appending the alternating keys and map values to args.
//
// Slices are flattened by appending the slice elements to args.
//
// Structs are flattened by appending the alternating names and values of
// exported fields to args. If v is a nil struct pointer, then nothing is
// appended. The 'redis' field tag overrides struct field names. See ScanStruct
// for more information on the use of the 'redis' field tag.
//
// Other types are appended to args as is.
func (args MyArgs) AddFlat(v interface{}) MyArgs {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Struct:
		args = flattenStruct(args, rv)
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			args = append(args, rv.Index(i).Interface())
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			args = append(args, k.Interface(), rv.MapIndex(k).Interface())
		}
	case reflect.Ptr:
		if rv.Type().Elem().Kind() == reflect.Struct {
			if !rv.IsNil() {
				args = flattenStruct(args, rv.Elem())
			}
		} else {
			args = append(args, v)
		}
	default:
		args = append(args, v)
	}
	return args
}

func flattenStruct(args MyArgs, v reflect.Value) MyArgs {
	ss := structSpecForType(v.Type())
	for _, fs := range ss.l {
		fv := v.FieldByIndex(fs.index)

		// -- ------------------------------------------------------------------
		// -- PATCH (from go objects to Redis)
		// -- ------------------------------------------------------------------
		switch fv.Type() {
		// convert time.Time to time formatted string(RFC3339)
		case TypeOfTime:
			t, ok := fv.Interface().(time.Time)
			if ok {
				tStr := t.Format(time.RFC3339)
				log.Debug("common.store.redis", "flattenStruct", "Flatten struct time is: %s", tStr)
				args = append(args, fs.name, tStr)
				continue
			}

		// // convert location object to marshaled json byte array
		// case TypeOfLocation:
		// 	l, ok := fv.Interface().(Location)
		// 	if ok {
		// 		body, err := json.Marshal(l)
		// 		if err == nil {
		// 			log.Debug("Flatten struct location is: %s", body)
		// 			args = append(args, fs.name, body)
		// 			continue
		// 		} else {
		// 			log.Error("Error when marshal location object: %#v(error: %s)", l, err.Error())
		// 		}
		// 	}

		// // convert location ptr object to marshaled json byte array
		// case TypeOfLocationPtr:
		// 	l, ok := fv.Interface().(*Location)
		// 	if ok {
		// 		body, err := json.Marshal(l)
		// 		if err == nil {
		// 			log.Debug("Flatten struct location pointer is: %s", body)
		// 			args = append(args, fs.name, body)
		// 			continue
		// 		} else {
		// 			log.Error("Error when marshal location object: %#v(error: %s)", l, err.Error())
		// 		}
		// 	}

		// // Convert Platform slice to comma seperated string
		// case TypeOfPlatformSlice:
		// 	platforms, ok := fv.Interface().([]Platform)
		// 	log.Warning("Platforms: %#v", platforms)

		// 	if ok {
		// 		ss := make([]string, 0)
		// 		for _, p := range platforms {
		// 			ss = append(ss, string(p))
		// 		}
		// 		body := strings.Join(ss, ",")
		// 		log.Debug("Flatten struct platform slice is: %s", body)
		// 		args = append(args, fs.name, body)
		// 		continue
		// 	}

		// 	// convert json.RawMessage to byte[]
		// case TypeOfRawMessage:
		// 	p, ok := fv.Interface().(json.RawMessage)
		// 	b := []byte(p)
		// 	if ok {
		// 		log.Debug("Flatten struct raw message is: %s", b)
		// 		args = append(args, fs.name, b)
		// 		continue
		// 	}

		// 	// convert *json.RawMessage to byte[]
		// case TypeOfRawMessagePtr:
		// 	p, ok := fv.Interface().(*json.RawMessage)
		// 	b := []byte(*p)
		// 	if ok {
		// 		log.Debug("Flatten struct raw message pointer is: %s", b)
		// 		args = append(args, fs.name, b)
		// 		continue
		// 	}

		default:
			if fv.Kind() == reflect.Interface {
				switch fv.Interface().(type) {
				// convert TextMessage to marshaled json byte array
				// case TextMessage:
				// 	data, ok := fv.Interface().(TextMessage)
				// 	if ok {
				// 		log.Debug("Flatten struct text message data is: %+v", data)
				// 		wrapData := RedisMessageBody{MessageTypeText, data}
				// 		body, err := json.Marshal(wrapData)
				// 		if err == nil {
				// 			log.Debug("Flatten struct text message is: %s", body)
				// 			args = append(args, fs.name, body)
				// 			continue
				// 		}
				// 	}

				// // convert *TextMessage to marshaled json byte array
				// case *TextMessage:
				// 	data, ok := fv.Interface().(*TextMessage)
				// 	if ok {
				// 		log.Debug("Flatten struct text message pointer data is: %+v", data)
				// 		wrapData := RedisMessageBody{MessageTypeText, data}
				// 		body, err := json.Marshal(wrapData)
				// 		if err == nil {
				// 			log.Debug("Flatten struct text message pointer is: %s", body)
				// 			args = append(args, fs.name, body)
				// 			continue
				// 		}
				// 	}

				// // convert ImageMessage to marshaled json byte array
				// case ImageMessage:
				// 	data, ok := fv.Interface().(ImageMessage)
				// 	if ok {
				// 		wrapData := RedisMessageBody{MessageTypeImage, data}
				// 		body, err := json.Marshal(wrapData)
				// 		if err == nil {
				// 			log.Debug("Flatten struct image message is: %s", body)
				// 			args = append(args, fs.name, body)
				// 			continue
				// 		}
				// 	}

				// // convert *ImageMessage to marshaled json byte array
				// case *ImageMessage:
				// 	data, ok := fv.Interface().(*ImageMessage)
				// 	if ok {
				// 		wrapData := RedisMessageBody{MessageTypeImage, data}
				// 		body, err := json.Marshal(wrapData)
				// 		if err == nil {
				// 			log.Debug("Flatten struct image message pointer is: %s", body)
				// 			args = append(args, fs.name, body)
				// 			continue
				// 		}
				// 	}

				default:
					log.Debug("common.store.redis", "flattenStruct", "Unhandled interface value: %#v", fv.Interface())
				}
			}
		} // end of swtich statement

		args = append(args, fs.name, fv.Interface())
	}
	return args
}
