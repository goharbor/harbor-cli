package quota

import (
	"errors"
	"os"
	"regexp"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/quota/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type QuotaUpdateReq struct {
	// The new hard limits for the quota
	Hard ResourceList `json:"hard,omitempty"`
}

type ResourceList map[string]int64

// UpdateQuotaCommand updates the quota
func UpdateQuotaCommand() *cobra.Command {
	var (
		quotaID int64
		storage string
	)

	cmd := &cobra.Command{
		Use:   "update [QuotaID]",
		Short: "update project quotas for projects",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var storageValue int64

			if len(args) > 0 {
				quotaID, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				quotaID = prompt.GetQuotaIDFromUser()
			}

			if err != nil {
				log.Errorf("failed to parse registry id: %v", err)
			}

			if storage != "" {
				if storage == "-1" {
					storageValue = -1
				} else {
					storageValue, err = storageStringToBytes(storage)
					if err != nil {
						log.Errorf("failed to parse storage: %v", err)
						os.Exit(1)
					}
				}
			} else {
				storage = update.UpdateQuotaView()
				storageValue, err = storageStringToBytes(storage)
			}

			hardlimit := &models.QuotaUpdateReq{
				Hard: models.ResourceList{"storage": storageValue},
			}

			err = api.UpdateQuota(quotaID, hardlimit)
			if err != nil {
				log.Errorf("failed to update quota: %v", err)
				os.Exit(1)
			}

			log.Infof("quota updated successfully!")
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&storage, "storage", "", "", "Enter storage size (e.g., 50GiB, 200MiB, 4TiB)")

	return cmd
}

func storageStringToBytes(storage string) (int64, error) {
	// Define the conversion multipliers
	multipliers := map[string]int64{
		"MiB": 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
	}

	// Define the regex to parse the input string
	re := regexp.MustCompile(`^(\d+)(MiB|GiB|TiB)$`)
	matches := re.FindStringSubmatch(storage)
	if matches == nil {
		return 0, errors.New("invalid storage format")
	}

	// Extract the value and unit from the matches
	valueStr, unit := matches[1], matches[2]
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Calculate the value in bytes
	bytes := value * multipliers[unit]

	// Check if the value exceeds 1024 TB
	maxBytes := 1024 * 1024 * 1024 * 1024 * 1024
	if bytes > int64(maxBytes) {
		return 0, errors.New("value exceeds 1024 TB")
	}

	return bytes, nil
}
