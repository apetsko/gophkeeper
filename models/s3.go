package models

// S3UploadData contains metadata for uploading an object to S3/MinIO.
//
// Fields:
//   - ObjectName: The name of the object in the bucket.
//   - MetaContent: Additional metadata content.
//   - FileName: The original file name.
//   - FileType: The MIME type of the file.
type S3UploadData struct {
	ObjectName  string
	MetaContent string
	FileName    string
	FileType    string
}
