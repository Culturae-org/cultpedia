package actions

import (
	"cultpedia/internal/models"
	"cultpedia/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var apiData models.APIData

const (
	defaultPort     = "8080"
	requestsPerMin  = 60
	cleanupInterval = 10 * time.Minute
)

type Visitor struct {
	lastSeen time.Time
	count    int
}

var (
	visitors = make(map[string]*Visitor)
	mu       sync.RWMutex
)

func RunAPIServer(serverPort ...string) {
	port := defaultPort
	if len(serverPort) > 0 && serverPort[0] != "" {
		port = serverPort[0]
	}

	if err := loadData(); err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	go cleanupVisitors()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/questions", handleQuestions)
	mux.HandleFunc("/api/geography/countries", handleCountries)
	mux.HandleFunc("/api/geography/regions", handleRegions)
	mux.HandleFunc("/api/geography/continents", handleContinents)
	mux.HandleFunc("/api/geography/flags/", handleFlags)
	mux.HandleFunc("/api/", handleRoot)
	mux.HandleFunc("/", handleRoot)

	handler := limitMiddleware(mux)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Cultpedia API server running on http://localhost:%s\n", port)
	fmt.Printf("Rate limit: %d requests per minute per IP\n", requestsPerMin)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]

		if !allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func allow(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		visitors[ip] = &Visitor{lastSeen: time.Now(), count: 1}
		return true
	}

	if time.Since(v.lastSeen) > time.Minute {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	if v.count >= requestsPerMin {
		v.lastSeen = time.Now()
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

func cleanupVisitors() {
	for {
		time.Sleep(cleanupInterval)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > cleanupInterval {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func loadData() error {
	var err error

	if err := loadManifests(); err != nil {
		return fmt.Errorf("error loading manifests: %w", err)
	}

	apiData.Questions, err = utils.LoadQuestions()
	if err != nil {
		return fmt.Errorf("error loading questions: %w", err)
	}

	apiData.Countries, err = utils.LoadCountries()
	if err != nil {
		return fmt.Errorf("error loading countries: %w", err)
	}

	apiData.Regions, err = utils.LoadRegions()
	if err != nil {
		return fmt.Errorf("error loading regions: %w", err)
	}

	apiData.Continents, err = utils.LoadContinents()
	if err != nil {
		return fmt.Errorf("error loading continents: %w", err)
	}

	return nil
}

func loadManifests() error {
	geoFile, err := os.Open(utils.GeographyManifestFile)
	if err != nil {
		return err
	}
	defer func() { _ = geoFile.Close() }()

	var geoManifest struct {
		Version   string `json:"version"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := json.NewDecoder(geoFile).Decode(&geoManifest); err != nil {
		return err
	}
	apiData.Manifests.Geography = models.ManifestInfo{
		Version:   geoManifest.Version,
		UpdatedAt: geoManifest.UpdatedAt,
	}

	qFile, err := os.Open(utils.ManifestFile)
	if err != nil {
		return err
	}
	defer func() { _ = qFile.Close() }()

	var qManifest struct {
		Version   string `json:"version"`
		UpdatedAt string `json:"updated_at"`
	}
	if err := json.NewDecoder(qFile).Decode(&qManifest); err != nil {
		return err
	}
	apiData.Manifests.GeneralKnowledge = models.ManifestInfo{
		Version:   qManifest.Version,
		UpdatedAt: qManifest.UpdatedAt,
	}

	return nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := models.APIRootResponse{
		Name:        "Cultpedia API",
		Version:     utils.Version,
		Description: "API for Cultpedia questions and geography data",
		Datasets: map[string]interface{}{
			"general_knowledge": map[string]string{
				"version":    apiData.Manifests.GeneralKnowledge.Version,
				"updated_at": apiData.Manifests.GeneralKnowledge.UpdatedAt,
			},
			"geography": map[string]string{
				"version":    apiData.Manifests.Geography.Version,
				"updated_at": apiData.Manifests.Geography.UpdatedAt,
			},
		},
		Endpoints: []models.EndpointInfo{
			{
				Path:        "/api/questions",
				Method:      "GET",
				Description: "Get questions with filters (theme, subtheme, tag, difficulty, type) and pagination (page, limit)",
			},
			{
				Path:        "/api/geography/countries",
				Method:      "GET",
				Description: "Get countries with filters (continent, region, independent) and pagination (page, limit)",
			},
			{
				Path:        "/api/geography/regions",
				Method:      "GET",
				Description: "Get all regions",
			},
			{
				Path:        "/api/geography/continents",
				Method:      "GET",
				Description: "Get all continents",
			},
			{
				Path:        "/api/geography/flags/{code}",
				Method:      "GET",
				Description: "Get country flag SVG (use ISO Alpha2 code)",
			},
		},
		Stats: map[string]int{
			"questions":  len(apiData.Questions),
			"countries":  len(apiData.Countries),
			"regions":    len(apiData.Regions),
			"continents": len(apiData.Continents),
		},
	}

	_ = json.NewEncoder(w).Encode(response)
}

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query()
	theme := query.Get("theme")
	subtheme := query.Get("subtheme")
	tag := query.Get("tag")
	difficulty := query.Get("difficulty")
	qtype := query.Get("type")

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if limit > 500 {
		limit = 500
	}

	filtered := make([]models.Question, 0)
	for _, q := range apiData.Questions {
		if theme != "" && q.Theme.Slug != theme {
			continue
		}
		if difficulty != "" && q.Difficulty != difficulty {
			continue
		}
		if qtype != "" && q.Qtype != qtype {
			continue
		}

		if subtheme != "" {
			found := false
			for _, s := range q.Subthemes {
				if s.Slug == subtheme {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if tag != "" {
			found := false
			for _, t := range q.Tags {
				if t.Slug == tag {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, q)
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginated := filtered[start:end]

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  paginated,
		"count": len(paginated),
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func handleCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query()
	continent := query.Get("continent")
	region := query.Get("region")
	independentStr := query.Get("independent")

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if limit > 500 {
		limit = 500
	}

	filtered := make([]models.Country, 0)
	for _, c := range apiData.Countries {
		if continent != "" && c.Continent != continent {
			continue
		}
		if region != "" && c.Region != region {
			continue
		}
		if independentStr != "" {
			if indep, err := strconv.ParseBool(independentStr); err == nil {
				if c.Independent != indep {
					continue
				}
			}
		}
		filtered = append(filtered, c)
	}

	total := len(filtered)

	start := (page - 1) * limit
	end := start + limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	paginated := filtered[start:end]

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  paginated,
		"count": len(paginated),
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func handleRegions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Regions,
		"count": len(apiData.Regions),
	})
}

func handleContinents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  apiData.Continents,
		"count": len(apiData.Continents),
	})
}

func handleFlags(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/api/geography/flags/")
	code = strings.TrimSuffix(code, ".svg")
	code = strings.ToLower(code)

	if code == "" {
		http.Error(w, "Country code required", http.StatusBadRequest)
		return
	}

	if len(code) != 2 || !isAlphaOnly(code) {
		http.Error(w, "Invalid country code format", http.StatusBadRequest)
		return
	}

	flagPath := filepath.Join(utils.FlagsSVGDir, code+".svg")

	if _, err := os.Stat(flagPath); os.IsNotExist(err) {
		http.Error(w, "Flag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, flagPath)
}

func isAlphaOnly(s string) bool {
	for _, c := range s {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}
