package server

import (
	"backend/internal/form"
	"backend/internal/mailchimp"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func Start() {
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/newsletter", newsletterHandler)
	http.HandleFunc("/form", formHandler)

	log.Println("serving requests at *:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(w, "ok")
}

func validateForm(req *http.Request) error {
	log.Printf("got a request: %v", req)

	if req.Method != http.MethodPost {
		return errors.New("only post requests are accepted")
	}

	err := req.ParseForm()
	if err != nil {
		return fmt.Errorf("could not parse form: %v", err)
	}

	if req.FormValue("fortytwo") != "" {
		return fmt.Errorf("hidden form field is set, gonna ditch this form: %v", req.PostForm)
	}

	return nil
}

func newsletterHandler(w http.ResponseWriter, req *http.Request) {
	err := validateForm(req)
	if err != nil {
		log.Printf("form validation failed: %s", err)
		return
	}

	email := req.PostFormValue("email")
	err = mailchimp.Subscribe(email)
	if err != nil {
		log.Printf("mailchimp subscription for %s failed: %s", email, err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"ok":true}`)
}

func formHandler(w http.ResponseWriter, req *http.Request) {
	err := validateForm(req)
	if err != nil {
		log.Printf("form validation failed: %s", err)
		return
	}

	m := make(map[string]string)
	for key := range req.PostForm {
		if key == "fortytwo" {
			continue
		}

		m[key] = req.PostFormValue(key)
	}

	err = form.Submit(m)
	if err != nil {
		log.Printf("form submission failed: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"ok":true}`)
}
