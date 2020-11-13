/*
Package config はconfigファイルを集めたパッケージ
将来的には環境変数から取得するようにする
*/
package config

// Config is a config for connecting DB
type Config struct {
	DB struct {
		Production struct {
			Host     string
			Username string
			Password string
			DBName   string
			Port     string
		}
		Test struct {
			Host     string
			Username string
			Password string
			DBName   string
			Port     string
		}
	}
	Routing struct {
		Port string
	}
}

// NewConfig is a constructor of Config
func NewConfig() *Config {

	c := new(Config)

	c.DB.Production.Host = "mysql-container"
	c.DB.Production.Username = "golang"
	c.DB.Production.Password = "golang"
	c.DB.Production.DBName = "golang"
	c.DB.Production.Port = ":3306"

	c.DB.Test.Host = "test-mysql-container"
	c.DB.Test.Username = "golang"
	c.DB.Test.Password = "golang"
	c.DB.Test.DBName = "golang"
	c.DB.Test.Port = ":3306"

	c.Routing.Port = ":8080"

	return c
}
