package vo

type PublicKeyEncryptVo struct {
	EncryptedData string `json:"data"`
}

type UploadImageVo struct {
	ImgUrl string `json:"img_url"`
}

type RASPublicKeyVo struct {
	PublicKey string `json:"public_key"`
}
