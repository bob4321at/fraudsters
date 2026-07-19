package shader

var BorderShader = `//kage:unit pixels
	package main

	func Fragment(targetCoords vec4, srcPos vec2, _ vec4) vec4 {
		col := imageSrc0At(srcPos.xy)
		if col.w == 0 {
			leftpix := imageSrc0At(vec2(srcPos.x-1, srcPos.y))
			if leftpix.w != 0 {
				return vec4(0,0,0,1)
			}
			rightpix := imageSrc0At(vec2(srcPos.x+1, srcPos.y))
			if rightpix.w != 0 {
				return vec4(0,0,0,1)
			}
			uppix := imageSrc0At(vec2(srcPos.x, srcPos.y+1))
			if uppix.w != 0 {
				return vec4(0,0,0,1)
			}
			downpix := imageSrc0At(vec2(srcPos.x, srcPos.y-1))
			if downpix.w != 0 {
				return vec4(0,0,0,1)
			}
		}

		return col
	}
`
