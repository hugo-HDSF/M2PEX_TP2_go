# GoLog Analyzer 🔍

![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

**Outil d'analyse de logs distribuée en Go** - TP2 Master Dev Manager Fullstack EFREI

## 🎯 Contexte

Projet développé pour le cours de Go de **Madame Axelle Lanca** - EFREI 2025.  
Démonstration des concepts : goroutines, erreurs personnalisées, CLI Cobra, JSON, architecture modulaire.

## ⚡ Installation rapide

```bash
git clone https://github.com/hugo-HDSF/m2pex_tp2_go.git
cd m2pex_tp2_go
make setup && make build
```

## 🚀 Utilisation

```bash
# Analyse basique
./loganalyzer analyze -c examples/config.json

# Avec export JSON horodaté
./loganalyzer analyze -c examples/config.json -o reports/2025/report.json

# Filtrage par statut
./loganalyzer analyze -c examples/config.json --status FAILED

# Ajout de log
./loganalyzer add-log --id server-3 --path /var/log/app.log --type app --file examples/config.json

# Démonstration complète
make demo
```

## 🎬 Démonstration

```bash
make demo  # Démo complète avec tous les cas d'usage
make help  # Liste des commandes Makefile
```

## 📋 Configuration JSON

```json
[
  {
  "id": "web-server-1",
  "path": "/var/log/nginx/access.log",
  "type": "nginx-access"
  }
]
```

## 🏗️ Architecture technique

```
m2pex_tp2_go/
├── main.go                 # Point d'entrée
├── Makefile               # Automatisation
├── cmd/                   # CLI Cobra
│   ├── root.go           # Commande racine
│   ├── analyze.go        # Analyse concurrente
│   └── add-log.go        # Ajout de configuration
├── internal/              # Packages métier
│   ├── config/           # Chargement JSON + validation
│   ├── analyzer/         # Goroutines + erreurs personnalisées
│   └── reporter/         # Export JSON
└── examples/             # Fichiers de test
```

### Implémentation clé

**Concurrence** : Goroutines + WaitGroup + Canal buffé
```go
var wg sync.WaitGroup
resultsChan := make(chan *AnalysisResult, len(logs))
for _, log := range logs {
  wg.Add(1)
  go func(l LogConfig) {
    defer wg.Done()
    resultsChan <- analyzer.AnalyzeLog(l)
  }(log)
}
```

**Erreurs personnalisées** : `FileNotFoundError` + `ParsingError` avec `errors.Is/As`
```go
type FileNotFoundError struct {
  FilePath string
  Err      error
}
func (e *FileNotFoundError) Unwrap() error { return e.Err }
```

**Simulation réaliste** :
- Vérification accessibilité fichiers
- Délai aléatoire 50-200ms
- Export JSON avec horodatage automatique (AAMMJJ)

## ✨ Fonctionnalités

**Core** : Analyse concurrente, erreurs robustes, CLI intuitive, export JSON  
**Bonus** : Horodatage auto, création répertoires, filtrage statut, ajout dynamique logs

## 👥 Équipe

**Développeurs** : Hugo Da Silva ([@hugo-HDSF](https://github.com/hugo-HDSF)), Shalom Ango  
**EFREI** : Master Dev Manager Fullstack - Go/Mme Axelle Lanca - Mai 2025

---
*Architecture modulaire • Programmation concurrente • Gestion d'erreurs avancée*