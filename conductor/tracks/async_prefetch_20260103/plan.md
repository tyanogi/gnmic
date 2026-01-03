# Implementation Plan - Async Prefetch for Prompt Mode

This plan outlines the steps to implement the asynchronous prefetch feature for the `gnmic prompt` mode.

## Phase 1: Core Logic Implementation (VariableStore & Prefetcher)

- [ ] Task: Implement Thread-Safe VariableStore
    - Create `pkg/app/variables.go`.
    - Define a structure to hold `map[string][]string` protected by `sync.RWMutex`.
    - Implement methods: `Set(category string, values []string)` and `Get(category string) []string`.
    - **Verification**: Write unit tests in `pkg/app/variables_test.go` to ensure concurrent read/write safety.

- [ ] Task: Implement Async Prefetcher Logic
    - Create `pkg/app/prefetch_oc.go`.
    - Implement `StartOpenConfigPrefetch(ctx context.Context, t *Target, store *VariableStore)`.
    - Implement helper function `extractLeafValues(resp *gnmi.GetResponse) []string`.
    - Define constant OpenConfig paths for `interface` and `netinst`.
    - Implement logic to execute gNMI Get requests in parallel using `sync.WaitGroup` with a timeout.
    - **Verification**: Write unit tests in `pkg/app/prefetch_oc_test.go` mocking the gNMI response to verify value extraction and store updates.

## Phase 2: Application Integration & UI Hook

- [ ] Task: Integrate VariableStore into App Context
    - Modify `pkg/app/app.go` (or relevant struct) to initialize and hold the `VariableStore` instance.
    - Update `Init` or constructor methods to ensure the store is ready before the prompt starts.

- [ ] Task: Hook Prefetcher into Prompt Startup
    - Modify `pkg/app/prompt.go` (specifically `runPrompt` or equivalent).
    - Insert the call to `StartOpenConfigPrefetch` immediately after the target connection is established.
    - Ensure it runs in a separate goroutine (`go app.StartOpenConfigPrefetch(...)`).

- [ ] Task: Implement Completer Logic with Store Lookup
    - Modify `completer` function in `pkg/app/prompt.go`.
    - Add logic to detect cursor context (e.g., regex match for `interface[name=` or `network-instance[name=`).
    - If a match is found, query the `VariableStore` for the corresponding category.
    - Convert retrieved values into `prompt.Suggest` items and return them.
    - **Verification**: Manual test using the `gnmic prompt` mode against a mock gNMI server or a lab device to verify suggestions appear.

## Phase 3: Final Verification and Cleanup

- [ ] Task: Error Handling and Logging Review
    - Review `prefetch_oc.go` to ensure all errors (connection, timeout) are logged as debug only and do not crash the UI.
    - Verify that the prompt remains responsive even if the prefetch hangs (simulate with a blocked mock server).

- [ ] Task: Conductor - User Manual Verification 'Final Verification and Cleanup' (Protocol in workflow.md)
