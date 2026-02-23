package main

import "athenamind/internal/snapshot"

func createSnapshot(root, createdBy, reason string) (snapshotManifest, error) {
	return snapshot.CreateSnapshot(root, createdBy, reason)
}

func listSnapshots(root string) ([]snapshotListRow, error) {
	return snapshot.ListSnapshots(root)
}

func restoreSnapshot(root, snapshotID string) error {
	return snapshot.RestoreSnapshot(root, snapshotID)
}

func loadSnapshotManifest(root, snapshotID string) (snapshotManifest, error) {
	return snapshot.LoadSnapshotManifest(root, snapshotID)
}

func validateSnapshotManifest(m snapshotManifest) error {
	return snapshot.ValidateSnapshotManifest(m)
}

func writeSnapshotAudit(root string, ev snapshotAuditEvent) error {
	return snapshot.WriteSnapshotAudit(root, ev)
}
