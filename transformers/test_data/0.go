package test_data

var configSet = SetTestConfig()

/*
Our test_data package makes use of constants derived from config files (e.g. contract addresses and ABIs). When other
packages import this one, they will first attempt to resolve all of its package dependencies before attempting to
assign initial values to all of test_data's variables - in lexical order.

If we attempt to assign a variable that depends on reading a value from config before the config file has been set, the
program will panic. Therefore, it is necessary that a config file be setup before other variables that depend on config
are evaluated. This `configSet` variable assignment includes a side effect that sets up a config file
(environments/testing.toml) for viper. It is assigned in 0.go so that the config will be setup before other variables
in subsequent files of this package are assigned.

If you rename this file, be aware that it could trigger panics during test execution if the config is not setup
before variables including values derived from config are assigned.
*/
