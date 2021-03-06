package image

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/docker/cli/internal/test"
	"github.com/docker/cli/internal/test/testutil"
	"github.com/docker/docker/api/types"
	"github.com/gotestyourself/gotestyourself/golden"
	"github.com/stretchr/testify/assert"
)

func TestNewInspectCommandErrors(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "wrong-args",
			args:          []string{},
			expectedError: "requires at least 1 argument.",
		},
	}
	for _, tc := range testCases {
		cmd := newInspectCommand(test.NewFakeCli(&fakeClient{}))
		cmd.SetOutput(ioutil.Discard)
		cmd.SetArgs(tc.args)
		testutil.ErrorContains(t, cmd.Execute(), tc.expectedError)
	}
}

func TestNewInspectCommandSuccess(t *testing.T) {
	imageInspectInvocationCount := 0
	testCases := []struct {
		name             string
		args             []string
		imageCount       int
		imageInspectFunc func(image string) (types.ImageInspect, []byte, error)
	}{
		{
			name:       "simple",
			args:       []string{"image"},
			imageCount: 1,
			imageInspectFunc: func(image string) (types.ImageInspect, []byte, error) {
				imageInspectInvocationCount++
				assert.Equal(t, "image", image)
				return types.ImageInspect{}, nil, nil
			},
		},
		{
			name:       "format",
			imageCount: 1,
			args:       []string{"--format='{{.ID}}'", "image"},
			imageInspectFunc: func(image string) (types.ImageInspect, []byte, error) {
				imageInspectInvocationCount++
				return types.ImageInspect{ID: image}, nil, nil
			},
		},
		{
			name:       "simple-many",
			args:       []string{"image1", "image2"},
			imageCount: 2,
			imageInspectFunc: func(image string) (types.ImageInspect, []byte, error) {
				imageInspectInvocationCount++
				if imageInspectInvocationCount == 1 {
					assert.Equal(t, "image1", image)
				} else {
					assert.Equal(t, "image2", image)
				}
				return types.ImageInspect{}, nil, nil
			},
		},
	}
	for _, tc := range testCases {
		imageInspectInvocationCount = 0
		cli := test.NewFakeCli(&fakeClient{imageInspectFunc: tc.imageInspectFunc})
		cmd := newInspectCommand(cli)
		cmd.SetOutput(ioutil.Discard)
		cmd.SetArgs(tc.args)
		err := cmd.Execute()
		assert.NoError(t, err)
		golden.Assert(t, cli.OutBuffer().String(), fmt.Sprintf("inspect-command-success.%s.golden", tc.name))
		assert.Equal(t, imageInspectInvocationCount, tc.imageCount)
	}
}
