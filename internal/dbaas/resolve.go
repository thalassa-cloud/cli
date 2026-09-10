package dbaasutil

import (
	"context"
	"fmt"

	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

// BackupStoreRef is a resolved backup store plus the optional source cluster.
type BackupStoreRef struct {
	StoreIdentity string
	Cluster       *dbaas.DbCluster
}

// ResolveBackupStore returns a backup store identity from an explicit store ID
// or from a cluster that has a backup store attached.
func ResolveBackupStore(ctx context.Context, client *dbaas.Client, clusterIdentity, backupStoreIdentity string) (*BackupStoreRef, error) {
	if backupStoreIdentity != "" && clusterIdentity == "" {
		return &BackupStoreRef{StoreIdentity: backupStoreIdentity}, nil
	}

	if clusterIdentity == "" && backupStoreIdentity == "" {
		return nil, fmt.Errorf("cluster identity or backup store identity is required")
	}

	if clusterIdentity == "" {
		return &BackupStoreRef{StoreIdentity: backupStoreIdentity}, nil
	}

	cluster, err := client.GetDbCluster(ctx, clusterIdentity)
	if err != nil {
		if tcclient.IsNotFound(err) {
			return nil, fmt.Errorf("database cluster not found: %s", clusterIdentity)
		}
		return nil, fmt.Errorf("failed to get database cluster: %w", err)
	}

	storeIdentity := backupStoreIdentity
	if storeIdentity == "" {
		if cluster.DbObjectStore == nil || cluster.DbObjectStore.Identity == "" {
			return nil, fmt.Errorf("database cluster %s has no backup store attached", clusterIdentity)
		}
		storeIdentity = cluster.DbObjectStore.Identity
	}

	return &BackupStoreRef{
		StoreIdentity: storeIdentity,
		Cluster:       cluster,
	}, nil
}

// ClusterHasReadyBackupStore reports whether the cluster's backup store is ready.
func ClusterHasReadyBackupStore(cluster *dbaas.DbCluster) bool {
	if cluster == nil || cluster.DbObjectStore == nil {
		return false
	}
	return cluster.DbObjectStore.Status == dbaas.ObjectStatusReady
}

// PITRAvailableForCluster reports whether the cluster currently advertises a PITR window.
func PITRAvailableForCluster(cluster *dbaas.DbCluster) bool {
	if cluster == nil || cluster.BackupRecoveryWindow == nil {
		return false
	}
	storeReady := ClusterHasReadyBackupStore(cluster)
	return cluster.BackupRecoveryWindow.PitrAvailable(storeReady, storeReady)
}
