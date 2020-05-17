package carte

type Jsonable interface {
	Json() (name string, value string)
}

type Jable struct {
	Name  string
	Value string
}

func (j *Jable) Json() (string, string) {
	return j.Name, j.Value
}
