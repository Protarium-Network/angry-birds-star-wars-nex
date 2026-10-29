package globals

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"os"

	pb "github.com/PretendoNetwork/grpc-go/account"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"google.golang.org/grpc/metadata"
)

func PasswordFromPID(pid *types.PID) (string, uint32) {
	if GRPCAccountClient != nil {
		ctx := metadata.NewOutgoingContext(context.Background(), GRPCAccountCommonMetadata)

		response, err := GRPCAccountClient.GetNEXPassword(ctx, &pb.GetNEXPasswordRequest{Pid: pid.LegacyValue()})
		if err == nil && response != nil && response.Password != "" {
			return response.Password, 0
		}
	}

	pwSecretHex := os.Getenv("PN_NEX_PASSWORD_SECRET")
	if pwSecretHex != "" {
		secret, err := hex.DecodeString(pwSecretHex)
		if err == nil && len(secret) >= 32 {
			// Must match account-server tokenservice (service.go): HMAC-SHA256 over the
			// 8-byte LE PID, base64 raw-URL of the full digest.
			pidBytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(pidBytes, uint64(pid.Value()))
			mac := hmac.New(sha256.New, secret)
			mac.Write(pidBytes)
			return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), 0
		}
	}

	return KerberosPassword, 0
}
