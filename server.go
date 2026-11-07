package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
)

type WaitlistRequest struct {
	Email string `json:"email"`
}

func handleWaitlist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WaitlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	err := sendWelcomeEmail(req.Email)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Added to SiGiL Enterprise waitlist"}`))
}

func sendWelcomeEmail(recipient string) error {
	from := "noreply@sigil-framework.space"
	password := "your-smtp-password"
	smtpHost := "smtp.mailgun.org"
	smtpPort := "587"

	msg := []byte("To: " + recipient + "\r\n" +
		"Subject: You're on the SiGiL Enterprise Waitlist\r\n\r\n" +
		"Thanks for requesting beta access to SiGiL Enterprise. We'll ping you as soon as slots open up.\n\n- KadeHood")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, msg)
}

func main() {
	http.HandleFunc("/api/waitlist", handleWaitlist)
	fmt.Println("Webhook listener active on port 8080...")
	http.ListenAndServe(":8080", nil)
}
