package renderer

type NetworkError struct{}

func (m *NetworkError) Error() string {
	return "network error"
}
