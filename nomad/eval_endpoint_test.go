package nomad

import (
	"reflect"
	"testing"

	"github.com/hashicorp/net-rpc-msgpackrpc"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/testutil"
)

func TestEvalEndpoint_GetEval(t *testing.T) {
	s1 := testServer(t, nil)
	defer s1.Shutdown()
	codec := rpcClient(t, s1)
	testutil.WaitForLeader(t, s1.RPC)

	// Create the register request
	eval1 := mockEval()
	s1.fsm.State().UpsertEval(1000, eval1)

	// Lookup the eval
	get := &structs.EvalSpecificRequest{
		EvalID:       eval1.ID,
		WriteRequest: structs.WriteRequest{Region: "region1"},
	}
	var resp structs.SingleEvalResponse
	if err := msgpackrpc.CallWithCodec(codec, "Eval.GetEval", get, &resp); err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Index != 1000 {
		t.Fatalf("Bad index: %d %d", resp.Index, 1000)
	}

	if !reflect.DeepEqual(eval1, resp.Eval) {
		t.Fatalf("bad: %#v %#v", eval1, resp.Eval)
	}

	// Lookup non-existing node
	get.EvalID = generateUUID()
	if err := msgpackrpc.CallWithCodec(codec, "Eval.GetEval", get, &resp); err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Index != 1000 {
		t.Fatalf("Bad index: %d %d", resp.Index, 1000)
	}
	if resp.Eval != nil {
		t.Fatalf("unexpected eval")
	}
}