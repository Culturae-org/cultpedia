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

Returns questions with their translations, answers, and metadata.

**Query Parameters:**
- `theme`: Filter by theme slug (e.g., `history`)
- `subtheme`: Filter by subtheme slug (e.g., `french-revolution`)
- `tag`: Filter by tag slug (e.g., `france`)
- `difficulty`: Filter by difficulty (`beginner`, `intermediate`, `advanced`, `pro`)
- `type`: Filter by question type (`single_choice`, `true_false`)
- `page`: Page number (default: `1`)
- `limit`: Number of items per page (default: `50`, max: `500`)

**Response Format:**
```json
{
  "data": [...],
  "count": 10,
  "total": 13,
  "page": 1,
  "limit": 10
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
| `shuffle_answers` | boolean | Randomize answer order |
| `i18n` | object | Translations (en, fr, es) |
| `answers` | array | Answer options |
| `sources` | array | Reference URLs |

---

### Countries

**Endpoint:** `GET /api/geography/countries`

Returns countries with geographic data, flags, population, etc.

**Query Parameters:**
- `continent`: Filter by continent slug (e.g., `europe`)
- `region`: Filter by region slug (e.g., `western_europe`)
- `independent`: Filter by sovereignty status (`true` or `false`)
- `page`: Page number (default: `1`)
- `limit`: Number of items per page (default: `50`, max: `500`)

**Response Format:**
```json
{
  "data": [...],
  "count": 50,
  "total": 250,
  "page": 1,
  "limit": 50
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

### Fetch questions by theme (JavaScript)

```javascript
fetch('http://localhost:8080/api/questions?theme=history&limit=10')
  .then(response => response.json())
  .then(data => {
    console.log(`Questions in history: ${data.total}`);
    data.data.forEach(question => {
      console.log(question.i18n.en.title);
    });
  });
```

### Fetch a country flag (HTML)

```html
<img src="http://localhost:8080/api/geography/flags/fr" alt="France flag" />
```

### Get sovereign countries in Europe with pagination (JavaScript)

```javascript
fetch('http://localhost:8080/api/geography/countries?continent=europe&independent=true&limit=20')
  .then(response => response.json())
  .then(data => {
    console.log(`European sovereign states: ${data.total}`);
    console.log(`Page: ${data.page}/${Math.ceil(data.total / data.limit)}`);
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
