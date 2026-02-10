# Cultpedia API Documentation

Cultpedia provides a REST API to access all datasets programmatically.

## Table of Contents

- [Getting Started](#getting-started)
  - [Running the API](#running-the-api)
  - [Docker Deployment](#docker-deployment)
- [API Reference](#api-reference)
  - [Root - API Information](#root---api-information)
  - [Questions](#questions)
  - [Countries](#countries)
  - [Regions](#regions)
  - [Continents](#continents)
  - [Country Flags](#country-flags)
- [Examples](#examples)

---

## Getting Started

### Running the API

**Local (with Go):**
```bash
./cultpedia api 8080
```

The API will be available at `http://localhost:8080`

### Docker Deployment

**Build and run:**
```bash
docker build -t cultpedia-api .
docker run -d -p 8080:8080 cultpedia-api
```

**Stop the container:**
```bash
docker stop $(docker ps -q --filter ancestor=cultpedia-api)
```

---

## API Reference

### Root - API Information

**Endpoint:** `GET /api/`

Returns API metadata, available endpoints, dataset versions, and statistics.

**Response Example:**
```json
{
  "name": "Cultpedia API",
  "version": "1.0",
  "description": "API for Cultpedia questions and geography data",
  "datasets": {
    "general_knowledge": {
      "version": "1.0.11",
      "updated_at": "2025-12-24T10:49:03.556078904Z"
    },
    "geography": {
      "version": "1.0.1",
      "updated_at": "2025-12-26T18:44:37Z"
    }
  },
  "endpoints": [
    {
      "path": "/api/questions",
      "method": "GET",
      "description": "Get all questions"
    },
    {
      "path": "/api/geography/countries",
      "method": "GET",
      "description": "Get all countries"
    },
    {
      "path": "/api/geography/regions",
      "method": "GET",
      "description": "Get all regions"
    },
    {
      "path": "/api/geography/continents",
      "method": "GET",
      "description": "Get all continents"
    },
    {
      "path": "/api/geography/flags/{code}",
      "method": "GET",
      "description": "Get country flag SVG (use ISO Alpha2 code)"
    }
  ],
  "stats": {
    "questions": 13,
    "countries": 250,
    "regions": 22,
    "continents": 6
  }
}
```

---

### Questions

**Endpoint:** `GET /api/questions`

Returns all questions with their translations, answers, and metadata.

**Response Format:**
```json
{
  "data": [
    {
      "kind": "question",
      "version": "1.0",
      "slug": "history-french-revolution-start-year",
      "theme": {"slug": "history"},
      "subthemes": [{"slug": "french-revolution"}],
      "tags": [{"slug": "revolution"}, {"slug": "france"}],
      "qtype": "single_choice",
      "difficulty": "beginner",
      "estimated_seconds": 10,
      "points": 0.5,
      "shuffle_answers": true,
      "i18n": {
        "en": {
          "title": "French Revolution",
          "stem": "In what year did the French Revolution begin?",
          "explanation": "The French Revolution began in 1789, marking a major turning point in French history."
        },
        "fr": {
          "title": "Révolution française",
          "stem": "En quelle année la Révolution française a-t-elle commencé ?",
          "explanation": "La Révolution française a commencé en 1789, marquant un tournant majeur dans l'histoire de France."
        },
        "es": {
          "title": "Revolución Francesa",
          "stem": "¿En qué año comenzó la Revolución Francesa?",
          "explanation": "La Revolución Francesa comenzó en 1789, marcando un punto de inflexión importante en la historia de Francia."
        }
      },
      "answers": [
        {
          "slug": "1789",
          "is_correct": true,
          "i18n": {
            "en": {"label": "1789"},
            "fr": {"label": "1789"},
            "es": {"label": "1789"}
          }
        },
        {
          "slug": "1774",
          "is_correct": false,
          "i18n": {
            "en": {"label": "1774"},
            "fr": {"label": "1774"},
            "es": {"label": "1774"}
          }
        }
      ],
      "sources": [
        "https://en.wikipedia.org/wiki/French_Revolution",
        "https://www.britannica.com/event/French-Revolution"
      ]
    }
  ],
  "count": 13
}
```

**Question Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `kind` | string | Always `"question"` |
| `version` | string | Question version |
| `slug` | string | Unique identifier |
| `theme` | object | Main theme |
| `subthemes` | array | Related subthemes |
| `tags` | array | Associated tags |
| `qtype` | string | `"single_choice"` or `"true_false"` |
| `difficulty` | string | `"beginner"`, `"intermediate"`, `"advanced"`, `"pro"` |
| `estimated_seconds` | number | Time to answer |
| `points` | number | Scoring weight (0.5 to 5.0) |
| `shuffle_answers` | boolean | Randomize answer order |
| `i18n` | object | Translations (en, fr, es) |
| `answers` | array | Answer options |
| `sources` | array | Reference URLs |

---

### Countries

**Endpoint:** `GET /api/geography/countries`

Returns all countries with geographic data, flags, population, etc.

**Response Format:**
```json
{
  "data": [
    {
      "slug": "fr",
      "iso_alpha2": "FR",
      "iso_alpha3": "FRA",
      "iso_numeric": "250",
      "name": {
        "en": "France",
        "fr": "France",
        "es": "Francia"
      },
      "official_name": {
        "en": "French Republic",
        "fr": "République française",
        "es": "República Francesa"
      },
      "capital": {
        "en": "Paris",
        "fr": "Paris",
        "es": "París"
      },
      "continent": "europe",
      "region": "western_europe",
      "coordinates": {
        "lat": 46,
        "lng": 2
      },
      "flag": "fr",
      "population": 67391582,
      "area_km2": 551695,
      "currency": {
        "code": "EUR",
        "name": "Euro",
        "symbol": "€"
      },
      "languages": ["fr"],
      "neighbors": ["and", "bel", "deu", "ita", "lux", "mco", "esp", "che"],
      "tld": ".fr",
      "phone_code": "+33",
      "driving_side": "right",
      "independent": true
    }
  ],
  "count": 250
}
```

**Country Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `slug` | string | ISO Alpha2 code (lowercase) |
| `iso_alpha2` | string | ISO 3166-1 alpha-2 code |
| `iso_alpha3` | string | ISO 3166-1 alpha-3 code |
| `iso_numeric` | string | ISO 3166-1 numeric code |
| `name` | object | Country name in en, fr, es |
| `official_name` | object | Official name in en, fr, es |
| `capital` | object | Capital city in en, fr, es |
| `continent` | string | Continent identifier |
| `region` | string | Geographic region |
| `coordinates` | object | Latitude and longitude |
| `flag` | string | Flag filename (without extension) |
| `population` | number | Population count |
| `area_km2` | number | Area in square kilometers |
| `currency` | object | Currency code, name, symbol |
| `languages` | array | ISO 639-1 language codes |
| `neighbors` | array | Neighboring country codes |
| `tld` | string | Top-level domain |
| `phone_code` | string | International calling code |
| `driving_side` | string | `"left"` or `"right"` |
| `independent` | boolean | Whether the country is a sovereign state (`true` for ~199 sovereign states, `false` for ~51 territories/dependencies) |

---

### Regions

**Endpoint:** `GET /api/geography/regions`

Returns all geographic regions grouped by continent.

**Response Format:**
```json
{
  "data": [
    {
      "slug": "western_europe",
      "name": {
        "en": "Western Europe",
        "fr": "Europe de l'Ouest",
        "es": "Europa Occidental"
      },
      "continent": "europe",
      "countries": ["be", "fr", "lu", "mc", "nl"]
    }
  ],
  "count": 22
}
```

---

### Continents

**Endpoint:** `GET /api/geography/continents`

Returns all continents with their countries, area, and population.

**Response Format:**
```json
{
  "data": [
    {
      "slug": "europe",
      "name": {
        "en": "Europe",
        "fr": "Europe",
        "es": "Europa"
      },
      "countries": ["ad", "al", "at", "..."],
      "area_km2": 10180000,
      "population": 747707351
    }
  ],
  "count": 6
}
```

---

### Country Flags

**Endpoint:** `GET /api/geography/flags/{code}`

Returns the SVG flag for a specific country. Use the ISO Alpha2 code (lowercase).

**Parameters:**
- `{code}` - ISO 3166-1 alpha-2 country code (e.g., `fr`, `us`, `jp`)

**Response:** SVG image
- **Content-Type:** `image/svg+xml`

**Examples:**
```
GET /api/geography/flags/fr     # France flag
GET /api/geography/flags/us     # United States flag
GET /api/geography/flags/jp     # Japan flag
GET /api/geography/flags/de     # Germany flag
```

**Error Responses:**
- `400 Bad Request` - Country code required
- `404 Not Found` - Flag not found

---

## Examples

### Fetch all questions (JavaScript)

```javascript
fetch('http://localhost:8080/api/questions')
  .then(response => response.json())
  .then(data => {
    console.log(`Total questions: ${data.count}`);
    data.data.forEach(question => {
      console.log(question.i18n.en.title);
    });
  });
```

### Fetch a country flag (HTML)

```html
<img src="http://localhost:8080/api/geography/flags/fr" alt="France flag" />
```

### Get countries by region (JavaScript)

```javascript
fetch('http://localhost:8080/api/geography/countries')
  .then(response => response.json())
  .then(data => {
    const europeanCountries = data.data.filter(
      country => country.continent === 'europe'
    );
    console.log(`European countries: ${europeanCountries.length}`);
  });
```

### Fetch API info (curl)

```bash
curl http://localhost:8080/api/
```

### Download a flag (curl)

```bash
curl http://localhost:8080/api/geography/flags/fr -o france.svg
```

---

## Need Help?

For issues or questions, open an issue on [GitHub](https://github.com/Culturae-org/cultpedia/issues).
