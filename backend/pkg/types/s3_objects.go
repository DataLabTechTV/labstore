package types

import (
	"encoding/xml"
	"io"
	"log/slog"
	"time"

	"github.com/IllumiKnowLabs/labstore/backend/internal/errs"
)

type BaseObject struct {
	Key          string
	ETag         string
	LastModified Timestamp
	Size         int64
}

type Object struct {
	BaseObject
	ChecksumAlgorithm []string
	ChecksumType      string
	Owner             *Owner
	RestoreStatus     RestoreStatus
	StorageClass      string
}

type ObjectIdentifier struct {
	BaseObject
	VersionId string
}

type RestoreStatus struct {
	IsRestoreInProgress bool
	RestoreExpiryDate   Timestamp
}

type GetObjectResult struct {
	Content      io.ReadSeekCloser
	ObjectSize   int
	DateModified time.Time
}

type DeleteObjectsRequest struct {
	XMLName xml.Name `xml:"Delete"`
	Object  []ObjectIdentifier
	Quiet   bool
}

type DeleteResult struct {
	Deleted []DeletedObject
	Error   []errs.S3Error
}

type DeletedObject struct {
	DeleteMarker          bool
	DeleteMarkerVersionId string
	Key                   string
	VersionId             string
}

func (req DeleteObjectsRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("XMLName", req.XMLName),
		slog.Int("Objects", len(req.Object)),
		slog.Bool("Quiet", req.Quiet),
	)
}
