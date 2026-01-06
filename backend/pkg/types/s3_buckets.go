package types

import "encoding/xml"

type Bucket struct {
	Name         string
	CreationDate Timestamp
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   Owner
	Buckets struct {
		Bucket []Bucket
	}
}

type BaseListBucketResult struct {
	XMLName        xml.Name `xml:"ListBucketResult"`
	Name           string
	Prefix         string
	MaxKeys        int
	Contents       []Object
	CommonPrefixes []CommonPrefixes
	IsTruncated    bool
	UntilKey       string `xml:"-"`
}

type ListBucketResult struct {
	BaseListBucketResult
	Marker     string
	NextMarker string
}

type ListBucketResultV2 struct {
	BaseListBucketResult
	KeyCount              int
	ContinuationToken     string
	NextContinuationToken string
	StartAfter            string
}

type CommonPrefixes struct {
	Prefix string
}
