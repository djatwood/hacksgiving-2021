package generator

func GetColor(value float32, simple bool) [4]uint8 {
	// water
	if value < 0.3 {
		return [4]uint8{127, 205, 255, 255}
	}

	// beach
	if value < 0.35 {
		return [4]uint8{210, 255, 228, 255}
	}

	if simple {
		return [4]uint8{2, 161, 78, 255}
	}

	// grass
	if value < 0.7 {
		return [4]uint8{2, 161, 78, 255}
	}

	// forest
	if value < 0.9 {
		return [4]uint8{1, 85, 41, 255}
	}

	// deep forest
	return [4]uint8{60, 42, 42, 255}
}
