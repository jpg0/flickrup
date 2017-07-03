package listen

import "github.com/Sirupsen/logrus"

type ChangeManger struct {
	expected *StringBag
}

func NewChangeManager() *ChangeManger {
	return &ChangeManger{
		expected: NewStringBag(),
	}
}

func (cm *ChangeManger) Expect(filepath string) {
	logrus.Debugf("Expecting %s", filepath)
	cm.expected.Add(filepath)
}

//returns true iff the change was expected
func (cm *ChangeManger) ChangeObserved(watchDir, filepath string) bool {

	if (IsStatusFile(watchDir, filepath)) {
		return false
	}

	removed := cm.expected.Remove(filepath)

	if (removed) {
		logrus.Debugf("Observed and ignored change: %s", filepath)
	} else {
		logrus.Debugf("Observed and registered change: %s", filepath)
	}

	return removed
}

//StringBag implementation

type StringBag struct {
	items map[string]int
}

func NewStringBag() *StringBag {
	return &StringBag{items: make(map[string]int)}
}

func (sb *StringBag) Add(s string) {
	count := sb.items[s]
	sb.items[s] = count + 1
}

func (sb *StringBag) Remove(s string) bool {
	count := sb.items[s]
	switch count {
	case 0: return false
	case 1: delete(sb.items, s); return true
	default: sb.items[s] = count - 1; return true
	}
}
