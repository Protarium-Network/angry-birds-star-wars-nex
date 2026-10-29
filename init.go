package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/ItzSwirlz/angry-birds-star-wars/database"
	"github.com/ItzSwirlz/angry-birds-star-wars/globals"
)

func init() {
	globals.Logger = plogger.NewLogger()

	_ = godotenv.Load(".env.production")
	_ = godotenv.Load()

	postgresURI := os.Getenv("PN_POSTGRES_URI")
	kerberosPassword := os.Getenv("PN_KERBEROS_PASSWORD")
	authenticationServerPort := os.Getenv("PN_AUTHENTICATION_SERVER_PORT")
	secureServerHost := os.Getenv("PN_SECURE_SERVER_HOST")
	secureServerPort := os.Getenv("PN_SECURE_SERVER_PORT")
	accountGRPCHost := os.Getenv("PN_ACCOUNT_GRPC_HOST")
	accountGRPCPort := os.Getenv("PN_ACCOUNT_GRPC_PORT")
	accountGRPCAPIKey := os.Getenv("PN_ACCOUNT_GRPC_API_KEY")

	if strings.TrimSpace(postgresURI) == "" {
		globals.Logger.Error("PN_POSTGRES_URI environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(kerberosPassword) == "" {
		globals.Logger.Warningf("PN_KERBEROS_PASSWORD environment variable not set. Using default password: %q", globals.KerberosPassword)
	} else {
		globals.KerberosPassword = kerberosPassword
	}

	globals.AuthenticationServerAccount = nex.NewAccount(types.NewPID(1), "Quazal Authentication", globals.KerberosPassword)
	globals.SecureServerAccount = nex.NewAccount(types.NewPID(2), "Quazal Rendez-Vous", globals.KerberosPassword)

	if strings.TrimSpace(authenticationServerPort) == "" {
		globals.Logger.Error("PN_AUTHENTICATION_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(authenticationServerPort); err != nil || port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_AUTHENTICATION_SERVER_PORT is not a valid port. Expected 0-65535, got %s", authenticationServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerHost) == "" {
		globals.Logger.Error("PN_SECURE_SERVER_HOST environment variable not set")
		os.Exit(0)
	}

	if strings.TrimSpace(secureServerPort) == "" {
		globals.Logger.Error("PN_SECURE_SERVER_PORT environment variable not set")
		os.Exit(0)
	}

	if port, err := strconv.Atoi(secureServerPort); err != nil || port < 0 || port > 65535 {
		globals.Logger.Errorf("PN_SECURE_SERVER_PORT is not a valid port. Expected 0-65535, got %s", secureServerPort)
		os.Exit(0)
	}

	if strings.TrimSpace(accountGRPCHost) != "" && strings.TrimSpace(accountGRPCPort) != "" {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", accountGRPCHost, accountGRPCPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			globals.Logger.Warningf("Failed to connect to account gRPC server: %v", err)
		} else {
			globals.GRPCAccountClientConnection = conn
			globals.GRPCAccountClient = pb.NewAccountClient(conn)
			globals.GRPCAccountCommonMetadata = metadata.Pairs("X-API-Key", accountGRPCAPIKey)
			globals.Logger.Success("Connected to Account gRPC!")
		}
	} else {
		globals.Logger.Info("Account gRPC not configured, running in standalone mode")
	}

	database.ConnectPostgres()
}
