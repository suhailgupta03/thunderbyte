package common

import S3Uploader "github.com/suhailgupta03/go-s3-uploader"

// UploadHandler Default upload handler for thunderbyte
// Checks for the form field with name="file". Returns an error
// if it does not find that field
func UploadHandler(context AppContext, s3Config S3Uploader.S3) (*S3Uploader.UploadId, error) {
	file, err := context.HTTPServerContext.FormFile("file")
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileName := file.Filename
	context.Logger.Info("Attempting to upload file", "filename", fileName)
	var fileBytes []byte
	src.Read(fileBytes)
	uploadId, err := s3Config.UploadFile(fileBytes, fileName)
	if err != nil {
		return nil, err
	}
	return uploadId, nil
}
