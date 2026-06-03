package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type ItemData struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderEvent struct {
	CustomerName  string     `json:"customerName"`
	CustomerEmail string     `json:"customerEmail"`
	OrderNumber   string     `json:"orderNumber"`
	Items         []ItemData `json:"items"`
	Subtotal      float64    `json:"subtotal"`
	ShippingFee   float64    `json:"shippingFee"`
	VatAmount     float64    `json:"vatAmount"`
	TotalAmount   float64    `json:"totalAmount"`
}

type NotificationService struct {
	cfg MailtrapConfig
}

func NewNotificationService(cfg MailtrapConfig) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (svc *NotificationService) ProcessAndSendEmail(payload string) {
	var event OrderEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("❌ Error parsing incoming JSON payload: %v", err)
		return
	}

	log.Printf("📦 Processing order #%s for %s", event.OrderNumber, event.CustomerEmail)

	if svc.cfg.Username == "" || svc.cfg.Password == "" {
		log.Printf("⚠️ Missing Mailtrap credentials in configuration! Skipping SMTP relay for Order #%s", event.OrderNumber)
		return
	}

	htmlBody, err := svc.compileTemplate(event)
	if err != nil {
		log.Printf("❌ Template compilation aborted for Order #%s: %v", event.OrderNumber, err)
		return
	}

	auth := smtp.PlainAuth("", svc.cfg.Username, svc.cfg.Password, svc.cfg.Host)
	to := []string{event.CustomerEmail}

	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: Roofi Homelabs Store <%s>\r\n"+
			"Subject: Order #%s Confirmation Successfully Received!\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		event.CustomerEmail, svc.cfg.From, event.OrderNumber, htmlBody,
	))

	addr := fmt.Sprintf("%s:%s", svc.cfg.Host, svc.cfg.Port)
	if err := smtp.SendMail(addr, auth, svc.cfg.From, to, msg); err != nil {
		log.Printf("❌ Failed to relay SMTP message to Mailtrap for Order #%s: %v", event.OrderNumber, err)
		return
	}

	log.Printf("📧 Email successfully blasted to Mailtrap for Order #%s!", event.OrderNumber)
}

func (svc *NotificationService) compileTemplate(event OrderEvent) (string, error) {
	funcMap := template.FuncMap{
		"formatIDR": FormatIDR,
	}

	tmpl, err := template.New("order_confirmation_template.html").Funcs(funcMap).ParseFiles("order_confirmation_template.html")
	if err != nil {
		return "", fmt.Errorf("failed to load order_confirmation_template.html: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, event); err != nil {
		return "", fmt.Errorf("failed to map fields to template execution token: %w", err)
	}

	return buf.String(), nil
}

func FormatIDR(val float64) string {
	intVal := int64(val)
	str := fmt.Sprintf("%d", intVal)

	var result []byte
	length := len(str)

	for i, char := range str {
		result = append(result, byte(char))
		if (length-i-1)%3 == 0 && i != length-1 {
			result = append(result, '.')
		}
	}
	return string(result)
}
