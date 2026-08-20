package stateprobe
import("testing";"time"; execution "github.com/example/wasm-sandbox-runtime/internal/domain/execution")
func TestCheckpointRecoveryReachesSucceeded(t *testing.T){e:=execution.Execution{State:execution.Checkpointed,Attempts:[]execution.Attempt{{StartedAt:time.Now()}}};if err:=e.FinishAttempt("sha256:ok",time.Now());err!=nil{t.Fatal(err)};if e.State!=execution.Succeeded||!e.HasResult(){t.Fatal("recovery did not succeed")}}
func TestRecoveredSummaryIsTerminal(t *testing.T){if !execution.Summarize(execution.Execution{State:execution.Succeeded}).Terminal(){t.Fatal("succeeded summary was not terminal")}}
func TestRecoveredFilterIncludesSucceeded(t *testing.T){if !(execution.Filter{States:[]execution.State{execution.Succeeded}}).Recovered(execution.Execution{State:execution.Succeeded}){t.Fatal("recovered state was filtered out")}}
func TestSuccessfulAttemptClearsPreviousError(t *testing.T){now:=time.Now();e:=execution.Execution{State:execution.Running,Error:"old",Attempts:[]execution.Attempt{{StartedAt:now}}};if err:=e.FinishAttempt("sha256:ok",now);err!=nil||e.Error!=""{t.Fatal("previous attempt error remained")}}
