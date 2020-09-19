package main

func main() {
	//var sqlDb *gorm.DB
	//var mongoClient *mongo.Client
	//var mongoDatabase *mongo.Database
	//var closer io.Closer
	//
	//cfg, err := config.NewConfig()
	//if err != nil {
	//	log.Fatalf("Error while creating app configuration: %v\n", err)
	//}
	//
	//ctx, cancel := context.WithCancel(context.Background())
	//
	//if cfg.StoreInSql {
	//	sqlDb, closer = connectToRelationalDatabaseAndMigrate(cfg.Database, cfg.DatabaseDSN)
	//	defer closer.Close()
	//} else {
	//	mongoClient, closer = connectToMongoAndMigrate(ctx, cfg.Mongo)
	//	mongoDatabase = mongoClient.Database(db.MongoDBName)
	//	defer closer.Close()
	//}
	//
	//app := fiber.New(fiber.Config{
	//	Prefork:      cfg.UsePrefork,
	//	ErrorHandler: handlers.Error,
	//})
	//
	//registerAPIHandlers(ctx, &cfg, mongoDatabase, sqlDb, app)
	//
	//if err := app.Listen(cfg.Port); err != nil {
	//	log.Fatalf("Error while starting Fiber Server: %v", err)
	//}
	//cancel()
}
