package faker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())
}

// Event structure
type Event struct {
	EventType   string                 `json:"event_type"`
	EventID     string                 `json:"event_id"`
	Timestamp   string                 `json:"timestamp"`
	Source      string                 `json:"source"`
	VisitorData map[string]interface{} `json:"visitor_data"`
	Data        map[string]interface{} `json:"data"`
	Identifiers map[string]interface{} `json:"identifiers"`
}

// Customer structure
type Customer struct {
	CustomerID         string                 `json:"customerId"`
	CreatedAt          string                 `json:"createdAt"`
	UpdatedAt          string                 `json:"updatedAt"`
	Merged             int                    `json:"merged"`
	Deleted            int                    `json:"deleted"`
	Identifiers        map[string]interface{} `json:"identifiers"`
	PrimaryIdentifiers map[string]interface{} `json:"primaryIdentifiers"`
	PersonalData       map[string]interface{} `json:"personalData"`
	ConfidenceScore    float64                `json:"confidenceScore"`
}

// Exported functions for random data generation
func RandomEvent() Event {
	visitorID := gofakeit.LetterN(3)
	sessionID := gofakeit.LetterN(3)
	return Event{
		EventType: "visitor_event",
		EventID:   "evt_" + gofakeit.LetterN(6),
		Timestamp: RandomTimestamp(),
		Source:    gofakeit.LetterN(8),
		VisitorData: map[string]interface{}{
			"behavior": map[string]interface{}{
				"interactions": []string{"scroll", "click_cta", "hover", "form_submit"}[:rand.Intn(4)+1],
				"pages_viewed": rand.Intn(10) + 1,
				"time_on_site": rand.Intn(591) + 10,
			},
			"device_info": map[string]interface{}{
				"device_type": []string{"desktop", "mobile", "tablet"}[rand.Intn(3)],
				"ip_address":  fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
				"user_agent":  gofakeit.LetterN(10),
			},
			"page_url":   RandomURL(),
			"referrer":   RandomReferrer(),
			"session_id": sessionID,
			"utm_params": map[string]interface{}{
				"utm_campaign": "camp" + gofakeit.LetterN(5),
				"utm_medium":   "med" + gofakeit.LetterN(5),
				"utm_source":   "src" + gofakeit.LetterN(4),
			},
			"visitor_id": visitorID,
		},
		Data: map[string]interface{}{
			"cookie": "cookie_" + gofakeit.LetterN(8),
			"email":  gofakeit.LetterN(10) + "@example.com",
			"phone":  gofakeit.LetterN(10),
		},
		Identifiers: map[string]interface{}{
			"cmec_contact_call_id":      "call_" + gofakeit.LetterN(5),
			"cmec_contact_chat_id":      "chat_" + gofakeit.LetterN(5),
			"cmec_contact_external_id":  "ext_" + gofakeit.LetterN(5),
			"cmec_contact_form2lead_id": "f2l_" + gofakeit.LetterN(5),
			"cmec_contact_tickets_id":   "ticket_" + gofakeit.LetterN(5),
			"cmec_visitor_id":           visitorID,
		},
	}
}

func RandomCustomer() Customer {
	visitorIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
	sessionIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
	email := gofakeit.Email()
	phone := gofakeit.Phone()
	return Customer{
		CustomerID: gofakeit.UUID(),
		CreatedAt:  gofakeit.Date().Format(time.RFC3339),
		UpdatedAt:  gofakeit.Date().Format(time.RFC3339),
		Merged:     gofakeit.Number(0, 1),
		Deleted:    gofakeit.Number(0, 1),
		Identifiers: map[string]interface{}{
			"visitor_ids": visitorIDs,
			"session_ids": sessionIDs,
		},
		PrimaryIdentifiers: map[string]interface{}{
			"email": email,
			"phone": phone,
		},
		PersonalData: map[string]interface{}{
			"name":              gofakeit.Name(),
			"company":           gofakeit.Company(),
			"title":             gofakeit.JobTitle(),
			"inferred_location": gofakeit.City() + ", " + gofakeit.Country(),
		},
		ConfidenceScore: gofakeit.Float64Range(0.6, 1.0),
	}
}
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// Helper for random float
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func MaybeEmpty(value string) string {
	if rand.Float64() > 0.3 {
		return value
	}
	return ""
}

func RandomString(prefix string, length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + string(b)
}

func RandomURL() string {
	return fmt.Sprintf("https://%s.com/%s", RandomString("site", 5), RandomString("page", 10))
}

func RandomReferrer() string {
	referrers := []string{
		RandomURL(),
		"/internal/path",
		"",
	}
	return referrers[rand.Intn(len(referrers))]
}

func RandomUTM() map[string]string {
	if rand.Float64() < 0.4 {
		return map[string]string{}
	}
	return map[string]string{
		"utm_source":   MaybeEmpty(RandomString("src", 10)),
		"utm_medium":   MaybeEmpty(RandomString("med", 10)),
		"utm_campaign": MaybeEmpty(RandomString("camp", 10)),
	}
}

func RandomDeviceInfo() map[string]string {
	if rand.Float64() < 0.2 {
		return map[string]string{}
	}
	return map[string]string{
		"user_agent":  MaybeEmpty(RandomString("ua", 20)),
		"ip_address":  MaybeEmpty(fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))),
		"device_type": MaybeEmpty([]string{"desktop", "mobile", "tablet", ""}[rand.Intn(4)]),
	}
}

func RandomBehavior() map[string]interface{} {
	if rand.Float64() < 0.2 {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"pages_viewed": rand.Intn(10) + 1,
		"time_on_site": rand.Intn(591) + 10,
		"interactions": []string{"scroll", "click_cta", "hover", "form_submit", "video_play"}[:rand.Intn(4)+1],
	}
}

func RandomTimestamp() string {
	now := time.Now().UTC()
	delta := time.Duration(rand.Intn(10000)-10000) * time.Minute
	return now.Add(delta).Format("2006-01-02T15:04:05Z")
}
