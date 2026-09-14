package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/koohq/forge-lite/internal/i18n"
	"github.com/koohq/forge-lite/internal/model"
)

const DefaultSaveFilePath = "save.json"

var ErrSaveNotFound = errors.New("save file not found")

// LoadSaveData loads and validates save.json.
func LoadSaveData(filePath string) (*model.SaveData, error) {
	if filePath == "" {
		filePath = DefaultSaveFilePath
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrSaveNotFound
		}
		return nil, err
	}

	var save model.SaveData
	if err := json.Unmarshal(data, &save); err != nil {
		return nil, fmt.Errorf("failed to parse save data: %w", err)
	}

	if save.Weapon.Slots == nil {
		save.Weapon.Slots = []model.OrbType{}
	}
	if save.StockOrbs == nil {
		save.StockOrbs = []model.OrbType{}
	}

	// Normalizing defaults
	defaultHP := 50
	if save.MaxHP == nil || *save.MaxHP <= 0 {
		save.MaxHP = &defaultHP
	}

	defaultDeepest := 0
	if save.DeepestFloor == nil || *save.DeepestFloor < 0 {
		save.DeepestFloor = &defaultDeepest
	}

	if save.Language == nil || (*save.Language != model.LanguageEN && *save.Language != model.LanguageJA) {
		detected := i18n.DetectDefaultLanguage()
		save.Language = &detected
	}

	defaultTheme := "classic_fantasy"
	if save.ThemeID == nil || *save.ThemeID == "" {
		save.ThemeID = &defaultTheme
	}

	return &save, nil
}

// SaveSaveData saves game state to JSON file.
func SaveSaveData(filePath string, data *model.SaveData) error {
	if filePath == "" {
		filePath = DefaultSaveFilePath
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode save data: %w", err)
	}

	return os.WriteFile(filePath, bytes, 0644)
}
