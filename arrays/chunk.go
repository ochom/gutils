package arrays

// Chunk will take a slice of any kind and a chunk size and return a slice of slices
func Chunk[T any](arr []T, chunkSize int) [][]T {
	chunks := [][]T{}
	if len(arr) == 0 {
		return chunks
	}

	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}

	return chunks
}
