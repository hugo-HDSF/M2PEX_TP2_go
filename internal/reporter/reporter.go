package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"m2pex_tp2_go/internal/analyzer"
)

// ExportToJSON exporte les résultats d'analyse vers un fichier JSON
func ExportToJSON(results []*analyzer.AnalysisResult, outputPath string) error {
	// Créer le fichier de sortie
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier de sortie '%s': %w", outputPath, err)
	}
	defer file.Close()

	// Encoder les résultats en JSON avec une indentation pour la lisibilité
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("erreur lors de l'encodage JSON: %w", err)
	}

	return nil
}

// GenerateReport crée un rapport formaté pour l'affichage console
func GenerateReport(results []*analyzer.AnalysisResult) string {
	if len(results) == 0 {
		return "Aucun résultat à afficher."
	}

	report := "=== RAPPORT D'ANALYSE ===\n\n"

	successCount := 0
	errorCount := 0

	for _, result := range results {
		status := "SUCCÈS"
		if result.Status == "FAILED" {
			status = "ÉCHEC"
			errorCount++
		} else {
			successCount++
		}

		report += fmt.Sprintf("ID: %s\n", result.LogID)
		report += fmt.Sprintf("Fichier: %s\n", result.FilePath)
		report += fmt.Sprintf("Statut: %s\n", status)
		report += fmt.Sprintf("Message: %s\n", result.Message)

		if result.ErrorDetails != "" {
			report += fmt.Sprintf("Détails erreur: %s\n", result.ErrorDetails)
		}

		report += "\n---\n\n"
	}

	report += fmt.Sprintf("TOTAL: %d | SUCCÈS: %d | ÉCHECS: %d\n", len(results), successCount, errorCount)

	return report
}
