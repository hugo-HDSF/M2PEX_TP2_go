# GoLog Analyzer - TP2 EFREI
# Hugo HDSF - Mai 2025

BINARY_NAME=loganalyzer

.PHONY: help setup build run clean

help: ## Affiche l'aide
	@echo "GoLog Analyzer - TP2"
	@echo "==================="
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-10s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

setup: ## Installation des dépendances
	@echo "📦 Installation des dépendances..."
	@go mod tidy
	@echo "✅ Setup terminé"

build: ## Compile l'application
	@echo "🔨 Compilation..."
	@go build -o $(BINARY_NAME) .
	@echo "✅ Binaire créé: $(BINARY_NAME)"

run: build ## Lance l'analyse avec les exemples
	@echo "🏃 Test de l'application..."
	@./$(BINARY_NAME) analyze -c examples/config.json

demo: build ## Démonstration complète
	@echo "🎬 Démonstration:"
	@echo "1. Version:"
	@./$(BINARY_NAME) --version
	@echo "\n2. Analyse basique:"
	@./$(BINARY_NAME) analyze -c examples/config.json
	@echo "\n3. Avec export:"
	@./$(BINARY_NAME) analyze -c examples/config.json -o reports/2025/report.json

.DEFAULT_GOAL := help