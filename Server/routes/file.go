package routes

import (
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"
	"upload-and-download-file/models"

	"github.com/gin-gonic/gin"
)

// File upload API
func (s *Server) Upload(c *gin.Context) {
	PrettyPrint("Start to Upload ...")

	// user, _ := c.Get("user")
	username := c.GetString("User")
	file, header, err := c.Request.FormFile("file")

	PrettyPrint(("file :"))
	PrettyPrint(file)
	PrettyPrint(("header :"))
	PrettyPrint(header)

	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid file"})
		return
	}

	// Convert the uploaded file to base64
	base64File, err := s.ConvertFileToBase64(file)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to convert file to base64"})
		return
	}

	// Store the base64-encoded file in the PostgreSQL database
	fileDetails, err := s.StoreBase64File(base64File, header, username)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to store file in the database"})
		return
	}

	c.JSON(200, gin.H{
		"message":     "File uploaded and stored successfully",
		"fileDetails": fileDetails,
	})
}

// File download API with download limit check
func (s *Server) Download(c *gin.Context) {
	URL := c.Param("URL")
	// var data map[string]string

	// if err := c.ShouldBindJSON(&data); err != nil {
	// 	c.JSON(400, gin.H{"message": "Bad Request"})
	// 	return
	// }

	// URL := data["URL"]

	PrettyPrint("Start to Download " + URL + " ...")

	// Check download limit
	var file models.StroeData

	// Perform download limit checking (e.g., by tracking downloads in the database)
	if err := s.gd.GetCorresponding(&file, "download_url = ?", URL); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Decode the Base64-encoded file
	decoded, err := base64.StdEncoding.DecodeString(file.FileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode Base64"})
		return
	}

	// Set the appropriate headers for file download
	// c.Header("Content-Description", "File Transfer")
	// c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Disposition", "attachment; filename="+file.FileName)
	// c.Header("Content-Type", "application/octet-stream")

	// Determine the user's home directory based on the operating system
	homeDir := ""
	if runtime.GOOS == "windows" {
		homeDir = "C:\\Users\\" + os.Getenv("USERPROFILE") + "\\Downloads"
	} else if runtime.GOOS == "darwin" {
		homeDir = "/Users/" + os.Getenv("USER") + "/Downloads"
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported operating system"})
		return
	}

	// Combine the home directory and file name to create the local file path
	localFilePath := homeDir + string(os.PathSeparator) + file.FileName

	PrettyPrint(localFilePath)
	PrettyPrint(localFilePath)
	PrettyPrint(localFilePath)

	// Create and write the decoded file to the local directory
	if err := writeToFile(localFilePath, decoded); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		PrettyPrint("Failed to save the file")
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "File downloaded and stored successfully"})
		PrettyPrint("File downloaded and stored successfully")
	}
}

// ConvertFileToBase64 converts a file to base64
func (s *Server) ConvertFileToBase64(file io.Reader) (string, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	base64File := base64.StdEncoding.EncodeToString(fileContent)

	return base64File, nil
}

// StoreBase64File stores the base64-encoded file in the database
func (s *Server) StoreBase64File(base64File string, header *multipart.FileHeader, username string) (*models.StroeData, error) {

	// TODO: Request expired time and download limit from frontend...
	uploadTime := time.Now()
	shareTime := time.Now().Add(time.Hour * 24)
	shareLimit := 20

	fileDetails := &models.StroeData{
		UploadTime:  uploadTime,
		ShareTime:   shareTime,  // Set to 24 hours from now
		ShareLimit:  shareLimit, // Initial shares to 0
		FileSize:    int64(header.Size),
		FileName:    header.Filename,
		FileType:    header.Header.Get("Content-Type"),
		FileContent: base64File,
		DownloadUrl: header.Filename + username,
	}

	// TODO: Generate the file URL using the file ID

	// Insert base64 file content into the database
	if err := s.gd.GetCorresponding(fileDetails, "name = ?", fileDetails.FileName); err != nil {
		if err := s.gd.Create(fileDetails); err != nil {
			panic(err)
		}
	} else {
		if err := s.gd.Update(fileDetails); err != nil {
			panic(err)
		}
	}

	fileDetails.FileContent = ""

	PrettyPrint(fileDetails)

	return fileDetails, nil
}

func (s *Server) CheckDownloadCount(c *gin.Context) {
	// In a real application, retrieve the download count from the database
	URL := c.Param("URL")

	// Check download limit
	var file models.StroeData

	// Perform download limit checking (e.g., by tracking downloads in the database)
	if err := s.gd.GetCorresponding(&file, "download_url = ?", URL); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// If the download limit has been reached, deny access
	if file.ShareLimit < 1 {
		c.JSON(401, gin.H{"message": "Download limit reached"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Download limit remaining"})
	}
}

func (s *Server) UpdateDownloadCount(c *gin.Context) {
	URL := c.Param("URL")

	var file models.StroeData

	if err := s.gd.GetCorresponding(&file, "download_url = ?", URL); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	file.ShareLimit = file.ShareLimit - 1

	if err := s.gd.Upsert(&file); err != nil {
		c.JSON(401, gin.H{"message": "failed to update download count"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Update download count"})
	}
}

// Helper function to write a byte slice to a local file
func writeToFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (s *Server) DeleteData(c *gin.Context) {
	var dataitem models.DataItem
	var data []models.StroeData

	if err := c.ShouldBindJSON(&dataitem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := s.gd.GetCorresponding(&data, "download_url = ?", dataitem.DownloadURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Id existed in database"})
	}

	if reflect.ValueOf(data).IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"message": "User data does not exist!"})
		return
	}

	// Delete data by their IDs
	if err := s.gd.Delete(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data deleted successfully"})
}

func (s *Server) DeleteDatas(c *gin.Context) {
	var frontend_ids models.FrontendRequest
	var data []models.StroeData

	if err := c.ShouldBindJSON(&frontend_ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(frontend_ids.FrontendIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided for deletion"})
		return
	}

	for _, fronted_id := range frontend_ids.FrontendIDs {

		if err := s.gd.GetCorresponding(&data, fronted_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Id existed in database"})
		}

		if reflect.ValueOf(data).IsZero() {
			c.JSON(http.StatusNotFound, gin.H{"message": "User data does not exist!"})
			return
		}

		// Delete data by their IDs
		if err := s.gd.Delete(&data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "data deleted successfully"})

	}
}

func (s *Server) SearchAllData(c *gin.Context) {

	var storedata []models.StroeData

	if err := s.gd.GetCorresponding(&storedata, "1 = ?", "1"); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, storedata)
}

func (s *Server) UserSearchAllData(c *gin.Context) {

	var storedata []models.StroeData
	var dataitem []models.DataItem

	if err := s.gd.GetCorresponding(&storedata, "1 = ?", "1"); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}

	for _, data := range storedata {

		item := models.DataItem{
			Dataname:    data.FileName,
			DownloadURL: data.DownloadUrl,
		}

		dataitem = append(dataitem, item)
	}

	PrettyPrint(dataitem)

	c.JSON(http.StatusOK, dataitem)
}

// A background goroutine to periodically delete expired files
// func (s *Server) DeleteExpiredFiles() {
// 	for {
// 		var files []models.StroeData

// 		currentTime := time.Now()

// 		s.gd.GetCorresponding(&files, "share_time <= ?", currentTime)

// 		s.gd.Delete(&files)

// 		// Check every 24 hours
// 		time.Sleep(24 * time.Minute)
// 	}
// }
