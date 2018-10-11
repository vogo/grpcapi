package auth

const (
	//KeyUserID grpc metadata user id key (lowercase required)
	KeyUserID = "x-uid"

	//KeyScope grpc metadata scope key
	KeyScope = "x-scp"

	//KeyRole grpc metadata role key
	KeyRole = "x-rol"

	//Authorization restful api auth header
	Authorization = "Authorization"

	//KeyRequestID grpc metadata request id
	KeyRequestID = "x-req-id"
)
