package response

type Identification struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func NewIdentification(id, t string) *Identification {
	return &Identification{ID: id, Type: t}
}
