module kms-plugin

go 1.14

replace src/v1beta1 => ./src/v1beta1/v1beta1/

require (
	golang.org/x/net v0.0.0-20201216054612-986b41b23924
	google.golang.org/grpc v1.34.0
	src/v1beta1 v0.0.0-00010101000000-000000000000
)
