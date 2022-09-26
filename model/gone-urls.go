package model

func GoneURLs() (urls []string, err error) {
	var goneURLS []GoneURL
	err = PgsqlDB.Find(&goneURLS).Error
	if err != nil {
		return nil, err
	}

	urls = make([]string, len(goneURLS))
	for i, g := range goneURLS {
		urls[i] = g.URL
	}
	return
}
