package server

func checkPolicy(accessKey, bucket, op string) bool {
	if polFunc, ok := userPolicies[accessKey]; ok {
		return polFunc(bucket, op)
	}
	return false
}
