package messages

import "github.com/IllumiKnowLabs/labstore/cli/types"

type (
	LoadProfilesMsg    struct{}
	ProfilesLoadedMsg  struct{ Entries []types.Entry }
	ProfilesFailedMsg  struct{ Err error }
	ProfileSelectedMsg struct{ Profile string }

	LoadBucketsMsg    struct{}
	BucketsLoadedMsg  struct{ Entries []types.Entry }
	BucketsFailedMsg  struct{ Err error }
	BucketSelectedMsg struct{ Bucket string }

	LoadRemoteMsg   struct{ Dirname *string }
	RemoteLoadedMsg struct {
		Entries []types.Entry
		Active  *string
	}
	RemoteFailedMsg struct{ Err error }

	LoadLocalMsg   struct{ Dirname *string }
	LocalLoadedMsg struct {
		Entries []types.Entry
		Active  *string
	}
	LocalFailedMsg struct{ Err error }

	RefreshAllMsg struct{}

	AlertInfoMsg  struct{ Title, Content string }
	AlertWarnMsg  struct{ Title, Content string }
	AlertErrorMsg struct {
		Title   string
		Content string
		Err     error
	}
	AlertHideMsg struct{ ID string }

	MoveDownMsg     struct{}
	MoveUpMsg       struct{}
	MoveToBottomMsg struct{}
	MoveToTopMsg    struct{}
	PageDownMsg     struct{}
	PageUpMsg       struct{}
	LevelUpMsg      struct{}
	OpenMsg         struct{}
	MarkMsg         struct{}

	StartUploadMsg    struct{}
	UploadProgressMsg struct {
		FileCount int
		FileIndex int
		Uploaded  int64
		Err       error
	}
	UploadFailedMsg struct{ Err error }
	UploadDoneMsg   struct{ FileCount int }

	StartDownloadMsg    struct{}
	DownloadProgressMsg struct {
		FileCount  int
		FileIndex  int
		Downloaded int64
		Err        error
	}
	DownloadFailedMsg struct{ Err error }
	DownloadDoneMsg   struct{}

	StatMsg   struct{}
	DeleteMsg struct{}
)
