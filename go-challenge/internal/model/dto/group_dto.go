package dto

type Group struct {
	Associations *AssociationsDTO
	Message      *MessageDTO
}

func NewGroup(associations *AssociationsDTO, message *MessageDTO) *Group {
	return &Group{Associations: associations, Message: message}
}
