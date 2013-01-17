// The package gl provides Go bindings for OpenGL.
// Some of the more awkward parts of the library are wrapped to provide idiomatic Go behaviour for e.g. error handling.
// Constants have their GL_ prefix removed when possible, i.e. unless they start with a number.
// This package uses the intersection of OpenGL 2.1 and OpenGL 3.2 core. Legacy features are not retained.
package gl

// #cgo darwin CFLAGS: -I/opt/local/include/
// #cgo linux LDFLAGS: -lGLEW 
// #cgo darwin LDFLAGS: -lGLEW -L/opt/local/lib/ -framework OpenGL
// #include <GL/glew.h>
// #undef GLEW_GET_FUN
// #define GLEW_GET_FUN(x) (*x)
import "C"
import "unsafe"
import "reflect"
import "errors"
import "image"
import "image/color"

func Init() {
	C.glewInit()
}

// Enable calls glEnable
func Enable(mask int) {
	C.glEnable(C.GLenum(mask))
}

// Disable calls glDisable
func Disable(mask int) {
	C.glDisable(C.GLenum(mask))
}

// ClearColor calls glClearColor
func ClearColor(r float64, g float64, b float64, a float64) {
	C.glClearColor(C.GLclampf(r), C.GLclampf(g), C.GLclampf(b), C.GLclampf(a))
}

// Clear calls glClear
func Clear(mask int) {
	C.glClear(C.GLbitfield(mask))
}

// Viewport calls glViewport
func Viewport(x int, y int, w int, h int) {
	C.glViewport(C.GLint(x), C.GLint(y), C.GLsizei(w), C.GLsizei(h))
}

// Add the interface to export image
func ReadPixels(x int, y int, w int, h int) image.Image {
	data := make([]uint16, 4*w*h)
	p := unsafe.Pointer(&data[0])
	C.glReadPixels(C.GLint(x), C.GLint(y), C.GLsizei(w), C.GLsizei(h), C.GLenum(RGBA), C.GLenum(UNSIGNED_SHORT), p)
	rec := image.Rect(x, y, x+w, y+h)
	rgba := image.NewRGBA64(rec)
	for i := 0; i < w*h; i++ {
		c := color.RGBA64{data[4*i], data[4*i+1], data[4*i+2], data[4*i+3]}
		rgba.Set(i%w, h-i/w, c)
	}
	return rgba
}

// DepthRange calls glDepthRange
func DepthRange(zNear, zFar float64) {
	C.glDepthRange(C.GLclampd(zNear), C.GLclampd(zFar))
}

// BlendFunc calls glBlendFunc
func BlendFunc(sfactor, dfactor int) {
	C.glBlendFunc(C.GLenum(sfactor), C.GLenum(dfactor))
}

// PolygonMode calls glPolygonMode
func PolygonMode(face, mode int) {
	C.glPolygonMode(C.GLenum(face), C.GLenum(mode))
}

// ColorMask calls glColorMask
func ColorMask(r, g, b, a bool) {
	R, G, B, A := FALSE, FALSE, FALSE, FALSE
	if r {
		R = TRUE
	}
	if g {
		G = TRUE
	}
	if b {
		B = TRUE
	}
	if a {
		A = TRUE
	}
	C.glColorMask(C.GLboolean(R), C.GLboolean(G), C.GLboolean(B), C.GLboolean(A))
}

func toCtype(data interface{}) (p unsafe.Pointer, t C.GLenum, ts int, s uintptr) {
	v := reflect.ValueOf(data)
	var et reflect.Type
	switch v.Type().Kind() {
	case reflect.Slice, reflect.Array:
		if !v.IsNil() {
			p = unsafe.Pointer(v.Index(0).UnsafeAddr())
			s = uintptr(v.Len())
		}
		et = v.Type().Elem()
	default:
		panic("not a pointer or slice")
	}
	switch et.Kind() {
	case reflect.Uint8:
		t = UNSIGNED_BYTE
	case reflect.Int8:
		t = BYTE
	case reflect.Uint16:
		t = UNSIGNED_SHORT
	case reflect.Int16:
		t = SHORT
	case reflect.Uint32:
		t = UNSIGNED_INT
	case reflect.Int32:
		t = INT
	case reflect.Float32:
		t = FLOAT
	case reflect.Float64:
		t = DOUBLE
	default:
		panic("unknown type: " + reflect.TypeOf(v).String())
	}
	ts = et.Bits() / 8
	s *= uintptr(et.Bits() / 8)

	return
}

// the type Buffer represents a buffer object
type Buffer struct {
	i  C.GLuint
	t  C.GLenum
	ts int
}

// NewBuffer creates a new buffer using glGenBuffers. If targ is not 0, it will call Buffer.Set with the given parameters.
func NewBuffer(targ int, data interface{}, usage int) *Buffer {
	var buf C.GLuint

	C.glGenBuffers(1, &buf)
	buff := &Buffer{}
	buff.i = buf
	if targ != 0 {
		buff.Set(targ, data, usage)
	}
	return buff
}

//DeleteBuffer delete the buffer using glDeleteBuffer
func DeleteBuffers(buffers ...*Buffer) {
	for _, buf := range buffers {
		C.glDeleteBuffers(1, &(buf.i))
	}
}

// Set calls glBufferData with appropriate arguments to load the data pointed to by data into the buffer. usage is passed along verbatim. targ is used for binding and it should most likely be ARRAY_BUFFER.
func (buf *Buffer) Set(targ int, data interface{}, usage int) {
	buf.Bind(targ)
	p, t, ts, s := toCtype(data)
	C.glBufferData(C.GLenum(targ), C.GLsizeiptr(s), p, C.GLenum(usage))
	buf.t = t
	buf.ts = ts
	buf.Unbind(targ)
}

func GetIntegerv(targ int, size int) (data []int) {
	data = make([]int, 4)
	var p []C.GLint = make([]C.GLint, size)
	C.glGetIntegerv(C.GLenum(targ), &p[0])
	for i := 0; i < size; i++ {
		data[i] = int(p[i])
	}
	return
}

// Bind calls glBindBuffer
func (buf *Buffer) Bind(targ int) {
	C.glBindBuffer(C.GLenum(targ), buf.i)
}

// Unbind calls glBindBuffer with a 0 argument
func (*Buffer) Unbind(targ int) {
	C.glBindBuffer(C.GLenum(targ), 0)
}

// The type Shader represents a shader.
type Shader C.GLuint

// NewShader creates a shader object of type typ, loads it with source code src and compiles it
func NewShader(typ int, src string) (Shader, error) {
	var val C.GLint
	shad := C.glCreateShader(C.GLenum(typ))
	s := (*C.GLchar)(C.CString(src))
	C.glShaderSource(shad, 1, &s, nil)
	C.glCompileShader(shad)
	C.glGetShaderiv(shad, COMPILE_STATUS, &val)
	if val != TRUE {
		C.glGetShaderiv(shad, INFO_LOG_LENGTH, &val)
		buf := make([]C.GLchar, val+1)
		C.glGetShaderInfoLog(shad, C.GLsizei(val), nil, &buf[0])
		C.glDeleteShader(shad)
		return Shader(0), errors.New(C.GoString((*C.char)(&buf[0])))
	}
	return Shader(shad), nil
}

// The type Program represents a shader program. It contains maps to cache the location of attributes and uniforms.
type Program struct {
	i    C.GLuint
	attr map[string]C.GLuint
	uni  map[string]C.GLint
}

// NewProgram creates an empty program
func NewProgram() *Program {
	return &Program{i: C.glCreateProgram()}
}

// Attach attaches a shader object
func (p *Program) Attach(s Shader) {
	C.glAttachShader(p.i, C.GLuint(s))
}

// Detach detaches a shader object
func (p *Program) Detach(s Shader) {
	C.glDetachShader(p.i, C.GLuint(s))
}

// Delete deletes the program object
func (p *Program) Delete() {
	C.glDeleteProgram(p.i)
}

// Use calls glUseProgram
func (p *Program) Use() {
	C.glUseProgram(p.i)
}

// Unuse calls glUseProgram with a 0 argument
func (p *Program) Unuse() {
	C.glUseProgram(C.GLuint(0))
}

// Link links the attached shader objects
func (p *Program) Link() error {
	var val, val2 C.GLint
	C.glLinkProgram(p.i)
	C.glGetProgramiv(p.i, LINK_STATUS, &val)
	if val != TRUE {
		C.glGetProgramiv(p.i, INFO_LOG_LENGTH, &val)
		buf := make([]C.GLchar, val+1)
		C.glGetProgramInfoLog(p.i, C.GLsizei(val), nil, &buf[0])
		return errors.New(C.GoString((*C.char)(&buf[0])))
	}
	p.attr = make(map[string]C.GLuint)
	C.glGetProgramiv(p.i, ACTIVE_ATTRIBUTES, &val)
	C.glGetProgramiv(p.i, ACTIVE_ATTRIBUTE_MAX_LENGTH, &val2)
	buf := make([]C.char, val2)
	for i := C.GLuint(0); i < C.GLuint(val); i++ {
		C.glGetActiveAttrib(p.i, i, C.GLsizei(val2), nil, nil, nil, (*C.GLchar)(&buf[0]))
		p.attr[C.GoString(&buf[0])] = C.GLuint(C.glGetAttribLocation(p.i, (*C.GLchar)(&buf[0])))
	}
	p.uni = make(map[string]C.GLint)
	C.glGetProgramiv(p.i, ACTIVE_UNIFORMS, &val)
	C.glGetProgramiv(p.i, ACTIVE_UNIFORM_MAX_LENGTH, &val2)
	buf = make([]C.char, val2)
	for i := C.GLuint(0); i < C.GLuint(val); i++ {
		C.glGetActiveUniform(p.i, i, C.GLsizei(val2), nil, nil, nil, (*C.GLchar)(&buf[0]))
		p.uni[C.GoString(&buf[0])] = C.glGetUniformLocation(p.i, (*C.GLchar)(&buf[0]))
	}
	return nil
}

// EnableAttrib calls glEnableVertexAttribArray and glVertexAttribPointer to activate an attribute and connect it to a buffer object.
// offset specifies the first vertex, stride specifies the distance from the beginning of one vertex to the next, size specifies the number of components in a vertex (all these arguments are in units of array elements, not bytes like the underlying API).
// The byte offset of component j of vertex i is thus calculated as: sizeof(data[0]) * (offset + stride * i + j), where data is the parameter passed to Buffer.Set
func (p *Program) EnableAttrib(loc string, buf *Buffer, offset int, size int, stride int, norm bool) {
	n := FALSE
	if norm {
		n = TRUE
	}
	buf.Bind(ARRAY_BUFFER)
	attr := p.attr[loc]
	C.glEnableVertexAttribArray(attr)
	C.glVertexAttribPointer(attr, C.GLint(size), buf.t, C.GLboolean(n), C.GLsizei(stride*buf.ts), unsafe.Pointer(uintptr(buf.ts*offset)))
	buf.Unbind(ARRAY_BUFFER)
}

// DisableAttrib calls glDisableVertexAttribArray
func (p *Program) DisableAttrib(loc string) {
	C.glDisableVertexAttribArray(p.attr[loc])
}

// SetUniform sets a uniform variable using the appropriate glUniform* or glUniformMatrix* call. It supports arrays of float32 and float64 or Mat4 objects.
// NB: The underlying API does not support double precision, being able to pass float64 values is for convenience only.
// BUG: It does not support non-square matrices.
func (p *Program) SetUniform(loc string, data interface{}) {
	uni := p.uni[loc]
	switch f := data.(type) {
	case float32:
		C.glUniform1f(uni, C.GLfloat(f))
	case [1]float32:
		C.glUniform1f(uni, C.GLfloat(f[0]))
	case [2]float32:
		C.glUniform2f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]))
	case [3]float32:
		C.glUniform3f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]))
	case [4]float32:
		C.glUniform4f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]), C.GLfloat(f[3]))
	case float64:
		C.glUniform1f(uni, C.GLfloat(f))
	case [1]float64:
		C.glUniform1f(uni, C.GLfloat(f[0]))
	case [2]float64:
		C.glUniform2f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]))
	case [3]float64:
		C.glUniform3f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]))
	case [4]float64:
		C.glUniform4f(uni, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]), C.GLfloat(f[3]))
	case int:
		C.glUniform1i(uni, C.GLint(f))
	case [1]int:
		C.glUniform1i(uni, C.GLint(f[0]))
	case [2]int:
		C.glUniform2i(uni, C.GLint(f[0]), C.GLint(f[1]))
	case [3]int:
		C.glUniform3i(uni, C.GLint(f[0]), C.GLint(f[1]), C.GLint(f[2]))
	case [4]int:
		C.glUniform4i(uni, C.GLint(f[0]), C.GLint(f[1]), C.GLint(f[2]), C.GLint(f[3]))
	case [2][2]float32:
		g := [4]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1])}
		C.glUniformMatrix2fv(uni, 1, FALSE, &g[0])
	case [2][2]float64:
		g := [4]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1])}
		C.glUniformMatrix2fv(uni, 1, FALSE, &g[0])
	case [3][3]float32:
		g := [9]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[2][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1]), C.GLfloat(f[2][1]), C.GLfloat(f[0][2]), C.GLfloat(f[1][2]), C.GLfloat(f[2][2])}
		C.glUniformMatrix3fv(uni, 1, FALSE, &g[0])
	case [3][3]float64:
		g := [9]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[2][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1]), C.GLfloat(f[2][1]), C.GLfloat(f[0][2]), C.GLfloat(f[1][2]), C.GLfloat(f[2][2])}
		C.glUniformMatrix3fv(uni, 1, FALSE, &g[0])
	case [4][4]float32:
		g := [16]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[2][0]), C.GLfloat(f[3][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1]), C.GLfloat(f[2][1]), C.GLfloat(f[3][1]), C.GLfloat(f[0][2]), C.GLfloat(f[1][2]), C.GLfloat(f[2][2]), C.GLfloat(f[3][2]), C.GLfloat(f[0][3]), C.GLfloat(f[1][3]), C.GLfloat(f[2][3]), C.GLfloat(f[3][3])}
		C.glUniformMatrix4fv(uni, 1, FALSE, &g[0])
	case [4][4]float64:
		g := [16]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[2][0]), C.GLfloat(f[3][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1]), C.GLfloat(f[2][1]), C.GLfloat(f[3][1]), C.GLfloat(f[0][2]), C.GLfloat(f[1][2]), C.GLfloat(f[2][2]), C.GLfloat(f[3][2]), C.GLfloat(f[0][3]), C.GLfloat(f[1][3]), C.GLfloat(f[2][3]), C.GLfloat(f[3][3])}
		C.glUniformMatrix4fv(uni, 1, FALSE, &g[0])
	case Mat4:
		g := [16]C.GLfloat{C.GLfloat(f[0][0]), C.GLfloat(f[1][0]), C.GLfloat(f[2][0]), C.GLfloat(f[3][0]), C.GLfloat(f[0][1]), C.GLfloat(f[1][1]), C.GLfloat(f[2][1]), C.GLfloat(f[3][1]), C.GLfloat(f[0][2]), C.GLfloat(f[1][2]), C.GLfloat(f[2][2]), C.GLfloat(f[3][2]), C.GLfloat(f[0][3]), C.GLfloat(f[1][3]), C.GLfloat(f[2][3]), C.GLfloat(f[3][3])}
		C.glUniformMatrix4fv(uni, 1, FALSE, &g[0])
	default:
		panic("invalid type passed to SetUniform()")
	}
}

// MakeProgram is a convenience routine which calls NewProgram(), NewShader(), Shader.Attach() and Program.Link() to create a shader program object.
func MakeProgram(vertex []string, fragment []string) (*Program, error) {
	p := NewProgram()
	for _, s := range vertex {
		shad, err := NewShader(VERTEX_SHADER, s)
		if err != nil {
			p.Delete()
			return nil, err
		}
		p.Attach(shad)
	}
	for _, s := range fragment {
		shad, err := NewShader(FRAGMENT_SHADER, s)
		if err != nil {
			p.Delete()
			return nil, err
		}
		p.Attach(shad)
	}
	err := p.Link()
	if err != nil {
		p.Delete()
		return nil, err
	}
	return p, nil
}

// DrawArrays calls glDrawArrays
func DrawArrays(mode, first, count int) {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}

// The type Texture represents a texture object.
type Texture C.GLuint

// NewTexture2D creates a new texture object from the given image using glTexImage2D. It uses RGBA as a color format and sets GL_TEXTURE_{MIN,MAG}_FILTER to GL_NEAREST
func NewTexture2D(img image.Image, border int) Texture {
	var t C.GLuint

	C.glGenTextures(1, &t)
	tt := Texture(t)
	tt.Bind(TEXTURE_2D)
	data := make([]uint16, img.Bounds().Dx()*img.Bounds().Dy()*4)
	r := img.Bounds()
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			R, G, B, A := img.At(x, y).RGBA()
			i := (y*r.Dx() + x) * 4
			data[i] = uint16(R)
			data[i+1] = uint16(G)
			data[i+2] = uint16(B)
			data[i+3] = uint16(A)
		}
	}
	C.glTexImage2D(TEXTURE_2D, 0, RGBA, C.GLsizei(img.Bounds().Dx()), C.GLsizei(img.Bounds().Dy()), C.GLint(border), RGBA, UNSIGNED_SHORT, unsafe.Pointer(&data[0]))
	C.glTexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, NEAREST)
	C.glTexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, NEAREST)
	tt.Unbind(TEXTURE_2D)
	return tt
}

// Bind calls glBindTexture
func (t Texture) Bind(targ int) {
	C.glBindTexture(C.GLenum(targ), C.GLuint(t))
}

// Unbind calls glBindTexture with a 0 argument
func (Texture) Unbind(targ int) {
	C.glBindTexture(C.GLenum(targ), 0)
}

// TexParameteri calls glTexParameteri on the texture. The targ argument is used for binding and should most likely be TEXTURE_2D.
func (t Texture) TexParameteri(targ, pname, param int) {
	t.Bind(targ)
	C.glTexParameteri(C.GLenum(targ), C.GLenum(pname), C.GLint(param))
	t.Unbind(targ)
}

// Enable calls glActiveTexture and Bind
func (t Texture) Enable(unit int, targ int) {
	C.glActiveTexture(TEXTURE0 + C.GLenum(unit))
	t.Bind(targ)
}

// Disable calls glActiveTexture and Unbind
func (t Texture) Disable(unit int, targ int) {
	C.glActiveTexture(TEXTURE0 + C.GLenum(unit))
	t.Unbind(targ)
}
