package ippanel

import (
	"testing"
)

// TestSend test send sms
func TestSend(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	MessageId, err := sms.Send("+983000505", []string{"989153621841"}, "ippanel is awesome", "Summary")
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(MessageId)
}

// TestErrors test api errors handling
func TestErrors(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	_, err := sms.Send("983000505", []string{"989153621841"}, "ippanel is awesome", "Summary")
	if err != nil {
		if e, ok := err.(Error); ok {
			switch e.Code {
			case ErrUnprocessableEntity:
				fieldErrors := e.Message.(FieldErrs)
				for field, fieldError := range fieldErrors {
					t.Log(field, fieldError)
				}
			default:
				errMsg := e.Message.(string)
				t.Log(errMsg)
			}
		}
	}
}

// TestGetMessage tests getMessage method
func TestGetMessage(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	// 391245018
	message, err := sms.GetMessage(391245018)
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(message)
}

// TestFetchStatuses test fetch message recipients status
func TestFetchStatuses(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	// 391245018
	statuses, paginationInfo, err := sms.FetchStatuses(391245018, ListParams{Page: 0, Limit: 10})
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(statuses, paginationInfo)
}

// TestFetchInbox fetch inbox test
func TestFetchInbox(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	messages, paginationInfo, err := sms.FetchInbox(ListParams{Page: 1, Limit: 10})
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(messages, paginationInfo)
}

// TestCreatePattern test create pattern
func TestCreatePattern(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	pattern, err := sms.CreatePattern("%name% is awesome, your code is %%code%", "description", map[string]string{
		"name": "string", "code": "string",
	}, "%", true)
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(pattern)
}

// TestSendPattern test send with pattern
func TestSendPattern(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")
	patternValues := map[string]string{
		"name": "IPPANEL",
		"code": "code",
	}

	MessageId, err := sms.SendPattern(
		"tvixsy9ed89ukfu", // pattern code
		"+983000505",      // originator
		"989153621841",    // recipient
		patternValues,     // pattern values
	)
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(MessageId)
}

func TestGetCredit(t *testing.T) {
	sms := New("6Aujg2MOwhOlIy3h_PeXJLtCl8TudIwEgv1mZMW_WMA=")

	credit, err := sms.GetCredit()
	if err != nil {
		t.Error("error occurred ", err)
	}

	t.Log(credit)
}
