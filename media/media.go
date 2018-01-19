package media

// Size is a struct, used for `GetSizes` method, it will return a slice of Size, media library will crop images automatically based on it
type Size struct {
	Width  int
	Height int
}