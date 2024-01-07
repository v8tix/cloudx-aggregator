package dto

type GroupDTO struct {
	AssociationsDTO *AssociationsDTO
	MessageDTO      *MessageDTO
}

func NewGroup(associations *AssociationsDTO, message *MessageDTO) *GroupDTO {
	return &GroupDTO{AssociationsDTO: associations, MessageDTO: message}
}
