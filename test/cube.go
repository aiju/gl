package main

import (
	"fmt"
	"github.com/neagix/Go-SDL/sdl"
	"github.com/aiju/gl"
	"image"
	_ "image/png"
	"os"
	"time"
)

var Vertices = []float64{
	-1, 1, -1, 0, 0,
	1, 1, -1, 1, 0,
	-1, -1, -1, 0, 1,
	1, -1, -1, 1, 1,
	1, -1, 1, 1, 0,
	1, 1, -1, 0, 1,
	1, 1, 1, 0, 0,
	-1, 1, -1, 1, 1,
	-1, 1, 1, 1, 0,
	-1, -1, -1, 0, 1,
	-1, -1, 1, 0, 0,
	1, -1, 1, 1, 0,
	-1, 1, 1, 0, 1,
	1, 1, 1, 1, 1,
}

var vertexShader = `
#version 110

attribute vec3 position;
uniform mat4 matrix;
attribute vec2 texcoord;
varying vec2 texco;

void main() {
	gl_Position = matrix * vec4(position.xyz, 1);
	texco = texcoord;
}
`

var fragmentShader = `
#version 110

uniform sampler2D tex;
varying vec2 texco;

void main() {
	gl_FragColor = texture2D(tex, texco);
}
`

func main() {
	sdl.Init(sdl.INIT_VIDEO)
	sdl.SetVideoMode(800, 600, 32, sdl.OPENGL|sdl.DOUBLEBUF|sdl.HWSURFACE)
	gl.Init()
	gl.Enable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, 800, 600)
	tick := time.Tick(time.Second / 50)
	timer := 0.0
	posbuf := gl.NewBuffer(gl.ARRAY_BUFFER, Vertices, gl.STATIC_DRAW)
	prog, err := gl.MakeProgram([]string{vertexShader}, []string{fragmentShader})
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.Open("glenda.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	tex := gl.NewTexture2D(img, 0)
	for {
		select {
		case ev := <-sdl.Events:
			if _, ok := ev.(sdl.QuitEvent); ok {
				return
			}
		case <-tick:
			gl.ClearColor(0, 0, 0, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			prog.Use()
			mat := gl.Mul4(gl.Frustum(45, 800./600, 0.01, 100), gl.Translate(0, 0, -8), gl.RotX(timer), gl.RotY(2*timer), gl.RotZ(3*timer))
			prog.EnableAttrib("position", posbuf, 0, 3, 5, false)
			prog.EnableAttrib("texcoord", posbuf, 3, 2, 5, false)
			prog.SetUniform("tex", 0)
			prog.SetUniform("matrix", mat)
			tex.Enable(0, gl.TEXTURE_2D)
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(Vertices)/5)
			prog.DisableAttrib("position")
			prog.Unuse()

			sdl.GL_SwapBuffers()
			timer += 1
		}
	}
}
