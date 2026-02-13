package infopane

type BucketInfoPane struct {
	Model
}

func NewBucket(label, value string) BucketInfoPane {
	return BucketInfoPane{Model: New(label, value)}
}

func (m *BucketInfoPane) Clear() {
	m.Model.Clear()
}
