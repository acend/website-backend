package mailchimp

import (
	"backend/internal/config"
	"crypto/md5"
	"fmt"
	"github.com/hanzoai/gochimp3"
	"log"
	"net/http"
	"strings"
	"time"
)

func Subscribe(email string) error {
	apiKey := config.Values.MailchimpAPIKey

	client := gochimp3.New(apiKey)
	client.Timeout = 5 * time.Second

	listID := "d19c0edee5"

	list, err := client.GetList(listID, nil)
	if err != nil {
		return fmt.Errorf("get mailchimp list %s: %w", listID, err)
	}

	id := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(email))))

	member, err := list.GetMember(id, nil)
	if err != nil {
		apiErr, ok := err.(*gochimp3.APIError)

		if !ok && apiErr.Status != http.StatusNotFound {
			return fmt.Errorf("get mailchimp member %s: %w", email, err)
		}
	}

	if member != nil && (member.Status == "subscribed" || member.Status == "pending") {
		log.Printf("%s is already subscribed", email)
		return nil
	}

	req := &gochimp3.MemberRequest{
		EmailAddress: email,
		Status:       "pending",
	}

	_, err = list.AddOrUpdateMember(id, req)
	if err != nil {
		return fmt.Errorf("subscribe %s to mailchimp list: %w", email, err)
	}
	log.Printf("subscribed %s to mailchimp list", email)

	return nil
}
