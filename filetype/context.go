package filetype

type ProcessingContext struct {
	Visibilty string
}

func NewProcessingContext() *ProcessingContext {
	return &ProcessingContext{
		Visibilty: "public",
	}
}