package response

const StatusType = "statuses"

type Status struct {
	Data struct {
		Identification
		Attributes struct {
			Name string `json:"name"`
		} `json:"attributes"`
	} `json:"data"`
}

func NewStatus() *Status {
	return &Status{
		Data: struct {
			Identification
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		}{
			Identification: Identification{
				ID:   "1",
				Type: StatusType,
			},
			Attributes: struct {
				Name string `json:"name"`
			}{
				Name: "success",
			},
		},
	}
}
