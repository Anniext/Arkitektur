package profile

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var (
	DefaultProfile *Profile
)

// InitSystemProfile 初始化配置文件
func InitSystemProfile(name string) error {
	profile := &Profile{}
	err := viper.Unmarshal(profile)
	profile.Name = name
	if err != nil {
		return err
	}

	if profile.Mode != "demo" && profile.Mode != "dev" && profile.Mode != "prod" {
		profile.Mode = "dev"
	}

	if profile.Mode == "prod" && profile.Tmp == "" {
		if runtime.GOOS == "windows" {
			profile.Config = filepath.Join(os.Getenv("ProgramData"), name)

			if _, err := os.Stat(profile.Tmp); os.IsNotExist(err) {
				if err := os.MkdirAll(profile.Tmp, 0770); err != nil {
					fmt.Printf("Failed to create data directory: %s, err: %+v\n", profile.Tmp, err)
					return err
				}
			}
		} else {
			profile.Tmp = "/var/opt/" + name
		}
	}

	DefaultProfile = profile

	return nil
}

func GetDefaultProfile() *Profile {
	return DefaultProfile
}
