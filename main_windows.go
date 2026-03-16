//go:build windows

package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type rtlOSVersionInfoEx struct {
	OSVersionInfoSize uint32
	MajorVersion      uint32
	MinorVersion      uint32
	BuildNumber       uint32
	PlatformID        uint32
	CSDVersion        [128]uint16
	ServicePackMajor  uint16
	ServicePackMinor  uint16
	SuiteMask         uint16
	ProductType       byte
	Reserved          byte
}

func shouldUseLegacyWindowStyle() bool {
	ver, err := getWindowsVersion()
	if err != nil {
		return false
	}

	// Windows 7 is NT 6.1.
	return ver.MajorVersion < 6 || (ver.MajorVersion == 6 && ver.MinorVersion <= 1)
}

func getWindowsVersion() (*rtlOSVersionInfoEx, error) {
	dll := windows.NewLazySystemDLL("ntdll.dll")
	proc := dll.NewProc("RtlGetVersion")

	if err := dll.Load(); err != nil {
		return nil, fmt.Errorf("load ntdll.dll: %w", err)
	}
	if err := proc.Find(); err != nil {
		return nil, fmt.Errorf("find RtlGetVersion: %w", err)
	}

	info := &rtlOSVersionInfoEx{OSVersionInfoSize: uint32(unsafe.Sizeof(rtlOSVersionInfoEx{}))}
	r1, _, callErr := proc.Call(uintptr(unsafe.Pointer(info)))
	if r1 != 0 {
		return nil, fmt.Errorf("RtlGetVersion failed: %w", callErr)
	}

	return info, nil
}
