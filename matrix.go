package gl

import "math"

const deg = math.Pi / 180

// The type Mat4 represents a double precision 4x4 matrix.
type Mat4 [4][4]float64

// An identity matrix
var Identity Mat4 = [4][4]float64{[4]float64{1, 0, 0, 0}, [4]float64{0, 1, 0, 0}, [4]float64{0, 0, 1, 0}, [4]float64{0, 0, 0, 1}}

func mul4(a, b Mat4) Mat4 {
	var r Mat4

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				r[i][k] += a[i][j] * b[j][k]
			}
		}
	}
	return r
}

// Mul4 multiplies an arbitrary number of Mat4 matrices.
func Mul4(a ...Mat4) Mat4 {
	r := Identity
	for i, j := range a {
		if i == 0 {
			r = j
		} else {
			r = mul4(r, j)
		}
	}
	return r
}

// RotZ returns a rotation matrix rotating r degrees around the z axis.
func RotZ(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{math.Cos(r), -math.Sin(r), 0, 0},
		[4]float64{math.Sin(r), math.Cos(r), 0, 0},
		[4]float64{0, 0, 1, 0},
		[4]float64{0, 0, 0, 1}}
}

// RotX returns a rotation matrix rotating r degrees around the x axis.
func RotX(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{1, 0, 0, 0},
		[4]float64{0, math.Cos(r), math.Sin(r), 0},
		[4]float64{0, -math.Sin(r), math.Cos(r), 0},
		[4]float64{0, 0, 0, 1}}
}

// RotY returns a rotation matrix rotating r degrees around the y axis.
func RotY(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{math.Cos(r), 0, math.Sin(r), 0},
		[4]float64{0, 1, 0, 0},
		[4]float64{-math.Sin(r), 0, math.Cos(r), 0},
		[4]float64{0, 0, 0, 1}}
}

// Translate returns a translation matrix
func Translate(x, y, z float64) Mat4 {
	return Mat4{[4]float64{1, 0, 0, x},
		[4]float64{0, 1, 0, y},
		[4]float64{0, 0, 1, z},
		[4]float64{0, 0, 0, 1}}
}

// Frustum returns a projection matrix similar to gluPerspective. The arguments are field of view angle in degrees, aspect ratio and near and far z clipping plane distance.
func Frustum(fov, aspect, zNear, zFar float64) Mat4 {
	f := 1 / math.Tan(fov*deg/2)
	return Mat4{[4]float64{f / aspect, 0, 0, 0},
		[4]float64{0, f, 0, 0},
		[4]float64{0, 0, (zNear + zFar) / (zNear - zFar), (2 * zNear * zFar) / (zNear - zFar)},
		[4]float64{0, 0, -1, 0}}
}

// Scale returns a scale matrix
func Scale(x, y, z float64) Mat4 {
	return Mat4{[4]float64{x, 0, 0, 0},
		[4]float64{0, y, 0, 0},
		[4]float64{0, 0, z, 0},
		[4]float64{0, 0, 0, 1}}
}
