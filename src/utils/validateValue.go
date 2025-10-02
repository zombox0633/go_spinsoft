package utils

import "fmt"

// ---------------------------------- Latitude and Longitude -------------------------

func ValidateCoordinates(lat, long float64) error {
	if err := validateLatitude(lat); err != nil {
		return err
	}
	if err := validateLongitude(long); err != nil {
		return err
	}
	return nil
}

func validateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	return nil
}

func validateLongitude(long float64) error {
	if long < -180 || long > 180 {
		return fmt.Errorf("invalid longitude: must be between -180 and 180")
	}
	return nil
}
