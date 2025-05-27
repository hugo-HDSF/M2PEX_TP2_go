package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LogConfig représente la configuration d'un log à analyser
type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// LoadConfig charge la configuration depuis un fichier JSON
func LoadConfig(configPath string) ([]LogConfig, error) {
	// Ouvrir le fichier de configuration
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir le fichier de configuration '%s': %w", configPath, err)
	}
	defer file.Close()

	// Décoder le JSON
	var logs []LogConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&logs); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage du JSON: %w", err)
	}

	// Valider la configuration
	if err := validateConfig(logs); err != nil {
		return nil, fmt.Errorf("configuration invalide: %w", err)
	}

	return logs, nil
}

// validateConfig valide la structure de la configuration
func validateConfig(logs []LogConfig) error {
	if len(logs) == 0 {
		return fmt.Errorf("aucun log configuré")
	}

	// Vérifier les champs requis et l'unicité des IDs
	seenIDs := make(map[string]bool)

	for i, log := range logs {
		if log.ID == "" {
			return fmt.Errorf("log à l'index %d: ID manquant", i)
		}
		if log.Path == "" {
			return fmt.Errorf("log '%s': chemin manquant", log.ID)
		}
		if log.Type == "" {
			return fmt.Errorf("log '%s': type manquant", log.ID)
		}

		// Vérifier l'unicité de l'ID
		if seenIDs[log.ID] {
			return fmt.Errorf("ID dupliqué: '%s'", log.ID)
		}
		seenIDs[log.ID] = true
	}

	return nil
}
