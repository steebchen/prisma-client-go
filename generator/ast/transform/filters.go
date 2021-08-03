package transform

// Method defines the method for the virtual types method
type Method struct {
	Name   string
	Action string
}

// Filter defines the data struct for the virtual types method
type Filter struct {
	// Scalar is the scalar name of a type, e.g. String, Int or DateTime
	Scalar  string
	Methods []Method
}
