package api

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/register --config register-ogen.yaml openapi.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/login --config login-ogen.yaml openapi.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/orders --config orders-ogen.yaml openapi.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/balance --config balance-ogen.yaml openapi.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/withdraw --config withdraw-ogen.yaml openapi.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --loglevel error --clean --target gen/withdrawals --config withdrawals-ogen.yaml openapi.yaml
