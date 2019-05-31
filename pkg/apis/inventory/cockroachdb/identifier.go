package cockroachdb

type identifier struct {
	guid string
}

func (i *identifier) GUID() string {
	return i.guid
}
