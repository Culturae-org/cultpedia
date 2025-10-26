# Data Format Specification

Cultpedia uses a standardized JSON format for questions. Data is stored in Newline-Delimited JSON (NDJSON) files for efficient streaming and version control.

## Question Structure

For the full schema, see [question.schema](../schemas/question.schema.json), or for an example question see [question.example](../schemas/question.example.json).


Each question is a JSON object with the following fields:

- `kind`: Always for the moment `"question"`
- `version`: Version string (e.g., `"1.0"`, incremented on edits ✗ not implemented yet)
- `slug`: Unique identifier (see Slug Format below)
- `theme`: Object with `slug` (e.g., `{"slug": "history"}`)
- `subthemes`: Array of objects with `slug` (e.g., `[{"slug": "ancient-history"}]`)
- `tags`: Array of objects with `slug` (e.g., `[{"slug": "capital-cities"}]`)
- `qtype`: `"single_choice"` (only supported type for the moment)
- `difficulty`: `"beginner"`, `"intermediate"`, or `"advanced"`
- `estimated_seconds`: Number (time to answer, e.g., 20)
- `points`: Number (scoring weight, e.g., 1.0 - between 0.5 and 5.0)
- `shuffle_answers`: Boolean (whether to randomize answer order)
- `i18n`: Object with translations for `fr`, `en`, `es`:
  - Each language has `title`, `stem`, `explanation`
- `answers`: Array of exactly 4 answer objects:
  - `slug`: Unique answer identifier
  - `is_correct`: Boolean (exactly one `true`)
  - `i18n`: Object with `label` for each language
- `sources`: Array of URLs (verifiable references)

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
- Version
- Export timestamp
- Counts (e.g., number of questions)
- SHA256 hashes for integrity verification

> [!NOTE]
> Manifest will be generated automatically by `CI`.