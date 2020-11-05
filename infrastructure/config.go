/*
infrastructure is Frameworks & Drivers.
DBのドライバー,Ginをinterfacesとつなぐ．
interfacesに依存．
*/
package infrastructure

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
