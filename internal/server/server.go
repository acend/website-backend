package server

import (
	"backend/internal/form"
	"backend/internal/mailchimp"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
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

func validateForm(req *http.Request, isNewsletter bool) error {
	log.Printf("got a request: %v", req)

	if req.Method != http.MethodPost {
		return errors.New("only post requests are accepted")
	}

	err := req.ParseForm()
	log.Printf("form contents: %s", req.PostForm)

	if err != nil {
		return fmt.Errorf("parse form: %v", err)
	}

	if req.FormValue("url") != "" {
		return fmt.Errorf("form submission is probably spam, gonna ditch it")
	}

	for k := range req.PostForm {
		v := req.PostFormValue(k)

		if strings.Contains(v, "http://") {
			return fmt.Errorf("form field %v contains an http (not https) link which is probably spam", k)
		}

		if strings.Contains(v, "//bit.ly/") {
			return fmt.Errorf("form field %v contains a bit.ly link which is probably spam", k)
		}
	}

	return nil
}

func newsletterHandler(w http.ResponseWriter, req *http.Request) {
	err := validateForm(req, true)
	if err != nil {
		log.Printf("form validation failed: %s", err)
		returnOK(w)
		return
	}

	email := req.PostFormValue("email")

	go func() {
		err = mailchimp.Subscribe(email)
		if err != nil {
			log.Printf("mailchimp subscription for %s failed: %s", email, err)
		}
	}()

	returnOK(w)
}

func formHandler(w http.ResponseWriter, req *http.Request) {
	err := validateForm(req, false)
	if err != nil {
		log.Printf("form validation failed: %s", err)
		returnOK(w)
		return
	}

	m := make(map[string]string)
	for key := range req.PostForm {
		if key == "url" || key == "human" {
			continue
		}

		m[key] = req.PostFormValue(key)
	}

	go func() {
		err = form.Submit(m)
		if err != nil {
			log.Printf("error on form submission: %v", err)
		}
	}()

	returnOK(w)
}

func returnOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"ok":true}`)
}
