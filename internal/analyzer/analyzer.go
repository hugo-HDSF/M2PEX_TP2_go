package analyzer

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"m2pex_tp2_go/internal/config"
)

// AnalysisResult représente le résultat de l'analyse d'un log
type AnalysisResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

// AnalyzeLog analyse un fichier de log et retourne le résultat
func AnalyzeLog(logConfig config.LogConfig) *AnalysisResult {
	result := &AnalysisResult{
		LogID:    logConfig.ID,
		FilePath: logConfig.Path,
	}

	// Vérifier si le fichier existe et est lisible
	if err := checkFileAccessibility(logConfig.Path); err != nil {
		var fileNotFoundErr *FileNotFoundError
		if errors.As(err, &fileNotFoundErr) {
			result.Status = "FAILED"
			result.Message = "Fichier introuvable."
			result.ErrorDetails = err.Error()
			return result
		}

		result.Status = "FAILED"
		result.Message = "Fichier inaccessible."
		result.ErrorDetails = err.Error()
		return result
	}

	// Simuler l'analyse avec un délai aléatoire
	analysisDelay := time.Duration(50+rand.Intn(151)) * time.Millisecond
	time.Sleep(analysisDelay)

	// Simuler une erreur de parsing aléatoire (10% de chance)
	if rand.Float32() < 0.1 {
		parsingErr := NewParsingError(logConfig.ID, logConfig.Path, "format de log non reconnu")
		result.Status = "FAILED"
		result.Message = "Erreur de parsing."
		result.ErrorDetails = parsingErr.Error()
		return result
	}

	// Succès
	result.Status = "OK"
	result.Message = "Analyse terminée avec succès."
	result.ErrorDetails = ""

	return result
}

// checkFileAccessibility vérifie si un fichier est accessible en lecture
func checkFileAccessibility(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileNotFoundError(filePath, err)
		}
		return err
	}

	// Vérifier si c'est un fichier (pas un répertoire)
	if info.IsDir() {
		return NewFileNotFoundError(filePath, errors.New("le chemin pointe vers un répertoire"))
	}

	// Tenter d'ouvrir le fichier pour vérifier les permissions de lecture
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}
