package validationFunc

import (
	"errors"
	"fullstack-journal/app/filters"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

func VEmail(c echo.Context, email string) error {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New(
			filters.Translate(c, "emailNotValid"),
		)
	}
	return nil
}

func VPhoneNumber(c echo.Context, phoneNumber string) error{
	if phoneNumber == "" {
		return errors.New(
			filters.Translate(c, "phoneNumberRequired"),
		)
	}

	match, _ := regexp.MatchString(`^\d{8,}$`, phoneNumber)
	if !match {
		return errors.New(
			filters.Translate(c, "phoneNumberNumeric"),
		)
	}

	return nil
}
func VFile(extAllowed []string, file *multipart.FileHeader, maxSizeInMb int, c echo.Context)error{
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range extAllowed{
		if ext == strings.ToLower(allowedExt){
			allowed = true
			break
		}
	}

	if !allowed{
		return errors.New(
			filters.Translate(c, "extFileNotValid"),
		)
	}

	if file.Size > int64(maxSizeInMb) * 1024 * 1024{
		return errors.New(
			filters.Translate(c, "maxFile")+" "+string(maxSizeInMb),
		)

	}

	return nil
}