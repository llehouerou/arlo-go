package arlo

// A Device is the device data, this can be a camera, basestation, arloq, etc.
type Device struct {
	arlo                          *Arlo        // Let's hold a reference to the parent arlo object since it holds the http.Client object and references to all devices.
	AnalyticsEnabled              bool         `json:"analyticsEnabled"`
	ArloMobilePlan                bool         `json:"arloMobilePlan"`
	ArloMobilePlanId              string       `json:"arloMobilePlanId"`
	ArloMobilePlanName            string       `json:"arloMobilePlanName"`
	ArloMobilePlanThreshold       int          `json:"arloMobilePlanThreshold"`
	Connectivity                  Connectivity `json:"connectivity"`
	CriticalBatteryState          bool         `json:"criticalBatteryState"`
	DateCreated                   float64      `json:"dateCreated"`
	DeviceId                      string       `json:"deviceId"`
	DeviceName                    string       `json:"deviceName"`
	DeviceType                    string       `json:"deviceType"`
	DisplayOrder                  uint8        `json:"displayOrder"`
	FirmwareVersion               string       `json:"firmwareVersion"`
	InterfaceVersion              string       `json:"interfaceVersion"`
	InterfaceSchemaVer            string       `json:"interfaceSchemaVer"`
	LastImageUploaded             string       `json:"lastImageUploaded"`
	LastModified                  float64      `json:"lastModified"`
	MigrateActivityZone           bool         `json:"migrateActivityZone"`
	MobileCarrier                 string       `json:"mobileCarrier"`
	MobileTrialUsed               bool         `json:"mobileTrialUsed"`
	PermissionsFilePath           string       `json:"permissionsFilePath"`
	PermissionsSchemaVer          string       `json:"permissionsSchemaVer"`
	PermissionsVerison            string       `json:"permissionsVerison"` // WTF? Netgear developers think this is OK... *sigh*
	PermissionsVersion            string       `json:"permissionsVersion"`
	PresignedFullFrameSnapshotUrl string       `json:"presignedFullFrameSnapshotUrl"`
	PresignedLastImageUrl         string       `json:"presignedLastImageUrl"`
	PresignedSnapshotUrl          string       `json:"presignedSnapshotUrl"`
	MediaObjectCount              uint8        `json:"mediaObjectCount"`
	ModelId                       string       `json:"modelId"`
	Owner                         Owner        `json:"owner"`
	ParentId                      string       `json:"parentId"`
	Properties                    Properties   `json:"properties"`
	UniqueId                      string       `json:"uniqueId"`
	UserId                        string       `json:"userId"`
	UserRole                      string       `json:"userRole"`
	State                         string       `json:"state"`
	XCloudId                      string       `json:"xCloudId"`
}

// Devices is an array of Device objects.
type Devices []Device

// A DeviceOrder holds a map of device ids and a numeric index. The numeric index is the device order.
// Device order is mainly used by the UI to determine which order to show the devices.
/*
{
  "devices":{
    "XXXXXXXXXXXXX":1,
    "XXXXXXXXXXXXX":2,
    "XXXXXXXXXXXXX":3
}
*/
type DeviceOrder struct {
	Devices map[string]int `json:"devices"`
}

// Find returns a device with the device id passed in.
func (ds *Devices) Find(deviceId string) *Device {
	for _, d := range *ds {
		if d.DeviceId == deviceId {
			return &d
		}
	}

	return nil
}

func (ds *Devices) FindCameras(basestationId string) Cameras {
	cs := new(Cameras)
	for _, d := range *ds {
		if d.ParentId == basestationId {
			*cs = append(*cs, Camera(d))
		}
	}

	return *cs
}

func (d Device) IsBasestation() bool {
	return d.DeviceType == DeviceTypeBasestation || d.DeviceId == d.ParentId
}

func (d Device) IsCamera() bool {
	return d.DeviceType == DeviceTypeCamera
}

// GetBasestations returns a Basestations object containing all devices that are NOT type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Cameras also includes devices of this type, so you can get the same data there or cast.
func (ds Devices) GetBasestations() Basestations {
	var basestations Basestations
	for _, d := range ds {
		if d.IsBasestation() || !d.IsCamera() {
			basestations = append(basestations, Basestation{Device: d})
		}
	}
	return basestations
}

// GetCameras returns a Cameras object containing all devices that are of type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Basestations also includes devices of this type, so you can get the same data there or cast.
func (ds Devices) GetCameras() Cameras {
	var cameras Cameras
	for _, d := range ds {
		if d.IsCamera() || !d.IsBasestation() {
			cameras = append(cameras, Camera(d))
		}
	}
	return cameras
}

// UpdateDeviceName sets the name of the given device to the name argument.
func (d *Device) UpdateDeviceName(name string) error {

	body := map[string]string{"deviceId": d.DeviceId, "deviceName": name, "parentId": d.ParentId}
	resp, err := d.arlo.put(DeviceRenameUri, d.XCloudId, body, nil)
	return checkRequest(*resp, err, "failed to update device name")
}
