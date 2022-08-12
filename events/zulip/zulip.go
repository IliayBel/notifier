package zulip

type EventsService service

type EventResponse struct {
	Msg         string `json:"msg"`
	Result      string `json:"result"`
	QueueID     string `json:"queue_id"`
	LastEventID int64  `json:"last_event_id"`
}

type EventRequest struct {
	EventTypes    []string `json:"type"`
	ApplyMarkdown bool     `json:"apply_markdown"`
}

func (s *EventsService) Register(types []string, applyMarkdown bool) (_ string, _ int64, err error) {
	body := &EventRequest{
		EventTypes:    types,
		ApplyMarkdown: applyMarkdown,
	}

	resp := new(EventResponse)
	_, err = s.client.sling.New().Post("register").BodyJSON(body).Receive(resp, nil)
	if err != nil {
		return "", 0, err
	}

	return resp.QueueID, resp.LastEventID, err
}
