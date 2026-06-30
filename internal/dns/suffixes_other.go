//go:build !windows

package dns

import "context"

func SearchSuffixes(ctx context.Context) ([]string, error) {
	return nil, nil
}
