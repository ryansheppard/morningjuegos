package discord

import (
	"encoding/json"
	"log/slog"
	"os"
)

type CopyPasta struct {
	PlayerID int64  `json:"player_id"`
	GuildID  int64  `json:"guild_id"`
	Message  string `json:"message"`
}

type CopyPastas struct {
	CopyPastas []CopyPasta `json:"copypastas"`
}

func (d *Discord) LoadCopyPastaFromJson(filepath string) error {
	dat, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	var copyPastas CopyPastas
	json.Unmarshal(dat, &copyPastas)

	formatted := make(map[int64]CopyPasta)
	for _, copyPasta := range copyPastas.CopyPastas {
		formatted[copyPasta.PlayerID] = copyPasta
	}

	if err != nil {
		slog.Error("Error loading copy pasta", "error", err)
		return err
	}

	d.copyPastas = formatted
	slog.Info("Loaded copy pasta", "copy_pastas", len(formatted))

	return nil
}
