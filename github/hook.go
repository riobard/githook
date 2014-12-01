package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Hook struct {
	secret string
	cmd    string
}

func NewHook(secret, command string) *Hook {
	return &Hook{secret, command}
}

func (h *Hook) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.UserAgent(), "GitHub-Hookshot/") {
		log.Printf("invalid User-Agent")
		http.Error(w, "invalid User-Agent", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("unsupported Content-Type")
		http.Error(w, "unsupported Content-Type", http.StatusBadRequest)
		return
	}

	// Split request body stream and feed into HMAC verification
	hash := hmac.New(sha1.New, []byte(h.secret))
	tee := io.TeeReader(r.Body, hash)

	var evt PushEvent
	if err := json.NewDecoder(tee).Decode(&evt); err != nil {
		log.Printf("invalid JSON payload")
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Verify the authenticity of the payload
	sig := r.Header.Get("X-Hub-Signature")
	sig2 := fmt.Sprintf("sha1=%x", hash.Sum(nil))
	if !hmac.Equal([]byte(sig2), []byte(sig)) {
		log.Printf("invalid signature")
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	uid := r.Header.Get("X-Github-Delivery")
	cmd := exec.Command(h.cmd)
	cmd.Env = []string{
		"GITHOOK_SOURCE=github",
		"GITHOOK_GITHUB_EVENT=" + r.Header.Get("X-Github-Event"),
		"GITHOOK_GITHUB_DELIVERY=" + r.Header.Get("X-GitHub-Delivery"),
	}

	if err := cmd.Run(); err != nil {
		log.Printf("GitHub delivery ID %s: failed to execute command %q: %s", uid, h.cmd, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
