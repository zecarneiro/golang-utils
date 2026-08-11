package common

func WithAppId(id string) {
	appId = id
}

func GetAppId() string {
	if len(appId) == 0 {
		appId = "GoApp"
	}
	return appId
}
