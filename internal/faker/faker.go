package faker

import (
	"fmt"
	"math/rand"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

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
	visitorID := RandomString("", 3)
	sessionID := RandomString("", 3)
	return Event{
		EventType: "visitor_event",
		EventID:   RandomString("evt_", 6),
		Timestamp: RandomTimestamp(),
		Source:    MaybeEmpty(RandomString("", 8)),
		VisitorData: map[string]interface{}{
			"visitor_id":  visitorID,
			"session_id":  sessionID,
			"page_url":    RandomURL(),
			"referrer":    RandomReferrer(),
			"utm_params":  RandomUTM(),
			"device_info": RandomDeviceInfo(),
			"behavior":    RandomBehavior(),
		},
		Data: map[string]interface{}{
			"phone":  MaybeEmpty(RandomString("", 10)),
			"email":  MaybeEmpty(RandomString("", 10) + "@example.com"),
			"cookie": MaybeEmpty(RandomString("cookie_", 8)),
		},
		Identifiers: map[string]interface{}{
			"cmec_visitor_id":           visitorID,
			"cmec_contact_call_id":      MaybeEmpty(RandomString("call_", 5)),
			"cmec_contact_chat_id":      MaybeEmpty(RandomString("chat_", 5)),
			"cmec_contact_external_id":  MaybeEmpty(RandomString("ext_", 5)),
			"cmec_contact_form2lead_id": MaybeEmpty(RandomString("f2l_", 5)),
			"cmec_contact_tickets_id":   MaybeEmpty(RandomString("ticket_", 5)),
		},
	}
}

func RandomCustomer() Customer {
	visitorIDs := []string{RandomString("", 3), RandomString("", 3)}
	sessionIDs := []string{RandomString("", 3), RandomString("", 3)}
	email := RandomString("", 8) + "@example.com"
	phone := RandomString("", 10)
	return Customer{
		CustomerID: RandomString("cust_unified_", 5),
		CreatedAt:  RandomTimestamp(),
		UpdatedAt:  RandomTimestamp(),
		Merged:     RandomInt(0, 1),
		Deleted:    RandomInt(0, 1),
		Identifiers: map[string]interface{}{
			"email":           []string{email},
			"phone":           []string{phone},
			"cmec_visitor_id": visitorIDs,
			"cmec_session_id": sessionIDs,
		},
		PrimaryIdentifiers: map[string]interface{}{
			"email":           email,
			"phone":           phone,
			"cmec_visitor_id": visitorIDs[0],
		},
		PersonalData: map[string]interface{}{
			"name":              MaybeEmpty(RandomString("", 8)),
			"company":           MaybeEmpty(RandomString("", 8)),
			"title":             MaybeEmpty(RandomString("", 8)),
			"inferred_location": MaybeEmpty(RandomString("", 8)),
		},
		ConfidenceScore: RandomFloat(0.7, 1.0),
	}
}

// Helper for random int
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
	return fmt.Sprintf("https://%s.com/%s", RandomString("site", 5), RandomString("page", 3))
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
		"utm_source":   MaybeEmpty(RandomString("src", 4)),
		"utm_medium":   MaybeEmpty(RandomString("med", 4)),
		"utm_campaign": MaybeEmpty(RandomString("camp", 4)),
	}
}

func RandomDeviceInfo() map[string]string {
	if rand.Float64() < 0.2 {
		return map[string]string{}
	}
	return map[string]string{
		"user_agent":  MaybeEmpty(RandomString("ua", 8)),
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
