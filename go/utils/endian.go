package utils

// FromLittleEndian converts a byte slice from little endian to big endian.
func FromLittleEndian(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return bytes
}
