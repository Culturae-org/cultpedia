# Data Format Specification

Cultpedia uses standardized JSON formats for its datasets. Data is stored in Newline-Delimited JSON (NDJSON) files for efficient streaming and version control.

## Table of Contents

- [Questions Dataset](#questions-dataset)
  - [Question Structure](#question-structure)
    - [Question Types](#question-types)
    - [Slug Format](#slug-format)
  - [NDJSON Files](#ndjson-files)
  - [Metadata](#metadata)
- [Geography Dataset](#geography-dataset)
  - [Country Structure](#country-structure)
    - [Country Fields](#country-fields)
  - [Continent Structure](#continent-structure)
  - [Region Structure](#region-structure)
  - [Supported Languages](#supported-languages)
  - [Flag Assets](#flag-assets)
  - [Data Sources](#data-sources)

---

# Questions Dataset

The `general-knowledge` dataset contains quiz questions and countries data in multiple languages.

## Question Structure

For the full schema, see [question.schema](../schemas/question.schema.json), or for an example question see [question.example](../schemas/question.example.json).


Each question is a JSON object with the following fields:

- `kind`: Always for the moment `"question"`
- `version`: Version string (e.g., `"1.0"`, incremented on edits ✗ not implemented yet)
- `slug`: Unique identifier (see Slug Format below)
- `theme`: Object with `slug` (e.g., `{"slug": "history"}`)
- `subthemes`: Array of objects with `slug` (e.g., `[{"slug": "ancient-history"}]`)
- `tags`: Array of objects with `slug` (e.g., `[{"slug": "capital-cities"}]`)
- `qtype`: `"single_choice"` or `"true_false"` (see Question Types below)
- `difficulty`: `"beginner"`, `"intermediate"`, `"advanced"`, or `"pro"`
- `estimated_seconds`: Number (time to answer, e.g., 20)
- `points`: Number (scoring weight, e.g., 1.0 - between 0.5 and 5.0)
- `shuffle_answers`: Boolean (whether to randomize answer order)
- `i18n`: Object with translations for `fr`, `en`, `es`:
  - Each language has `title`, `stem`, `explanation`
- `answers`: Array of answer objects (see Question Types for count requirements):
  - `slug`: Unique answer identifier
  - `is_correct`: Boolean (exactly one `true`)
  - `i18n`: Object with `label` for each language
- `sources`: Array of URLs (verifiable references)

---

### Question Types

#### Single Choice (`single_choice`)
Standard multiple choice questions with exactly **4 answers**.
- One answer must be correct (`is_correct: true`)
- Answer slugs can be any valid identifier

#### True/False (`true_false`)
Binary choice questions with exactly **2 answers**.
- One answer must be correct (`is_correct: true`)
- Answer slugs must be exactly `"true"` and `"false"`
- The `shuffle_answers` field is ignored for this type

---

### Slug Format

Recommended: `{theme}-{subtheme}-{key-element}-{specific-detail}`

- `{theme}`: Theme slug (e.g., `"art"`, `"science"`)
- `{subtheme}`: Primary subtheme (e.g., `"renaissance"`, `"chemistry"`)
- `{key-element}`: Main element (e.g., `artist name`, `chemical element`)
- `{specific-detail}`: Unique detail (e.g., `artwork`, `atomic number`)

Examples:
- `art-renaissance-botticelli-birth-of-venus`
- `science-chemistry-element-hydrogen-atomic-number`

> [!IMPORTANT]
> Slugs must be unique, lowercase, and use hyphens.

## NDJSON Files

- `questions.ndjson`: All questions, one per line.

Automatically generated files:
- `themes.ndjson`: List of themes.
- `subthemes.ndjson`: List of subthemes.
- `tags.ndjson`: List of tags.

## Metadata

`manifest.json` contains:
- Schema version
- Dataset info
- Type (e.g., "questions" or "geography" for the moment)
- Version
- Export timestamp
- Counts (e.g., number of questions)
- SHA256 hashes for integrity verification

> [!NOTE]
> Manifest will be generated automatically by `CI`.
> JSON Schemas are available at `schemas/manifest-questions.schema.json` and `schemas/manifest-geography.schema.json`.
> Examples: `schemas/manifest-questions.example.json` and `schemas/manifest-geography.example.json`.

---

# Geography Dataset

The geography dataset is a **reference dataset** (not quiz questions) containing geographic data for map games, flag quizzes, capital quizzes, etc.

## Country Structure

Each country in `countries.ndjson` is a JSON object, follow this to see the schema [country.schema](../schemas/country.schema.json) or an example [country.example](../schemas/country.example.json).

### Country Fields

| Field | Type | Description |
|-------|------|-------------|
| `kind` | string | Always `"country"` |
| `version` | string | Version string (e.g., `"1.0"`) |
| `slug` | string | Unique identifier (ISO 3166-1 alpha-2 code, lowercase) |
| `iso_alpha2` | string | ISO 3166-1 alpha-2 code (uppercase) |
| `iso_alpha3` | string | ISO 3166-1 alpha-3 code |
| `iso_numeric` | string | ISO 3166-1 numeric code |
| `name` | object | Country name in `en`, `fr`, `es` |
| `official_name` | object | Official name in `en`, `fr`, `es` |
| `capital` | object | Capital city in `en`, `fr`, `es` |
| `continent` | string | Continent ID |
| `region` | string | Geographic region ID |
| `coordinates` | object | `lat` and `lng` values |
| `flag` | string | Flag filename (without extension) |
| `population` | number | Population count |
| `area_km2` | number | Area in square kilometers |
| `currency` | object | `code`, `name`, `symbol` |
| `languages` | array | ISO 639-1 language codes |
| `neighbors` | array | Neighboring country IDs (alpha-3 lowercase) |
| `tld` | string | Top-level domain |
| `phone_code` | string | International calling code |
| `driving_side` | string | `"left"` or `"right"` |
| `independent` | boolean | Whether the country is a sovereign state |

> **`independent` field:** `true` for sovereign states (~199 entries), including non-UN members (Vatican, Kosovo, Palestine, Taiwan, Cook Islands, Niue). `false` for territories and dependencies (~51 entries: overseas territories, special administrative regions, Antarctica, etc.). This allows clients to filter: normal mode shows only sovereign states, hard mode includes everything.

## Continent Structure

```json
{
  "id": "europe",
  "name": {
    "en": "Europe",
    "fr": "Europe",
    "es": "Europa"
  },
  "countries": ["ad", "al", "at", "..."],
  "area_km2": 10180000,
  "population": 746000000
}
```

## Region Structure

```json
{
  "id": "western_europe",
  "name": {
    "en": "Western Europe",
    "fr": "Europe de l'Ouest",
    "es": "Europa Occidental"
  },
  "continent": "europe",
  "countries": ["at", "be", "ch", "de", "fr", "li", "lu", "mc", "nl"]
}
```

## Supported Languages

| Code | Language |
|------|----------|
| `en` | English (default, used for IDs) |
| `fr` | French |
| `es` | Spanish |

## Flag Assets

Flags are stored as SVG files in `assets/flags/svg/` with the country's ISO alpha-2 code as filename:

- `fr.svg` - France
- `de.svg` - Germany
- `us.svg` - United States

## Data Sources

- [APICountries](https://www.apicountries.com/) - Country data
- [REST Countries](https://restcountries.com) - Country data (MPL 2.0)
- [Flag-icons](https://github.com/lipis/flag-icons) - SVG flags (MIT License)
- [Natural Earth](https://naturalearthdata.com) - Geographic data (Public Domain)