package iost

import (
	"time"

	"github.com/iost-official/go-iost/common"
	"github.com/iost-official/go-sdk/pb"
)

// Config of tx
type Config struct {
	GasLimit   float64
	GasRatio   float64
	Expiration int64
	Delay      int64
}

// DefaultTxConfig .
var DefaultTxConfig = Config{
	GasLimit:   1000000,
	GasRatio:   1,
	Expiration: time.Second.Nanoseconds() * 90,
	Delay:      0,
}

// NewTx make a tx with config
func NewTx(config Config) *rpcpb.TransactionRequest {
	ret := &rpcpb.TransactionRequest{
		Time:          time.Now().UnixNano(),
		Actions:       []*rpcpb.Action{},
		Signers:       []string{},
		GasLimit:      config.GasLimit,
		GasRatio:      config.GasRatio,
		Expiration:    time.Second.Nanoseconds()*config.Expiration + time.Now().UnixNano(),
		PublisherSigs: []*rpcpb.Signature{},
		Delay:         config.Delay * 1e9,
		AmountLimit:   []*rpcpb.AmountLimit{},
		ChainId:       1024,
	}
	return ret
}

// AddAction add calls to a tx
func AddAction(t *rpcpb.TransactionRequest, contractID, abi, args string) {
	t.Actions = append(t.Actions, newAction(contractID, abi, args))
}

func actionToBytes(a *rpcpb.Action) []byte {
	se := common.NewSimpleEncoder()
	se.WriteString(a.Contract)
	se.WriteString(a.ActionName)
	se.WriteString(a.Data)
	return se.Bytes()
}

func amountToBytes(a *rpcpb.AmountLimit) []byte {
	se := common.NewSimpleEncoder()
	se.WriteString(a.Token)
	se.WriteString(a.Value)
	return se.Bytes()
}

func signatureToBytes(s *rpcpb.Signature) []byte {
	se := common.NewSimpleEncoder()
	se.WriteByte(byte(s.Algorithm))
	se.WriteBytes(s.Signature)
	se.WriteBytes(s.PublicKey)
	return se.Bytes()
}

func txToBytes(t *rpcpb.TransactionRequest, withSign bool) []byte {
	se := common.NewSimpleEncoder()
	se.WriteInt64(t.Time)
	se.WriteInt64(t.Expiration)
	se.WriteInt64(int64(t.GasRatio * 100))
	se.WriteInt64(int64(t.GasLimit * 100))
	se.WriteInt64(t.Delay)
	se.WriteInt32(int32(t.ChainId))
	se.WriteBytes(nil)
	se.WriteStringSlice(t.Signers)

	actionBytes := make([][]byte, 0, len(t.Actions))
	for _, a := range t.Actions {
		actionBytes = append(actionBytes, actionToBytes(a))
	}
	se.WriteBytesSlice(actionBytes)

	amountBytes := make([][]byte, 0, len(t.AmountLimit))
	for _, a := range t.AmountLimit {
		amountBytes = append(amountBytes, amountToBytes(a))
	}
	se.WriteBytesSlice(amountBytes)

	if withSign {
		signBytes := make([][]byte, 0, len(t.Signatures))
		for _, sig := range t.Signatures {
			signBytes = append(signBytes, signatureToBytes(sig))
		}
		se.WriteBytesSlice(signBytes)
	}

	return se.Bytes()
}


func newAction(contract string, name string, data string) *rpcpb.Action {
	return &rpcpb.Action{
		Contract:   contract,
		ActionName: name,
		Data:       data,
	}
}
