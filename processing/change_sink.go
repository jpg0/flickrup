package processing

type ChangeSink interface {
	Expect(change string)
}