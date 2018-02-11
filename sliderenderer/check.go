package sliderenderer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func (sr *SlideRenderer) CheckSlides() error {
	log.Printf("Checking slides")
	if sr.ShouldRecache() {
		log.Printf("Reading and preprocessing slides")
		if err := sr.RecacheSlides(); err != nil {
			return fmt.Errorf("failed to load slides: %s", err)
		}
	}
	if len(sr.CachedSlides) == 0 {
		return fmt.Errorf("0 slides detected")
	}
	for i := range sr.CachedSlides {
		rr := httptest.NewRecorder()
		sr.Serve(i, rr, httptest.NewRequest(http.MethodGet, "http://local/", nil))
		if rr.Code != http.StatusOK {
			return fmt.Errorf("rendering of slide %d failed with code %d", i, rr.Code)
		}
	}
	log.Printf("All checks passed.")
	return nil
}
