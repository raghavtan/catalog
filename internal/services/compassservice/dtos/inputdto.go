package dtos

type ValidationFunc func() error

type InputDTO struct {
	PreValidationFunc ValidationFunc
}

func (dto *InputDTO) GetPreValidationFunc() ValidationFunc {
	if dto.PreValidationFunc != nil {
		return dto.PreValidationFunc
	}
	return func() error { return nil }
}
