package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

// rootCmd représente la commande de base quand appelée sans sous-commandes
var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Un outil d'analyse de logs distribuée",
	Long: `GoLog Analyzer est un outil en ligne de commande développé en Go 
pour analyser des fichiers de logs provenant de diverses sources de manière concurrente.

Cet outil permet de:
- Traiter plusieurs logs en parallèle avec des goroutines
- Gérer les erreurs de manière robuste avec des erreurs personnalisées
- Exporter les résultats au format JSON
- Centraliser l'analyse de logs multiples

Développé dans le cadre du TP2 - Master Dev Manager Fullstack EFREI
Cours de Go - Madame Axelle Lanca`,
	Version: version,
}

// Execute ajoute toutes les commandes enfants à la commande racine et définit les drapeaux appropriés.
// Ceci est appelé par main.main(). Il ne doit être appelé qu'une seule fois sur rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Ici vous pouvez définir vos drapeaux et votre configuration.

	// Cobra prend également en charge les drapeaux persistants, qui, s'ils sont définis ici,
	// seront globaux pour votre application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "fichier de configuration (par défaut est $HOME/.loganalyzer.yaml)")

	// Cobra prend également en charge les drapeaux locaux, qui ne s'exécuteront
	// que lorsque cette action sera appelée directement.
	rootCmd.Flags().BoolP("toggle", "t", false, "Message d'aide pour toggle")
}
