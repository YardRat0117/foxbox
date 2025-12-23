package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/YardRat0117/foxbox/internal/domain"
)

// ParseVolumes converts CLI volume arguments into domain VolumeSpec.
// $PWD and ~ are expanded
func parseVolumes(vols []string) []domain.VolumeSpec {
	specs := make([]domain.VolumeSpec, 0, len(vols))
	for _, v := range vols {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 {
			continue
		}

		host := parts[0]
		// expand ~
		if strings.HasPrefix(host, "~") {
			home, err := os.UserHomeDir()
			if err == nil {
				host = filepath.Join(home, host[1:])
			}
		}

		// expand $PWD
		host = os.ExpandEnv(host)

		specs = append(specs, domain.VolumeSpec{
			Host:  host,
			Guest: parts[1],
		})
	}
	return specs
}
