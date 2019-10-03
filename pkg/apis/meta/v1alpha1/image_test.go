package v1alpha1

import (
	"testing"
)

func TestNewOCIImageRef(t *testing.T) {
	tests := []struct {
		in, out string
		err     bool
	}{
		{
			in:  "weaveworks/ignite-kernel:4.19.47",
			out: "weaveworks/ignite-kernel:4.19.47",
		},
		{
			in:  "weaveworks/ignite-ubuntu:v0.6.0",
			out: "weaveworks/ignite-ubuntu:v0.6.0",
		},
		{
			in:  "centos",
			out: "centos:latest",
		},
		{
			in:  "weaveworks/ignite-ubuntu@sha256:ad984fa5f6f2db55a0a48860d263a02a0f77aee3bbefa0d71648b4bc287ac13c",
			out: "weaveworks/ignite-ubuntu@sha256:ad984fa5f6f2db55a0a48860d263a02a0f77aee3bbefa0d71648b4bc287ac13c",
		},
		{
			in:  "weaveworks/ignite-kubeadm@sha256:928f55c6162a50d3278e334ac0d7a6e0bf7c765d9c6dd835b935f0ca73e398ca",
			out: "weaveworks/ignite-kubeadm@sha256:928f55c6162a50d3278e334ac0d7a6e0bf7c765d9c6dd835b935f0ca73e398ca",
		},
		{
			in:  "skjjnfnskj//bs::777",
			err: true,
		},
	}

	for _, rt := range tests {
		t.Run(rt.in, func(t *testing.T) {
			actual, err := NewOCIImageRef(rt.in)
			if (err != nil) != rt.err {
				t.Errorf("\nexpected error: %t\nactual:  %v", rt.err, err)
			}

			if actual.String() != rt.out {
				t.Errorf("\nexpected %q\nactual:  %q", rt.out, actual.String())
			}
		})
	}
}

func TestParseOCIString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  OCIContentID
		err  bool
	}{
		{
			name: "oci string with repo name",
			in:   "oci://docker.io/library/hello-world@sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
			out: OCIContentID{
				repoName: "docker.io/library/hello-world",
				digest:   "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
			},
		},
		{
			name: "oci string local docker sha",
			in:   "docker://sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
			out: OCIContentID{
				digest: "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
			},
		},
		{
			name: "invalid oci string",
			in:   "sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e",
			err:  true,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			actual, err := parseOCIString(rt.in)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			if !rt.err {
				if actual.String() != rt.out.String() {
					t.Errorf("expected %q, actual: %q", rt.out.String(), actual.String())
				}
			}
		})
	}
}
