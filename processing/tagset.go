package processing

type TagSet struct {
	tags map[string]struct{}
}

func NewEmptyTagSet() *TagSet {
	return &TagSet{
		make(map[string]struct{}),
	}
}

func NewTagSet(ss []string) *TagSet {
	rv := NewEmptyTagSet()

	for _, s := range ss {
		rv.Add(s)
	}

	return rv
}

func (t *TagSet) Add(s string) {
	t.tags[s] = struct {}{}
}

func  (t *TagSet) Remove(s string) {
	delete(t.tags, s)
}

func (t *TagSet) AddAll(other *TagSet) {
	for k := range other.tags {
		t.Add(k)
	}
}

func (t *TagSet) Slice() []string {
	keys := make([]string, len(t.tags))

	i := 0
	for k := range t.tags {
		keys[i] = k
		i++
	}

	return keys
}

func (t *TagSet) Size() int {
	return len(t.tags)
}

func (t *TagSet) Contains(s string) bool {
	_, contains := t.tags[s]
	return contains
}