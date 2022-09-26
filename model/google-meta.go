package model

func AllGoogleMetas() (metas []string, err error) {
	var gMetas []GoogleMeta
	err = PgsqlDB.Find(&gMetas).Error
	if err != nil {
		return nil, err
	}

	metas = make([]string, len(gMetas))
	for i, g := range gMetas {
		metas[i] = g.Content
	}
	return
}
