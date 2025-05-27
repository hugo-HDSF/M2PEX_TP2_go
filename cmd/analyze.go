package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"m2pex_tp2_go/internal/analyzer"
	"m2pex_tp2_go/internal/config"
	"m2pex_tp2_go/internal/reporter"

	"github.com/spf13/cobra"
)

var (
	configPath   string
	outputPath   string
	statusFilter string
)

// analyzeCmd représente la commande analyze
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse les fichiers de logs de manière concurrente",
	Long: `La commande analyze traite plusieurs fichiers de logs en parallèle.

Elle prend en entrée un fichier de configuration JSON contenant la liste des logs
à analyser et génère un rapport détaillé des résultats.

Fonctionnalités:
- Traitement concurrent avec goroutines
- Gestion d'erreurs personnalisées
- Export JSON des résultats
- Filtrage par statut (optionnel)

Exemple d'utilisation:
  loganalyzer analyze -c config.json -o rapport.json
  loganalyzer analyze --config logs_config.json --output /tmp/analyse_240524.json
  loganalyzer analyze -c config.json --status FAILED`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Drapeaux de la commande analyze
	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Chemin vers le fichier de configuration JSON (requis)")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Chemin de sortie pour le rapport JSON (optionnel)")
	analyzeCmd.Flags().StringVar(&statusFilter, "status", "", "Filtrer les résultats par statut (OK, FAILED)")

	// Marquer le drapeau config comme requis
	analyzeCmd.MarkFlagRequired("config")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Démarrage de l'analyse des logs...")

	// Vérifier si le fichier de configuration existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("le fichier de configuration '%s' n'existe pas", configPath)
	}

	// Charger la configuration
	logs, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement de la configuration: %w", err)
	}

	fmt.Printf("📋 Configuration chargée: %d logs à analyser\n", len(logs))

	// Canal pour collecter les résultats
	resultsChan := make(chan *analyzer.AnalysisResult, len(logs))

	// WaitGroup pour synchroniser les goroutines
	var wg sync.WaitGroup

	// Démarrer l'analyse de chaque log dans une goroutine séparée
	for _, logConfig := range logs {
		wg.Add(1)
		go func(lc config.LogConfig) {
			defer wg.Done()

			result := analyzer.AnalyzeLog(lc)
			resultsChan <- result
		}(logConfig)
	}

	// Goroutine pour fermer le canal une fois toutes les analyses terminées
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collecter tous les résultats
	var results []*analyzer.AnalysisResult
	for result := range resultsChan {
		// Appliquer le filtre de statut si spécifié
		if statusFilter == "" || result.Status == statusFilter {
			results = append(results, result)
		}
	}

	// Afficher le résumé sur la console
	displaySummary(results)

	// Exporter vers JSON si le chemin de sortie est spécifié
	if outputPath != "" {
		// Bonus: Ajouter l'horodatage au nom du fichier
		timestampedOutput := addTimestampToFilename(outputPath)

		// Bonus: Créer les répertoires si nécessaire
		if err := os.MkdirAll(filepath.Dir(timestampedOutput), 0755); err != nil {
			return fmt.Errorf("erreur lors de la création des répertoires: %w", err)
		}

		if err := reporter.ExportToJSON(results, timestampedOutput); err != nil {
			return fmt.Errorf("erreur lors de l'export JSON: %w", err)
		}

		fmt.Printf("\n💾 Rapport exporté vers: %s\n", timestampedOutput)
	}

	return nil
}

func displaySummary(results []*analyzer.AnalysisResult) {
	fmt.Println("\n📊 RÉSUMÉ DE L'ANALYSE")
	fmt.Println("====================")

	successCount := 0
	errorCount := 0

	for _, result := range results {
		status := "✅"
		if result.Status == "FAILED" {
			status = "❌"
			errorCount++
		} else {
			successCount++
		}

		fmt.Printf("%s [%s] %s - %s\n", status, result.LogID, result.FilePath, result.Message)
		if result.ErrorDetails != "" {
			fmt.Printf("   🔴 Erreur: %s\n", result.ErrorDetails)
		}
	}

	fmt.Printf("\n📈 STATISTIQUES\n")
	fmt.Printf("Total: %d | Succès: %d | Échecs: %d\n", len(results), successCount, errorCount)
}

func addTimestampToFilename(filename string) string {
	// Bonus: Ajouter l'horodatage au format AAMMJJ
	now := time.Now()
	timestamp := now.Format("060102") // Format AAMMJJ

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	nameWithoutExt := base[:len(base)-len(ext)]

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s", timestamp, nameWithoutExt, ext))
}
