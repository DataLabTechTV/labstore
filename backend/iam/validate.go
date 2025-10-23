package iam

func CheckPolicy(accessKey, bucket, op string) bool {
	if polFunc, ok := Policies[accessKey]; ok {
		return polFunc(bucket, op)
	}
	return false
}
