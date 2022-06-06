package models

// import "time"

type ResponseProfil struct {
	ID            string     `json:"id"`
	AppName       string     `json:"appName"`
	MessageProfil string     `json:"message"`
	Version       string     `json:"version"`
	Datas         DataDetail `json:"data"`
}

type DataDetail struct {
	ID              uint         `json:"id"`
	Firstname             string      `json:"firstname"`
	Lastname              string      `json:"lastname"`
	Email                 string      `json:"email"`
	Phone                 string      `json:"phone"`
	DateOfBirth           string      `json:"dateOfBirth"`
	IdentityNumber        string      `json:"identityNumber"`
	MemberStatus          string      `json:"memberStatus"`
	MemberExpiredDate     string      `json:"memberExpiredDate"`
	MemberPoint           string      `json:"memberPoint"`
	ExpiryPoint           int         `json:"expiryPoint"`
	ExpiryDate            string      `json:"expiryDate"`
	TotalCoupons          string      `json:"totalCoupons"`
	Gender                string      `json:"gender"`
	IsSubscribeNewsletter string      `json:"isSubscribeNewsletter"`
	CustomerGroup         string      `json:"customerGroup"`
	Attributes            Attributess `json:"attributes"`
}

type Attributess struct {
	AboutYou         AboutYous         `json:"aboutYou"`
	SocialMedia      SocialMedias      `json:"socialMedia"`
	CompletionStatus CompletionStatuss `json:"completionStatus"`
}

type AboutYous struct {
	Status          string `json:"status"`
	Job             string `json:"job"`
	Hobby           string `json:"hobby"`
	LifestyleBudget string `json:"lifestyleBudget"`
	Celebrate       string `json:"celebrate"`
}

type SocialMedias struct {
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
}

type CompletionStatuss struct {
	PrivacyData string `json:"privacyData"`
	SocialMedia string `json:"socialMedia"`
	AboutYou    string `json:"aboutYou"`
	Address     string `json:"address"`
}
