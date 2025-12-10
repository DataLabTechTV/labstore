package iam

type Store struct {
	Users    map[string]*User
	Groups   map[string]*Group
	Policies map[string]*Policy
}

func ensureSchema() {

}
