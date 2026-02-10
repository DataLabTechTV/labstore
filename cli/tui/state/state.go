package state

type State struct {
	profile    string
	hasProfile bool

	bucket    string
	hasBucket bool

	localPath    string
	hasLocalPath bool

	remotePath    string
	hasRemotePath bool
}

func (s *State) SetProfile(profile string) {
	s.profile = profile
	s.hasProfile = true
}

func (s *State) UnsetProfile() {
	s.profile = ""
	s.hasProfile = false
}

func (s *State) HasProfile() bool {
	return s.hasProfile
}

func (s *State) Profile() string {
	return s.profile
}

func (s *State) SetBucket(bucket string) {
	s.bucket = bucket
	s.hasBucket = true
}

func (s *State) UnsetBucket() {
	s.bucket = ""
	s.hasBucket = false
}

func (s *State) HasBucket() bool {
	return s.hasBucket
}

func (s *State) Bucket() string {
	return s.bucket
}

func (s *State) SetLocalPath(path string) {
	s.localPath = path
	s.hasLocalPath = true
}

func (s *State) UnsetLocalPath() {
	s.localPath = ""
	s.hasLocalPath = false
}

func (s *State) HasLocalPath() bool {
	return s.hasLocalPath
}

func (s *State) LocalPath() string {
	return s.localPath
}

func (s *State) SetRemotePath(path string) {
	s.remotePath = path
	s.hasRemotePath = true
}

func (s *State) UnsetRemotePath() {
	s.remotePath = ""
	s.hasRemotePath = false
}

func (s *State) HasRemotePath() bool {
	return s.hasRemotePath
}

func (s *State) RemotePath() string {
	return s.remotePath
}
