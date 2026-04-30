package linux

import (
	"testing"
)

func TestDistroString(t *testing.T) {
	d := DistroDebian
	if d.String() != "debian" {
		t.Errorf("Expected debian, got %s", d.String())
	}
}

func TestIsDebianBased(t *testing.T) {
	tests := []struct {
		distro Distro
		want   bool
	}{
		{DistroDebian, true},
		{DistroUbuntu, true},
		{DistroRHEL, false},
		{DistroRocky, false},
	}
	
	for _, tt := range tests {
		if got := tt.distro.IsDebianBased(); got != tt.want {
			t.Errorf("%s.IsDebianBased() = %v, want %v", tt.distro, got, tt.want)
		}
	}
}

func TestIsRedHatBased(t *testing.T) {
	tests := []struct {
		distro Distro
		want   bool
	}{
		{DistroRHEL, true},
		{DistroRocky, true},
		{DistroFedora, true},
		{DistroDebian, false},
	}
	
	for _, tt := range tests {
		if got := tt.distro.IsRedHatBased(); got != tt.want {
			t.Errorf("%s.IsRedHatBased() = %v, want %v", tt.distro, got, tt.want)
		}
	}
}

func TestPackageManager(t *testing.T) {
	tests := []struct {
		distro Distro
		want   string
	}{
		{DistroDebian, "apt"},
		{DistroUbuntu, "apt"},
		{DistroRHEL, "dnf"},
		{DistroRocky, "dnf"},
		{DistroFedora, "dnf"},
		{DistroArch, "pacman"},
		{DistroAlpine, "apk"},
		{DistroSUSE, "zypper"},
	}
	
	for _, tt := range tests {
		if got := tt.distro.PackageManager(); got != tt.want {
			t.Errorf("%s.PackageManager() = %s, want %s", tt.distro, got, tt.want)
		}
	}
}
