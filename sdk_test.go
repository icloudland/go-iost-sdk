package iost

import (
	"fmt"
	rpcpb "github.com/iost-official/go-sdk/pb"
	"testing"
	"time"

	"encoding/json"

	"github.com/iost-official/go-iost/common"
)

//var addr = "localhost:30002"
var addr = "54.180.196.80:30002"

func TestGet(t *testing.T) {
	client := NewClient()
	err := client.Dial(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(client.ChainInfo())
	t.Log(client.NodeInfo())
	// ...
	//t.Log(client.TxByHash("abc"))
	t.Log(client.BlockByNumber(25261079, true))

	block, err := client.BlockByNumber(25261079, true)
	if err != nil {
		t.Fatal(err)
	}

	for _, tx := range block.Block.Transactions {
		t.Log(tx.Hash)


		//if tx.TxReceipt.StatusCode != rpcpb.TxReceipt_SUCCESS {
		//	continue
		//}

		for _, action := range tx.Actions {
			if action.Contract != "token.iost" ||
				action.ActionName != "transfer" {
				continue
			}

			var datas []string
			err := json.Unmarshal([]byte(action.Data), &datas)
			if err != nil {
				t.Log(err)
			}
			t.Log("length", len(datas))
			for _, v := range datas {
				t.Log(v)
			}
		}
	}

	client.Close()
}

func TestGetAccount(t *testing.T)  {

	client := NewClient()
	err := client.Dial(addr)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(client.Account("yuxiaobei", true))

	a, err := client.Account("yuxiaobei1", true)

	t.Log(a)
}

func TestGetRatio(t *testing.T)  {
	client := NewClient()
	err := client.Dial(addr)
	if err != nil {
		t.Fatal("error in dial", err)
	}
	defer client.Close()

	t.Log(client.GasRatio())

}

func TestSendTx(t *testing.T) {
	client := NewClient()
	err := client.Dial(addr)
	if err != nil {
		t.Fatal("error in dial", err)
	}
	defer client.Close()

	tx := NewTx(Config{
		GasRatio:   1,
		GasLimit:   100000,
		Delay:      0,
		Expiration: 90,
	})
	args, err := json.Marshal([]string{"iost", "yuxiaobei", "guohua", "2.000", "test"})
	if err != nil {
		t.Fatal(err)
	}
	AddAction(tx, "token.iost", "transfer", string(args))

	kc := NewKeychain("yuxiaobei")
	kc.AddKey(common.Base58Decode("privatekey"), "active")

	tx.AmountLimit = []*rpcpb.AmountLimit{{Token: "*", Value: "unlimited"}}

	kc.SignTx(tx)

	fmt.Println(tx.Time)

	//err = sdk.VerifySignature(tx)
	//if err != nil {
	//	return  err
	//}

	handler := NewHandler(tx, client)
	hash, err := handler.Send()
	t.Log(hash, err)
}

func TestAbc(t *testing.T)  {

	fmt.Println(time.Now().Unix())

	time.Sleep(time.Duration(3*1000) * time.Millisecond)

	fmt.Println(time.Now().Unix())
}
