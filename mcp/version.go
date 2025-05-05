package mcp

type Version string

const (
	DefaultVersion          = Experimental
	Experimental    Version = "experimental"
	Latest          Version = Version20250326
	Version20250326 Version = "2025-03-26"
	Version20241105 Version = "2024-11-05"
)
