package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"m2pex_tp2_go/internal/config"

	"github.com/spf13/cobra"
)

var (
	addLogID   string
	addLogPath string
	addLogType string
	configFile string
)

// addLogCmd représente la commande add-log
var addLogCmd = &cobra.Command{
	Use:   "add-log",
	Short: "Ajoute une nouvelle configuration de log au fichier config.json",
	Long: `La commande add-log permet d'ajouter manuellement une nouvelle configuration
de log à un fichier config.json existant.

Exemple d'utilisation:
  loganalyzer add-log --id web-server-3 --path /var/log/apache/access.log --type apache-access --file config.json`,
	RunE: runAddLog,
}

func init() {
	rootCmd.AddCommand(addLogCmd)

	// Drapeaux requis
	addLogCmd.Flags().StringVar(&addLogID, "id", "", "Identifiant unique du log (requis)")
	addLogCmd.Flags().StringVar(&addLogPath, "path", "", "Chemin vers le fichier de log (requis)")
	addLogCmd.Flags().StringVar(&addLogType, "type", "", "Type du log (requis)")
	addLogCmd.Flags().StringVar(&configFile, "file", "", "Chemin vers le fichier config.json (requis)")

	// Marquer tous les drapeaux comme requis
	addLogCmd.MarkFlagRequired("id")
	addLogCmd.MarkFlagRequired("path")
	addLogCmd.MarkFlagRequired("type")
	addLogCmd.MarkFlagRequired("file")
}

func runAddLog(cmd *cobra.Command, args []string) error {
	// Vérifier si le fichier de configuration existe
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("le fichier de configuration '%s' n'existe pas", configFile)
	}

	// Charger la configuration existante
	existingLogs, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement de la configuration: %w", err)
	}

	// Vérifier l'unicité de l'ID
	for _, log := range existingLogs {
		if log.ID == addLogID {
			return fmt.Errorf("un log avec l'ID '%s' existe déjà", addLogID)
		}
	}

	// Créer la nouvelle configuration
	newLog := config.LogConfig{
		ID:   addLogID,
		Path: addLogPath,
		Type: addLogType,
	}

	// Ajouter à la liste existante
	updatedLogs := append(existingLogs, newLog)

	// Sauvegarder dans le fichier
	if err := saveConfig(updatedLogs, configFile); err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde: %w", err)
	}

	fmt.Printf("✅ Log '%s' ajouté avec succès au fichier '%s'\n", addLogID, configFile)
	fmt.Printf("   Chemin: %s\n", addLogPath)
	fmt.Printf("   Type: %s\n", addLogType)

	return nil
}

func saveConfig(logs []config.LogConfig, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(logs)
}
