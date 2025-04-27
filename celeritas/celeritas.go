package celeritas

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/fouched/celeritas/cache"
	"github.com/fouched/celeritas/render"
	"github.com/fouched/celeritas/session"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "1.0.0"

var redisCache *cache.RedisCache

type Celeritas struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	JetViews      *jet.Set
	config        config // no reason to export this
	EncryptionKey string
	Cache         cache.Cache
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
}

func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := c.Init(pathConfig)
	if err != nil {
		return err
	}

	err = c.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	// read .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers
	infoLog, errorLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errorLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	// connect to database if specified
	if os.Getenv("DATABASE_TYPE") != "" {
		db, err := c.OpenDB(os.Getenv("DATABASE_TYPE"), c.BuildDSN())
		//TODO: build in retry
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		c.DB = Database{
			Type: os.Getenv("DATABASE_TYPE"),
			Pool: db,
		}
	}

	if os.Getenv("CACHE") == "redis" || os.Getenv("SESSION_TYPE") == "redis" {
		redisCache = c.createClientRedisCache()
		c.Cache = redisCache
	}

	// setup config
	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database: databaseConfig{
			dsn:      c.BuildDSN(),
			database: os.Getenv("DATABASE_TYPE"),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}

	// create session
	s := session.Session{
		CookieLifetime: c.config.cookie.lifetime,
		CookiePersist:  c.config.cookie.persist,
		CookieSecure:   c.config.cookie.secure,
		CookieName:     c.config.cookie.name,
		CookieDomain:   c.config.cookie.domain,
		SessionType:    c.config.sessionType,
	}

	switch c.config.sessionType {
	case "redis":
		s.RedisPool = redisCache.Conn
	case "mysql", "mariadb", "postgres", "postgresql":
		s.DBPool = c.DB.Pool
	}

	c.Session = s.InitSession()
	c.EncryptionKey = os.Getenv("KEY")

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)),
		jet.InDevelopmentMode(),
	)
	c.JetViews = views

	c.createRenderer()

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath
	for _, path := range p.folderNames {
		// create folder if it does not exist
		err := c.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}

// ListenAndServe starts the web server
func (c *Celeritas) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer c.DB.Pool.Close()

	c.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	c.ErrorLog.Fatal(err)
}

func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExist(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) createRenderer() {
	myRenderer := render.Render{
		Renderer: c.config.renderer,
		RootPath: c.RootPath,
		Port:     c.config.port,
		JetViews: c.JetViews,
		Session:  c.Session,
	}
	c.Render = &myRenderer
}

func (c *Celeritas) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	}

	return dsn
}

func (c *Celeritas) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Conn:   c.createRedisPool(),
		Prefix: c.config.redis.prefix,
	}
	return &cacheClient
}

func (c *Celeritas) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				c.config.redis.host,
				redis.DialPassword(c.config.redis.password))
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}
