package ippanel

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageConfirmState message confirm state
type MessageConfirmState string

const (
	// MessageConfirmeStatePending pending
	MessageConfirmeStatePending MessageConfirmState = "notconfirm"
	// MessageConfirmeStateConfirmed confirmed
	MessageConfirmeStateConfirmed MessageConfirmState = "approve"
	// MessageConfirmeStateRejected rejected
	MessageConfirmeStateRejected MessageConfirmState = "reject"
)

// PatternStatus ...
type PatternStatus string

const (
	// PatternStatusActive active
	PatternStatusActive PatternStatus = "active"
	// PatternStatusInactive inactive
	PatternStatusInactive PatternStatus = "inactive"
	// PatternStatusPending pending
	PatternStatusPending PatternStatus = "pending"
)

// Message message model
type Message struct {
	MessageId      int64               `json:"message_id"`
	Number         string              `json:"number"`
	Message        string              `json:"message"`
	State          string              `json:"state"`
	Type           string              `json:"type"`
	Valid          MessageConfirmState `json:"valid"`
	Time           time.Time           `json:"time"`
	TimeSend       time.Time           `json:"time_sent"`
	RecipientCount int64               `json:"recipient_count"`
	ExitCount      int64               `json:"exit_count"`
	Part           int64               `json:"part"`
	Cost           float64             `json:"cost"`
	ReturnCost     float64             `json:"return_cost"`
	Summary        string              `json:"summary"`
}

// MessageRecipient message recipient status
type MessageRecipient struct {
	Recipient string `json:"recipient"`
	Status    int    `json:"status"`
}

// InboxMessage inbox message
type InboxMessage struct {
	To        string    `json:"to"`
	Message   string    `json:"message"`
	From      string    `json:"from"`
	CreatedAt time.Time `json:"created_at"`
	Type      string    `json:"type"`
}

// Pattern pattern
type Pattern struct {
	Code    string        `json:"code"`
	Status  PatternStatus `json:"status"`
	Message string        `json:"message"`
	IsShare bool          `json:"is_share"`
}

// sendSMSReqType request type for send sms
type sendSMSReqType struct {
	Sender      string            `json:"sender"`
	Recipient   []string          `json:"recipient"`
	Message     string            `json:"message"`
	Description map[string]string `json:"description"`
}

// sendResType response type for send sms
type sendResType struct {
	MessageId int64 `json:"message_id"`
}

// fetchMessageStatusesResType get message statuses response template
type fetchMessageStatusesResType struct {
	Statuses []MessageRecipient `json:"deliveries"`
}

type PatternVariable struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// createPatternReqType create pattern request type
type createPatternReqType struct {
	Pattern     string            `json:"message"`
	Delimiter   string            `json:"delimiter"`
	Description string            `json:"description"`
	Variable    []PatternVariable `json:"variable"`
	IsShare     bool              `json:"is_share"`
}

type PatternResType struct {
	Code string `json:"code"`
}

// sendPatternReqType send sms with pattern request template
type sendPatternReqType struct {
	Code      string            `json:"code"`
	Sender    string            `json:"sender"`
	Recipient string            `json:"recipient"`
	Variable  map[string]string `json:"variable"`
}

// Send send a message
func (sms *Ippanel) Send(sender string, recipients []string, message string, summary string) (int64, error) {
	data := sendSMSReqType{
		Sender:    sender,
		Recipient: recipients,
		Message:   message,
		Description: map[string]string{
			"summary":         summary,
			"count_recipient": fmt.Sprint(len(recipients)),
		},
	}

	_res, err := sms.post("/sms/send/panel/single", "application/json", data)
	if err != nil {
		return 0, err
	}

	res := sendResType{}
	if err = json.Unmarshal(_res.Data, &res); err != nil {
		return 0, err
	}

	return res.MessageId, nil
}

// GetMessage get a message by message_id
func (sms *Ippanel) GetMessage(MessageId int64) (*Message, error) {
	_res, err := sms.get("/sms/message/all", map[string]string{"message_id": fmt.Sprint(MessageId)})
	if err != nil {
		return nil, err
	}

	res := []Message{}
	if err = json.Unmarshal(_res.Data, &res); err != nil {
		return nil, err
	}

	return &res[0], nil
}

// FetchStatuses get message recipients statuses
func (sms *Ippanel) FetchStatuses(MessageId int64, pp ListParams) ([]MessageRecipient, *PaginationInfo, error) {
	_res, err := sms.get(fmt.Sprintf("/sms/message/show-recipient/message-id/%d", MessageId), map[string]string{
		"page":     fmt.Sprintf("%d", pp.Page),
		"per_page": fmt.Sprintf("%d", pp.Limit),
	})
	if err != nil {
		return nil, nil, err
	}

	res := &fetchMessageStatusesResType{}
	if err = json.Unmarshal(_res.Data, res); err != nil {
		return nil, nil, err
	}

	return res.Statuses, _res.Meta, nil
}

// FetchInbox fetch inbox messages list
func (sms *Ippanel) FetchInbox(pp ListParams) ([]InboxMessage, *PaginationInfo, error) {
	_res, err := sms.get("/inbox", map[string]string{
		"page":     fmt.Sprintf("%d", pp.Page),
		"per_page": fmt.Sprintf("%d", pp.Limit),
	})
	if err != nil {
		return nil, nil, err
	}

	res := []InboxMessage{}
	if err = json.Unmarshal(_res.Data, &res); err != nil {
		return nil, nil, err
	}

	return res, _res.Meta, nil
}

// CreatePattern create new pattern
func (sms *Ippanel) CreatePattern(pattern string, description string, variables map[string]string, delimiter string, isShare bool) (string, error) {
	data := createPatternReqType{
		Pattern:     pattern,
		Description: description,
		Variable:    []PatternVariable{},
		Delimiter:   delimiter,
		IsShare:     isShare,
	}

	for k, v := range variables {
		data.Variable = append(data.Variable, PatternVariable{Name: k, Type: v})
	}

	_res, err := sms.post("/sms/pattern/normal/store", "application/json", data)
	if err != nil {
		return "", err
	}

	res := []Pattern{}
	if err = json.Unmarshal(_res.Data, &res); err != nil {
		return "", err
	}

	return res[0].Code, nil
}

// SendPattern send a message with pattern
func (sms *Ippanel) SendPattern(patternCode string, originator string, recipient string, values map[string]string) (int64, error) {
	data := sendPatternReqType{
		Code:      patternCode,
		Sender:    originator,
		Recipient: recipient,
		Variable:  values,
	}

	_res, err := sms.post("/sms/pattern/normal/send", "application/json", data)
	if err != nil {
		return 0, err
	}

	res := sendResType{}
	if err = json.Unmarshal(_res.Data, &res); err != nil {
		return 0, err
	}

	return res.MessageId, nil
}
