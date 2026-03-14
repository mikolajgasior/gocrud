package builder

// Options configuration for validation and field generation
type Options struct {
	RestrictFields  map[string]struct{}
	ExcludeFields   map[string]struct{}
	TagName         string
	IDPrefix        string
	NamePrefix      string
	OverwriteValues map[string]string
	Values          bool
}
