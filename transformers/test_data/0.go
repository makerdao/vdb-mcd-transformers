package test_data

/// Config must be set, in order for test data vars that depend on config to resolve.
/// Vars are evaluated in order, so this must be the first var in the first file (0.go).
var configSet = SetTestConfig()
