package tui

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/alert"
	"github.com/IllumiKnowLabs/labstore/cli/tui/infopane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/pane"
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
	"github.com/IllumiKnowLabs/labstore/cli/tui/state"
	"github.com/IllumiKnowLabs/labstore/cli/tui/statusbar"
)

type FocusedPane int

const (
	FocusBuckets FocusedPane = iota
	FocusProfiles
	FocusRemote
	FocusLocal
)

type Model struct {
	AppState state.State

	bucketsPane  pane.BucketsPane
	profilesPane pane.ProfilesPane
	localPane    pane.LocalPane
	remotePane   pane.RemotePane

	profileInfoPane infopane.ProfileInfoPane
	bucketInfoPane  infopane.BucketInfoPane

	statusBar statusbar.Model
	alerts    []alert.Model

	s3BucketProvider providers.S3BucketProvider
	profilesProvider providers.ProfilesProvider
	s3FSProvider     providers.S3FSProvider
	fsProvider       providers.FSProvider

	focusedPane FocusedPane
	width       int
	height      int
}

func New() Model {
	s3BucketProvider := providers.NewS3BucketProvider()
	profilesProvider := providers.NewProfilesProvider()
	s3FSProvider := providers.NewS3FSProvider()
	fsProvider := providers.NewFSProvider()

	bucketPane := pane.NewBuckets(1, "Buckets", pane.WithFocus(), pane.WithSimpleList())
	profilesPane := pane.NewProfiles(2, "Profiles", pane.WithSimpleList())
	remotePane := pane.NewRemote(3, "Remote", pane.WithFileList())
	localPane := pane.NewLocal(4, "Local", pane.WithFileList())

	bucketInfo := infopane.NewBucket("Active Bucket", infopane.ValueNone)
	profileInfo := infopane.NewProfile("Active Profile", infopane.ValueNone)

	m := Model{
		bucketsPane:  bucketPane,
		profilesPane: profilesPane,
		remotePane:   remotePane,
		localPane:    localPane,

		bucketInfoPane:  bucketInfo,
		profileInfoPane: profileInfo,

		statusBar: statusbar.New(DefaultHomeKeyMap.HelpKeys()),

		s3BucketProvider: *s3BucketProvider,
		profilesProvider: *profilesProvider,
		s3FSProvider:     *s3FSProvider,
		fsProvider:       *fsProvider,
	}

	return m
}
