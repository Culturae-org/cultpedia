<img width="1200" height="475" alt="Cultpedia-banner" src="https://github.com/user-attachments/assets/9fd8628c-3ad4-4f8b-b396-81c049dbb7e0" />

# Cultpedia

Knowledge game distributed server.

Cultpedia is a repository of standardized, multilingual questions, and countries data for educational platforms. Designed for Culturae, this project provides high-quality, schema-validated questions across various themes.

The Goal of Cultpedia is to offer a centralized question bank that can be easily integrated into different learning management systems (LMS) and quiz applications.

- [Features](#features)
- [API](#api)
- [Contributing](#contributing)
- [Project Structure](#project-structure)
- [Support](#support)

> [!WARNING]
> The main Culturae platform is not yet available, but Cultpedia is being developed in parallel to provide ready-to-use content once the platform is live.

## Features

- **Multilingual Support**: English, French, and Spanish.
- **Schema Validation**: JSON Schema ensures data integrity.
- **Versioning**: Automatic versioning with manifest updates.
- **CLI Tool**: Go-based tool for adding, validating, and managing questions.
- **SHA256 Checksums**: Data integrity verification for imports.
- **Full compatibility with Culturae**: Seamless integration with the Culturae platform.

## Importing Datasets

> [!NOTE]
> To import datasets into Culturae (or any compatible platform).

Culturae use manifest.json files to import datasets with all the needed metadata and sha256 checksums.

### Importing General Knowledge Dataset

```
https://raw.githubusercontent.com/Culturae-org/cultpedia/refs/heads/main/datasets/general-knowledge/manifest.json
```

### Importing Geography Dataset

```
https://raw.githubusercontent.com/Culturae-org/cultpedia/refs/heads/main/datasets/geography/manifest.json
```

## API

Cultpedia provides a REST API to access all datasets programmatically.

API is available at [api.culturae.me](https://api.culturae.me).

**Available endpoints:**
- `GET /api/` - API information and stats
- `GET /api/questions` - All questions
- `GET /api/geography/countries` - All countries
- `GET /api/geography/regions` - All regions
- `GET /api/geography/continents` - All continents
- `GET /api/geography/flags/{code}` - Country flag (SVG format)
- `GET /api/geography/flags/{code}/png` - Country flag (PNG format, 512px)
- `GET /api/geography/flags/{code}/png/1024` - Country flag (PNG format, 1024px)

**[Full API Documentation](docs/API.md)**

## Contributing

If you wish to contribute, please refer to the [contributing guide](docs/CONTRIBUTING.md) for detailed instructions on how to add questions.

> [!IMPORTANT]
> For the moment we are accepting contributions only for the "general-knowledge" dataset. Future datasets may be added later.

Check the [Format](docs/FORMAT.md) to understand the json question, and countrie structure.

### Create your own dataset with Cultpedia format.

> [!NOTE]
> Only available for questions datasets (not geography).

You can create your own dataset following the Cultpedia format. Use the *cultpedia* CLI to generate the essential files. 

```
./cultpedia init {dataset-name}
```

## AI-Assisted Question Generation

Cultpedia includes a built-in AI question generator powered by [Perplexity AI](https://www.perplexity.ai/). It uses the [LLMS.md](LLMS.md) specification as system prompt to ensure generated questions follow the Cultpedia format.

### Setup

```bash
# Set your Perplexity API key (see .env.example)
export PPLX_API_KEY="pplx-xxxx..."
```

### CLI Commands

```bash
# Generate questions on a topic
./cultpedia generate "French Revolution" --theme history --count 5

# Options
#   --theme <slug>        Theme slug (default: science)
#   --count <n>           Number of questions, 1-10 (default: 3)
#   --difficulty <level>  beginner, intermediate, advanced (default: intermediate)
#   --type <qtype>        single_choice or true_false (default: single_choice)

# Review pending questions
./cultpedia pending

# Approve a question (adds it to the main dataset)
./cultpedia approve <index>

# Reject a question (removes it from pending)
./cultpedia reject <index>
```

### Workflow

1. **Generate** -- Questions are generated via Perplexity AI and saved to `datasets/pending-questions.ndjson`
2. **Review** -- Use `pending` to preview each question (slug, stem, theme, sources)
3. **Approve/Reject** -- Approved questions are moved to the main dataset; rejected ones are discarded
4. **Validate & PR** -- Run `./cultpedia validate`, then create a PR

> [!TIP]
> Always verify the sources provided by the AI before approving. Factual accuracy is our top priority.

### Using LLMS.md with external AI

You can also use [LLMS.md](LLMS.md) directly with any AI assistant:

```
Using the LLMS.md guide from Cultpedia, generate a quiz question about [TOPIC].
```

**[Full AI Guide](LLMS.md)**


## Project Structure

```
.
├── build.bat                   # Build script for Windows
├── build.sh                    # Build script for Unix
├── cmd/
│   └── main.go                 # CLI entry point
├── datasets/
│   ├── general-knowledge/
│   │   ├── manifest.json       # Metadata and hashes
│   │   ├── questions.ndjson    # Main questions file
│   │   ├── subthemes.ndjson    # Subthemes
│   │   ├── tags.ndjson         # Tags
│   │   └── themes.ndjson       # Available themes
│   ├── new-question.json       # New question template
│   │
│   └── geography/
│       ├── manifest.json       # Metadata and hashes
│       ├── countries.ndjson    # Main Countries file
│       ├── continents.ndjson   # Continents file
│       ├── regions.ndjson      # Regions file
│       └── assets/
│           └── flags/
│               ├── png_512/         # Country flags (PNG, 512px)
│               ├── png_1024/        # Country flags (PNG, 1024px)
│               └── svg/             # Country flags (SVG format)
├── docs/
│   ├── CONTRIBUTING.md         # Contribution guidelines
│   ├── FORMAT.md               # Data format specification
│   └── MEDIA                   # All Media
│
├── flake.lock                  # Nix lock file
├── flake.nix                   # Nix configuration
├── go.mod                      # Go module
├── go.sum                      # Go sum file
│
├── internal/
│   ├── actions/
│   │   ├── actions.go          # Question CRUD, theme sync, versioning
│   │   ├── api.go              # REST API server
│   │   └── generate.go         # AI question generation (Perplexity)
│   ├── checks/
│   │   ├── checks.go           # Question validation
│   │   └── geography.go        # Geography validation
│   ├── models/
│   │   ├── question.go         # Question data models
│   │   └── geography.go        # Geography data models
│   └── utils/
│       └── utils.go            # File I/O, path constants
|
├── schemas/
│   ├── manifest-geography.example.json  # Geography manifest example
│   ├── manifest-geography.schema.json   # Geography manifest schema
│   ├── manifest-questions.example.json  # Questions manifest example
│   ├── manifest-questions.schema.json   # Questions manifest schema
│   ├── question.example.json            # Question example
│   ├── question.schema.json             # Question schema
│   ├── country.example.json             # Country example
│   └── country.schema.json              # Country schema
```

# Todo

- [x] Cultpedia CLI
- [x] QCM dataset structure 
- [x] CI question validation
- [x] CI sync + bump version
- [x] Auto check version
- [x] Countries data
- [x] Add true / false questions
- [ ] CLI countries tool
- [ ] CLI edit tool
- [ ] Geography validator
- [ ] Exporter (csv for example)
- [ ] Branchs by theme
- [x] Add flags format
- [ ] Contribution cli help
- [ ] Contribution Helper
- [ ] More questions !

Missings svg flags for countries.

- [x] cz : Czech Republic
- [x] hn : Honduras
- [x] io : British Indian Ocean Territory
- [x] iq : Iraq
- [x] mm : Myanmar
- [x] qa : Qatar
- [x] th : Thailand
- [x] tw : Taiwan
- [x] ws : Samoa
- [x] xk : Kosovo

## Support

For questions or support, open an issue on GitHub or contact the Culturae/Cultpedia maintainers or open an issue on GitHub.

## License

[LICENSE](./LICENSE) - [Culturae](https://culturae.me)
