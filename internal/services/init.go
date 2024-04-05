package services

func InitServices() {
	initGlobalSharedTicket()
	initBloomFilter()
	initCaptchaClient()
}
