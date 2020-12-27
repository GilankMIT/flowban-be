package fileUploader

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"os"
	"tools_kecamatan/helper/stringBuilder"
)

func FormFileUploader(c *gin.Context, formName string, destinationFolder string,
	additionalPrefix string) (string, error) {
	_, fileHeader, err := c.Request.FormFile(formName)
	if err != nil {
		log.Error().Msg("cannot retrieve " + formName + err.Error())
		return "", err
	} else {
		if fileHeader.Size > 0 { //check if file size is bigger than 0

			destinationFile, err := saveFile(fileHeader, destinationFolder, additionalPrefix)
			if err != nil {
				return "", err
			}

			return destinationFile, nil
		}

		return "", errors.New("cannot retrieve file : file size is 0")
	}
}

func MultiFileUploader(c *gin.Context, formName string, destinationFolder, additionalPrefix string) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	var uploadedFilePathList []string

	//retrieve files
	files := form.File[formName]

	for _, file := range files {

		filePath, err := saveFile(file, destinationFolder, additionalPrefix)
		if err != nil {
			return nil, err
		}
		uploadedFilePathList = append(uploadedFilePathList, filePath)
	}

	return uploadedFilePathList, nil
}

func saveFile(fileHeader *multipart.FileHeader, destinationFolder, additionalPrefix string) (filePath string, err error) {
	//get file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	//build destination path
	uniqueSuffix := stringBuilder.GenerateRandom(12)
	destinationFile := destinationFolder + "/" + additionalPrefix + uniqueSuffix + "_" + fileHeader.Filename

	//create file and copy
	out, err := os.Create(destinationFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return destinationFile, nil
}
