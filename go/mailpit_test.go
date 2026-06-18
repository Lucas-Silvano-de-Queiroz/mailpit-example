package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"testing"
	"time"
)

/**
 * Example of integration tests with Mailpit in Go.
 * Uses only the standard library:
 * - testing: test runner
 * - net/http: client for Mailpit API
 * - net/smtp: client for sending emails
 * - encoding/json: to process API responses
 */

const mailpitAPI = "http://localhost:8025/api/v1"

// Structs to map the Mailpit API response
type MailpitMessages struct {
	Total    int `json:"total"`
	Messages []struct {
		ID      string `json:"ID"`
		Subject string `json:"Subject"`
		From    struct {
			Address string `json:"Address"`
		} `json:"From"`
		To []struct {
			Address string `json:"Address"`
		} `json:"To"`
	} `json:"messages"`
}

type MailpitMessageDetail struct {
	ID   string `json:"ID"`
	HTML string `json:"HTML"`
	Text string `json:"Text"`
}

func TestMailpitIntegration(t *testing.T) {
	// 1. Clear existing emails
	t.Run("Initial Cleanup", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, mailpitAPI+"/messages", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error cleaning messages: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 on delete, got %d", resp.StatusCode)
		}
	})

	// 2. Send email via SMTP
	t.Run("Email Sending", func(t *testing.T) {
		from := "go-test@example.com"
		to := []string{"go-recipient@example.com"}
		subject := "Go API Test 🐹"
		body := "This email was sent during the Go integration test."

		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n<b>%s</b>",
			from, to[0], subject, body)

		err := smtp.SendMail("localhost:1025", nil, from, to, []byte(msg))
		if err != nil {
			t.Fatalf("Error sending email via SMTP: %v", err)
		}
	})

	// Brief pause for Mailpit indexing
	time.Sleep(300 * time.Millisecond)

	var lastMessageID string

	// 3. Validate if the email was received via API
	t.Run("Validate Receipt", func(t *testing.T) {
		resp, err := http.Get(mailpitAPI + "/messages")
		if err != nil {
			t.Fatalf("Error fetching messages: %v", err)
		}
		defer resp.Body.Close()

		var result MailpitMessages
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Error decoding JSON: %v", err)
		}

		if result.Total != 1 {
			t.Errorf("Expected 1 message, found %d", result.Total)
		}

		msg := result.Messages[0]
		if msg.Subject != "Go API Test 🐹" {
			t.Errorf("Incorrect subject: %s", msg.Subject)
		}

		if msg.To[0].Address != "go-recipient@example.com" {
			t.Errorf("Incorrect recipient: %s", msg.To[0].Address)
		}

		lastMessageID = msg.ID
		fmt.Printf("✅ Email found! ID: %s\n", lastMessageID)
	})

	// 4. Validate detailed content
	t.Run("Validate Content", func(t *testing.T) {
		if lastMessageID == "" {
			t.Skip("Skipping as we don't have the message ID")
		}

		resp, err := http.Get(mailpitAPI + "/message/" + lastMessageID)
		if err != nil {
			t.Fatalf("Error fetching message details: %v", err)
		}
		defer resp.Body.Close()

		var detail MailpitMessageDetail
		if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
			t.Fatalf("Error decoding details: %v", err)
		}

		if !strings.Contains(detail.HTML, "sent during the Go integration test") {
			t.Errorf("HTML content does not contain the expected text")
		}
		fmt.Println("✅ HTML content validated successfully.")
	})
}
