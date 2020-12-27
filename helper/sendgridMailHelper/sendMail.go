package sendgridMailHelper

import (
	"encoding/json"
	"errors"
	"flowban/helper/httpRequest"
)

const (
	SendgridEmailSendURL = "https://api.sendgrid.com/v3/mail/send"
)

type sendGridInternalMailStructure struct {
	Personalization []sendGridPersonalization `json:"personalizations"`

	Content []sendGridContent `json:"content"`

	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
}

type sendGridReceiver struct {
	Email string `json:"email"`
}

type sendGridPersonalization struct {
	To      []sendGridReceiver `json:"to"`
	Subject string             `json:"subject"`
}

type sendGridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (actualStruct sendGridInternalMailStructure) ToInternalStruct(emailData Email, senderEmail string) sendGridInternalMailStructure {
	var internalStruct sendGridInternalMailStructure
	var sendGridRecs []sendGridReceiver
	for _, subject := range emailData.Receivers {
		sendGridRecs = append(sendGridRecs, sendGridReceiver{Email: subject})
	}

	internalStruct.Personalization = append(internalStruct.Personalization, sendGridPersonalization{
		To:      sendGridRecs,
		Subject: emailData.Subject,
	})
	internalStruct.Content = append(internalStruct.Content, sendGridContent{
		Type:  emailData.ContentType,
		Value: emailData.Content,
	})
	internalStruct.From.Email = senderEmail
	return internalStruct
}

type Email struct {
	Receivers   []string `json:"receivers"`
	Subject     string   `json:"subject"`
	ContentType string   `json:"content_type"`
	Content     string   `json:"content"`
}

type SendGridMailService struct {
	apiKey string
}

func NewEmailService(apiKey string) SendGridMailService {
	return SendGridMailService{apiKey: apiKey}
}

func (mail *SendGridMailService) SendMail(email Email, senderEmail string) error {
	sendgridStructActual := sendGridInternalMailStructure{}.ToInternalStruct(email, senderEmail)

	//parse to json and send to html
	jsonPayload, err := json.Marshal(sendgridStructActual)
	if err != nil {
		return err
	}
	headers := httpRequest.HttpHeaders{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + mail.apiKey,
	}
	code, respBody, err := httpRequest.PostData(SendgridEmailSendURL, jsonPayload, headers, 120)
	if err != nil {
		return err
	}

	if code > 299 {
		return errors.New(string(respBody))
	}
	return nil
}
