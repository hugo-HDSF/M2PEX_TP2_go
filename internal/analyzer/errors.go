package analyzer

import "fmt"

// FileNotFoundError représente une erreur de fichier introuvable
type FileNotFoundError struct {
	FilePath string
	Err      error
}

func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf("fichier introuvable: %s", e.FilePath)
}

func (e *FileNotFoundError) Unwrap() error {
	return e.Err
}

// ParsingError représente une erreur de parsing
type ParsingError struct {
	FilePath string
	LogID    string
	Reason   string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("erreur de parsing pour %s (%s): %s", e.LogID, e.FilePath, e.Reason)
}

// NewFileNotFoundError crée une nouvelle erreur de fichier introuvable
func NewFileNotFoundError(filePath string, err error) *FileNotFoundError {
	return &FileNotFoundError{
		FilePath: filePath,
		Err:      err,
	}
}

// NewParsingError crée une nouvelle erreur de parsing
func NewParsingError(logID, filePath, reason string) *ParsingError {
	return &ParsingError{
		LogID:    logID,
		FilePath: filePath,
		Reason:   reason,
	}
}
