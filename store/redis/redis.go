package redis

// // NewClient creates aredis connection
// func NewClient(opts *redis.Options) (*redis.Client, error) {
// 	const op errors.Op = "redis.NewClient"
// 	// opts.ReadTimeout = time.Minute

// 	client := redis.NewClient(opts)
// 	_, err := client.Ping().Result()
// 	if err != nil {
// 		return nil, errors.E(op, err, errors.KindUnexpected)
// 	}
// 	return client, nil
// }
