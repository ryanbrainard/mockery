package mockery

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile)

	iface, err := parser.Find("Requester")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type Requester struct {
	mock.Mock
}

func (m *Requester) Name_Get() string {
	return "Get"
}
func (m *Requester) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *Requester) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *Requester) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *Requester) Get(path string) (string, error) {
	ret := m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorSingleReturn(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile2)

	iface, err := parser.Find("Requester2")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type Requester2 struct {
	mock.Mock
}

func (m *Requester2) Name_Get() string {
	return "Get"
}
func (m *Requester2) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *Requester2) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *Requester2) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *Requester2) Get(path string) error {
	ret := m.Called(path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNoArguments(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester3.go"))

	iface, err := parser.Find("Requester3")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type Requester3 struct {
	mock.Mock
}

func (m *Requester3) Name_Get() string {
	return "Get"
}
func (m *Requester3) MockOn_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester3) MockOnTyped_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester3) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester3) Get() error {
	ret := m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNoNothing(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester4.go"))

	iface, err := parser.Find("Requester4")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type Requester4 struct {
	mock.Mock
}

func (m *Requester4) Name_Get() string {
	return "Get"
}
func (m *Requester4) MockOn_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester4) MockOnTyped_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester4) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *Requester4) Get() {
	m.Called()
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorUnexported(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_unexported.go"))

	iface, err := parser.Find("requester")

	gen := NewGenerator(iface)
	gen.ip = true

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type mockRequester struct {
	mock.Mock
}

func (m *mockRequester) Name_Get() string {
	return "Get"
}
func (m *mockRequester) MockOn_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *mockRequester) MockOnTyped_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *mockRequester) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get")
}
func (m *mockRequester) Get() {
	m.Called()
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorPrologue(t *testing.T) {
	parser := NewParser()
	parser.Parse(testFile)

	iface, err := parser.Find("Requester")
	assert.NoError(t, err)

	gen := NewGenerator(iface)

	gen.GeneratePrologue("mocks")

	expected := `package mocks

import "github.com/ryanbrainard/mockery/mockery/fixtures"
import "github.com/stretchr/testify/mock"

`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorProloguewithImports(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ns.go"))

	iface, err := parser.Find("RequesterNS")
	assert.NoError(t, err)

	gen := NewGenerator(iface)

	gen.GeneratePrologue("mocks")

	expected := `package mocks

import "github.com/ryanbrainard/mockery/mockery/fixtures"
import "github.com/stretchr/testify/mock"

import "net/http"

`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorPointers(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ptr.go"))

	iface, err := parser.Find("RequesterPtr")
	assert.NoError(t, err)

	gen := NewGenerator(iface)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterPtr struct {
	mock.Mock
}

func (m *RequesterPtr) Name_Get() string {
	return "Get"
}
func (m *RequesterPtr) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterPtr) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterPtr) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *RequesterPtr) Get(path string) (*string, error) {
	ret := m.Called(path)

	var r0 *string
	if rf, ok := ret.Get(0).(func(string) *string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorSlice(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_slice.go"))

	iface, err := parser.Find("RequesterSlice")
	assert.NoError(t, err)

	gen := NewGenerator(iface)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterSlice struct {
	mock.Mock
}

func (m *RequesterSlice) Name_Get() string {
	return "Get"
}
func (m *RequesterSlice) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterSlice) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterSlice) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *RequesterSlice) Get(path string) ([]string, error) {
	ret := m.Called(path)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorArrayLiteralLen(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_array.go"))

	iface, err := parser.Find("RequesterArray")
	assert.NoError(t, err)

	gen := NewGenerator(iface)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterArray struct {
	mock.Mock
}

func (m *RequesterArray) Name_Get() string {
	return "Get"
}
func (m *RequesterArray) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterArray) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterArray) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *RequesterArray) Get(path string) ([2]string, error) {
	ret := m.Called(path)

	var r0 [2]string
	if rf, ok := ret.Get(0).(func(string) [2]string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([2]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorVarArg(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_vararg.go"))

	iface, err := parser.Find("RequesterVarArg")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterVarArg struct {
	mock.Mock
}

func (m *RequesterVarArg) Name_Get() string {
	return "Get"
}
func (m *RequesterVarArg) MockOn_Get(paths interface{}) *mock.Call {
	return m.Mock.On("Get", paths)
}
func (m *RequesterVarArg) MockOnTyped_Get(paths ...string) *mock.Call {
	return m.Mock.On("Get", paths)
}
func (m *RequesterVarArg) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *RequesterVarArg) Get(paths ...string) error {
	ret := m.Called(paths)

	var r0 error
	if rf, ok := ret.Get(0).(func(...string) error); ok {
		r0 = rf(paths...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorNamespacedTypes(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_ns.go"))

	iface, err := parser.Find("RequesterNS")
	assert.NoError(t, err)

	gen := NewGenerator(iface)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterNS struct {
	mock.Mock
}

func (m *RequesterNS) Name_Get() string {
	return "Get"
}
func (m *RequesterNS) MockOn_Get(path interface{}) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterNS) MockOnTyped_Get(path string) *mock.Call {
	return m.Mock.On("Get", path)
}
func (m *RequesterNS) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything)
}
func (m *RequesterNS) Get(path string) (http.Response, error) {
	ret := m.Called(path)

	var r0 http.Response
	if rf, ok := ret.Get(0).(func(string) http.Response); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(http.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorHavingNoNamesOnArguments(t *testing.T) {
	parser := NewParser()

	parser.Parse(filepath.Join(fixturePath, "custom_error.go"))

	iface, err := parser.Find("KeyManager")
	assert.NoError(t, err)

	gen := NewGenerator(iface)
	assert.NoError(t, err)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type KeyManager struct {
	mock.Mock
}

func (m *KeyManager) Name_GetKey() string {
	return "GetKey"
}
func (m *KeyManager) MockOn_GetKey(_a0 interface{}, _a1 interface{}) *mock.Call {
	return m.Mock.On("GetKey", _a0, _a1)
}
func (m *KeyManager) MockOnTyped_GetKey(_a0 string, _a1 uint16) *mock.Call {
	return m.Mock.On("GetKey", _a0, _a1)
}
func (m *KeyManager) MockOnAny_GetKey() *mock.Call {
	return m.Mock.On("GetKey", mock.Anything, mock.Anything)
}
func (m *KeyManager) GetKey(_a0 string, _a1 uint16) ([]byte, *test.Err) {
	ret := m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, uint16) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 *test.Err
	if rf, ok := ret.Get(1).(func(string, uint16) *test.Err); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*test.Err)
		}
	}

	return r0, r1
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorElidedType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "requester_elided.go"))

	iface, err := parser.Find("RequesterElided")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type RequesterElided struct {
	mock.Mock
}

func (m *RequesterElided) Name_Get() string {
	return "Get"
}
func (m *RequesterElided) MockOn_Get(path interface{}, url interface{}) *mock.Call {
	return m.Mock.On("Get", path, url)
}
func (m *RequesterElided) MockOnTyped_Get(path string, url string) *mock.Call {
	return m.Mock.On("Get", path, url)
}
func (m *RequesterElided) MockOnAny_Get() *mock.Call {
	return m.Mock.On("Get", mock.Anything, mock.Anything)
}
func (m *RequesterElided) Get(path string, url string) error {
	ret := m.Called(path, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(path, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorFuncType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "func_type.go"))

	iface, err := parser.Find("Fooer")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type Fooer struct {
	mock.Mock
}

func (m *Fooer) Name_Foo() string {
	return "Foo"
}
func (m *Fooer) MockOn_Foo(f interface{}) *mock.Call {
	return m.Mock.On("Foo", f)
}
func (m *Fooer) MockOnTyped_Foo(f func(string) string) *mock.Call {
	return m.Mock.On("Foo", f)
}
func (m *Fooer) MockOnAny_Foo() *mock.Call {
	return m.Mock.On("Foo", mock.Anything)
}
func (m *Fooer) Foo(f func(string) string) error {
	ret := m.Called(f)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(string) string) error); ok {
		r0 = rf(f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (m *Fooer) Name_Bar() string {
	return "Bar"
}
func (m *Fooer) MockOn_Bar(f interface{}) *mock.Call {
	return m.Mock.On("Bar", f)
}
func (m *Fooer) MockOnTyped_Bar(f func([]int) ) *mock.Call {
	return m.Mock.On("Bar", f)
}
func (m *Fooer) MockOnAny_Bar() *mock.Call {
	return m.Mock.On("Bar", mock.Anything)
}
func (m *Fooer) Bar(f func([]int) ) {
	m.Called(f)
}
func (m *Fooer) Name_Baz() string {
	return "Baz"
}
func (m *Fooer) MockOn_Baz(path interface{}) *mock.Call {
	return m.Mock.On("Baz", path)
}
func (m *Fooer) MockOnTyped_Baz(path string) *mock.Call {
	return m.Mock.On("Baz", path)
}
func (m *Fooer) MockOnAny_Baz() *mock.Call {
	return m.Mock.On("Baz", mock.Anything)
}
func (m *Fooer) Baz(path string) func(string) string {
	ret := m.Called(path)

	var r0 func(string) string
	if rf, ok := ret.Get(0).(func(string) func(string) string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(string) string)
		}
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}

func TestGeneratorChanType(t *testing.T) {
	parser := NewParser()
	parser.Parse(filepath.Join(fixturePath, "async.go"))

	iface, err := parser.Find("AsyncProducer")

	gen := NewGenerator(iface)

	err = gen.Generate()
	assert.NoError(t, err)

	expected := `type AsyncProducer struct {
	mock.Mock
}

func (m *AsyncProducer) Name_Input() string {
	return "Input"
}
func (m *AsyncProducer) MockOn_Input() *mock.Call {
	return m.Mock.On("Input")
}
func (m *AsyncProducer) MockOnTyped_Input() *mock.Call {
	return m.Mock.On("Input")
}
func (m *AsyncProducer) MockOnAny_Input() *mock.Call {
	return m.Mock.On("Input")
}
func (m *AsyncProducer) Input() chan<- bool {
	ret := m.Called()

	var r0 chan<- bool
	if rf, ok := ret.Get(0).(func() chan<- bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan<- bool)
		}
	}

	return r0
}
func (m *AsyncProducer) Name_Output() string {
	return "Output"
}
func (m *AsyncProducer) MockOn_Output() *mock.Call {
	return m.Mock.On("Output")
}
func (m *AsyncProducer) MockOnTyped_Output() *mock.Call {
	return m.Mock.On("Output")
}
func (m *AsyncProducer) MockOnAny_Output() *mock.Call {
	return m.Mock.On("Output")
}
func (m *AsyncProducer) Output() <-chan bool {
	ret := m.Called()

	var r0 <-chan bool
	if rf, ok := ret.Get(0).(func() <-chan bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan bool)
		}
	}

	return r0
}
func (m *AsyncProducer) Name_Whatever() string {
	return "Whatever"
}
func (m *AsyncProducer) MockOn_Whatever() *mock.Call {
	return m.Mock.On("Whatever")
}
func (m *AsyncProducer) MockOnTyped_Whatever() *mock.Call {
	return m.Mock.On("Whatever")
}
func (m *AsyncProducer) MockOnAny_Whatever() *mock.Call {
	return m.Mock.On("Whatever")
}
func (m *AsyncProducer) Whatever() chan bool {
	ret := m.Called()

	var r0 chan bool
	if rf, ok := ret.Get(0).(func() chan bool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan bool)
		}
	}

	return r0
}
`

	assert.Equal(t, expected, gen.buf.String())
}
