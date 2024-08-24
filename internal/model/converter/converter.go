package converter

import (
	"github.com/Foxtrot1388/MessaggioTask/internal/entity"
	"github.com/Foxtrot1388/MessaggioTask/internal/model"
)

func GetOutputMessage(e *entity.OutputMessage) model.OutputMessage {

	return model.OutputMessage{
		ID: e.ID,
	}

}

func GetStatMessage(e *entity.StatMessage) model.StatMessage {

	return model.StatMessage{
		Count: e.Count,
		Day:   e.Day,
	}

}
