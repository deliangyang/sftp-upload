package ssh

type User struct {
	Name    string `json:"name"`
	RSAPath string `json:"rsaPath"`
	Client  string `json:"client"`
	Dir     string `json:"dir"`
	Chmod   string `json:"chmod"`
	Watch   string `json:"watch"`
}