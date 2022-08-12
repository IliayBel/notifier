package zulip

type MessagesService service

type MessageResponse struct {
	Msg    string `json:"msg"`
	Result string `json:"result"`
	ID     int64  `json:"id"`
}

type MessageRequest struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Subject string `json:"subject,omitempty"`
	Content string `json:"content"`
}

func (s *MessagesService) Send(type_ string, to string, subject string, content string) (_ int64, err error) {
	body := &MessageRequest{
		Type:    type_,
		To:      to,
		Subject: subject,
		Content: content,
	}
	resp := new(MessageResponse)

	_, err = s.client.sling.New().Post("messages").BodyJSON(body).Receive(resp, nil)
	if err != nil {
		return 0, err
	}

	return resp.ID, err
}
