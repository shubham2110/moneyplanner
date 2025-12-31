package persons

type PersonCreationRequest struct {
	PersonName string `json:"person_name"`
	Alias      string `json:"alias,omitempty"`
}

type PersonUpdateRequest struct {
	PersonName *string `json:"person_name,omitempty"`
	Alias      *string `json:"alias,omitempty"`
}
