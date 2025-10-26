<div>
<img src="docs/Cultpedia-banner.png">
</div>

# Cultpedia

Cultpedia is a repository of standardized, multilingual questions for educational platforms. Designed for Culturae, this project provides high-quality, schema-validated questions across various themes.

The Goal of Cultpedia is to offer a centralized question bank that can be easily integrated into different learning management systems (LMS) and quiz applications.

- [Features](#features)
- [Contributing](#contributing)
- [Project Structure](#project-structure)
- [Support](#support)

> [!IMPORTANT]
> The main Culturae platform is not yet available, but Cultpedia is being developed in parallel to provide ready-to-use content once the platform is live.

## Features

- **Multilingual Support**: Questions in English, French, and Spanish.
- **Schema Validation**: JSON Schema ensures data integrity.
- **Versioning**: Automatic versioning with manifest updates.
- **Interactive CLI**: Go-based tool for adding, validating, and managing questions.
- **SHA256 Checksums**: Data integrity verification for imports.
- **Full compatibility with Culturae**: Seamless integration with the Culturae platform.

## Contributing

If you wish to contribute, please refer to the [contributing guide](docs/CONTRIBUTING.md) for detailed instructions on how to add questions.

> [!IMPORTANT]
> For the moment we are accepting contributions only for the "general-knowledge" dataset. Future datasets may be added later.

Check the [Format](docs/FORMAT.md) to understand the json question structure.

## Project Structure

```
.
├── build.bat                   # Build script for Windows
├── build.sh                    # Build script for Unix
├── cmd/
│   └── main.go                 # CLI entry point
├── cultpedia.excalidraw        # Project diagram
├── datasets/
│   ├── general-knowledge/
│   │   ├── manifest.json       # Metadata and hashes
│   │   ├── questions.ndjson    # Main questions file
│   │   ├── subthemes.ndjson    # Subthemes
│   │   ├── tags.ndjson         # Tags
│   │   └── themes.ndjson       # Available themes
│   └── new-question.json       # New question template
├── docs/
│   ├── CONTRIBUTING.md         # Contribution guidelines
│   ├── FORMAT.md               # Data format specification
├── flake.lock                  # Nix lock file
├── flake.nix                   # Nix configuration
├── go.mod                      # Go module
├── go.sum                      # Go sum file
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
├── schemas/
│   ├── manifest.schema.json    # Manifest schema
│   ├── question.example.json   # Question example
│   └── question.schema.json    # Question schema
```

## Support

For questions or support, open an issue on GitHub or contact the Culturae/Cultpedia maintainers.