package resource

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
)

// Cleanup deletes the resource group created for the sample
func Cleanup(ctx context.Context) error {
	log.Println("deleting resources")
	_, err := DeleteGroup(ctx, azureutil.GetAzureResourceGP())
	return err
}
