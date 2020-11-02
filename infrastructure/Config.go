package infrastructure

type Config struct {
    DB struct {
        Production struct {
            Host string
            Username string
            Password string
            DBName string
        }
        Test struct {
            Host string
            Username string
            Password string
            DBName string
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

    c.DB.Test.Host = "test-mysql-container"
    c.DB.Test.Username = "golang"
    c.DB.Test.Password = "golang"
    c.DB.Test.DBName = "golang"
    
    c.Routing.Port = ":3306"

    return c
}
