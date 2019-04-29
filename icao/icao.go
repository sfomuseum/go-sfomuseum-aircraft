package icao

type AircraftTypes []AircraftType

type AircraftType struct {
	ModelFullName string
	Description string
    WTC string
    Designator string
    ManufacturerCode string
    AircraftDescription string
    EngineCount int 
    EngineType string
}
