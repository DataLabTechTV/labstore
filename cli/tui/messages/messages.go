package messages

import (
	"github.com/IllumiKnowLabs/labstore/cli/tui/providers"
)

type (
	LoadProfilesMsg    struct{}
	ProfilesLoadedMsg  struct{ Entries []providers.Entry }
	ProfilesFailedMsg  struct{ Err error }
	ProfileSelectedMsg struct{ Profile string }

	LoadBucketsMsg    struct{}
	BucketsLoadedMsg  struct{ Entries []providers.Entry }
	BucketsFailedMsg  struct{ Err error }
	BucketSelectedMsg struct{ Bucket string }

	LoadRemoteMsg   struct{ Dirname *string }
	RemoteLoadedMsg struct {
		Entries []providers.Entry
		Active  *string
	}
	RemoteFailedMsg struct{ Err error }

	LoadLocalMsg   struct{ Dirname *string }
	LocalLoadedMsg struct {
		Entries []providers.Entry
		Active  *string
	}
	LocalFailedMsg struct{ Err error }

	RefreshAllMsg struct{}

	AlertInfoMsg  struct{ Title, Content string }
	AlertWarnMsg  struct{ Title, Content string }
	AlertErrorMsg struct{ Err error }
	AlertHideMsg  struct{ ID string }

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{}
	OpenMsg         struct{}
	MarkMsg         struct{}

	UploadMsg   struct{}
	DownloadMsg struct{}
	StatMsg     struct{}
	DeleteMsg   struct{}
)
