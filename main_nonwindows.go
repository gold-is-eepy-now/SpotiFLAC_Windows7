//go:build !windows

package main

func shouldUseLegacyWindowStyle() bool {
	return false
}
