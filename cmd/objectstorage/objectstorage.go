package objectstorage

import (
	"github.com/spf13/cobra"
)

// ObjectStorageCmd represents the object storage command
var ObjectStorageCmd = &cobra.Command{
	Use:     "object-storage",
	Aliases: []string{"objectstorage", "os", "s3", "buckets"},
	Short:   "Manage object storage buckets",
}

func init() {
}

