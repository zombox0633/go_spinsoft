package station

import (
	"encoding/json"

	"github.com/zombox0633/go_spinsoft/src/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoJSONPointModel struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

type StationModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StationID     int                `bson:"id" json:"station_id"`
	StationCode   int                `bson:"station_code" json:"station_code"`
	Name          string             `bson:"name" json:"name"`
	EnName        string             `bson:"en_name" json:"en_name"`
	ThShort       string             `bson:"th_short" json:"th_short"`
	EnShort       string             `bson:"en_short" json:"en_short"`
	ChName        string             `bson:"chname" json:"chname"`
	ControlDiv    int                `bson:"controldivision" json:"controldivision"`
	ExactKM       int                `bson:"exact_km" json:"exact_km"`
	ExactDistance int                `bson:"exact_distance" json:"exact_distance"`
	KM            int                `bson:"km" json:"km"`
	Class         int                `bson:"class" json:"class"`
	Lat           float64            `bson:"lat" json:"lat"`
	Long          float64            `bson:"long" json:"long"`
	Location      *GeoJSONPointModel `bson:"location,omitempty" json:"coordinates,omitempty"`
	Active        int                `bson:"active" json:"active"`
	Giveway       int                `bson:"giveway" json:"giveway"`
	DualTrack     int                `bson:"dual_track" json:"dual_track"`
	Comment       string             `bson:"comment" json:"comment"`
}

func (s *StationModel) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.StationID = utils.ToInt(raw["id"])
	s.StationCode = utils.ToInt(raw["station_code"])
	s.Name = utils.ToString(raw["name"])
	s.EnName = utils.ToString(raw["en_name"])
	s.ThShort = utils.ToString(raw["th_short"])
	s.EnShort = utils.ToString(raw["en_short"])
	s.ChName = utils.ToString(raw["chname"])
	s.ControlDiv = utils.ToInt(raw["controldivision"])
	s.ExactKM = utils.ToInt(raw["exact_km"])
	s.ExactDistance = utils.ToInt(raw["exact_distance"])
	s.KM = utils.ToInt(raw["km"])
	s.Class = utils.ToInt(raw["class"])
	s.Lat = utils.ToFloat64(raw["lat"])
	s.Long = utils.ToFloat64(raw["long"])
	s.Active = utils.ToInt(raw["active"])
	s.Giveway = utils.ToInt(raw["giveway"])
	s.DualTrack = utils.ToInt(raw["dual_track"])
	s.Comment = utils.ToString(raw["comment"])

	if s.Lat != 0 && s.Long != 0 {
		s.Location = &GeoJSONPointModel{
			Type:        "Point",
			Coordinates: []float64{s.Long, s.Lat}, // x, y
		}
	}

	return nil
}
