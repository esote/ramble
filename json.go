package ramble

type SendReq struct {
	Sender    string   `json:"sender"`
	Recipient []string `json:"recipient"`
	Msg       string   `json:"message"`
	Guid      int      `json:"guid"`
}

type SendResp struct {
	Guid int `json:"guid"`
}

type ViewReq1 struct {
	Fingerprint string `json:"fingerprint"`
}

type ViewResp1 struct {
	NONCE int `json:"NONCE"`
	Guid  int `json:"guid"`
}

type ViewReq2 struct {
	SignedNONCE string `json:"signedNONCE"`
	Guid        int    `json:"guid"`
}

type ViewResp2 struct {
	Guids []int `json:"guids"`
}

type DeleteReq1 ViewReq1

type DeleteResp ViewResp1

type DeleteReq2 ViewReq2
