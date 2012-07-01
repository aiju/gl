package gl

import "math"

const deg = math.Pi / 180

type Mat4 [4][4]float64

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

func RotZ(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{math.Cos(r), -math.Sin(r), 0, 0},
		[4]float64{math.Sin(r), math.Cos(r), 0, 0},
		[4]float64{0, 0, 1, 0},
		[4]float64{0, 0, 0, 1}}
}

func RotX(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{1, 0, 0, 0},
		[4]float64{0, math.Cos(r), math.Sin(r), 0},
		[4]float64{0, -math.Sin(r), math.Cos(r), 0},
		[4]float64{0, 0, 0, 1}}
}

func RotY(r float64) Mat4 {
	r *= deg
	return Mat4{[4]float64{math.Cos(r), 0, math.Sin(r), 0},
		[4]float64{0, 1, 0, 0},
		[4]float64{-math.Sin(r), 0, math.Cos(r), 0},
		[4]float64{0, 0, 0, 1}}
}

func Translate(x, y, z float64) Mat4 {
	return Mat4{[4]float64{1, 0, 0, x},
		[4]float64{0, 1, 0, y},
		[4]float64{0, 0, 1, z},
		[4]float64{0, 0, 0, 1}}
}

func Frustum(fov, aspect, zNear, zFar float64) Mat4 {
	f := 1 / math.Tan(fov*deg/2)
	return Mat4{[4]float64{f / aspect, 0, 0, 0},
		[4]float64{0, f, 0, 0},
		[4]float64{0, 0, (zNear + zFar) / (zNear - zFar), (2 * zNear * zFar) / (zNear - zFar)},
		[4]float64{0, 0, -1, 0}}
}

func Scale(x, y, z float64) Mat4 {
	return Mat4{[4]float64{x, 0, 0, 0},
		[4]float64{0, y, 0, 0},
		[4]float64{0, 0, z, 0},
		[4]float64{0, 0, 0, 1}}
}
