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
- **Interactive CLI**: Go-based tool for adding, validating, and managing questions.
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
- `GET /api/geography/flags/{code}` - Country flag SVG

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
│               └── svg/        # Country flags (SVG format)
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
│   │   └── actions.go          # Actions logic
│   ├── checks/
│   │   └── checks.go           # Validation checks
│   ├── models/
│   │   └── question.go         # Data models
│   ├── ui/
│   │   └── ui.go               # TUI interface
│   └── utils/
│       └── utils.go            # Utilities
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
- [ ] Add flags format
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

MIT - [Culturae](https://culturae.me)
