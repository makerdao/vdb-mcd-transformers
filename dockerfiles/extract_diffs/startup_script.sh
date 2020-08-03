#!/bin/sh
# Starts the extractDiffs command

# Verify required args present
MISSING_VAR_MESSAGE=" is required and no value was given"

function testDatabaseVariables() {
  for a in DATABASE_NAME DATABASE_HOSTNAME DATABASE_PORT DATABASE_USER DATABASE_PASSWORD
  do
    eval arg="$"$a
    test $arg
    if [ $? -ne 0 ]; then
      echo $a $MISSING_VAR_MESSAGE
      exit 1
    fi
  done
}

if test -z "$VDB_PG_CONNECT"; then
  # Exits if the variable tests fail
  testDatabaseVariables
  if [ $? -ne 0 ]; then
    exit 1
  fi

  # Construct the connection string for postgres
  VDB_PG_CONNECT=postgresql://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOSTNAME:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable
fi

# Run the DB migrations
echo "Connecting with: $VDB_PG_CONNECT"
./goose -dir db/migrations postgres "$VDB_PG_CONNECT" up

if [ $? -ne 0 ]; then
  echo "Could not run migrations. Are the database details correct?"
  exit 1
fi

# Run extractDiffs
echo "Running extractDiffs..."
./vulcanizedb extractDiffs --watchedAddresses=0x78f2c2af65126834c51822f56be0d7469d7a523e \
    --watchedAddresses=0x5ef30b9986345249bc32d8928b7ee64de9435e39 \
    --watchedAddresses=0xdfe0fb1be2a52cdbf8fb962d5701d7fd0902db9f \
    --watchedAddresses=0xaa745404d55f88c108a28c86abe7b5a1e7817c07 \
    --watchedAddresses=0xd8a04f5412223f513dc55f839574430f5ec15531 \
    --watchedAddresses=0xAbBCB9Ae89cDD3C27E02D279480C7fF33083249b \
    --watchedAddresses=0x5432b2f3c0dff95aa191c45e5cbd539e2820ae72 \
    --watchedAddresses=0xba3f6a74BD12Cf1e48d4416c7b50963cA98AfD61 \
    --watchedAddresses=0xE6ed1d09a19Bd335f051d78D5d22dF3bfF2c28B1 \
    --watchedAddresses=0xec25Ca3fFa512afbb1784E17f1D414E16D01794F \
    --watchedAddresses=0x3E115d85D4d7253b05fEc9C0bB5b08383C2b0603 \
    --watchedAddresses=0x08c89251FC058cC97d5bA5F06F95026C0A5CF9B0 \
    --watchedAddresses=0x4d95a049d5b0b7d32058cd3f2163015747522e99 \
    --watchedAddresses=0x19c0976f590d67707e62397c87829d896dc0f1f1 \
    --watchedAddresses=0x197e90f9fad81970ba7976f33cbd77088e5d7cf7 \
    --watchedAddresses=0x65c79fcb50ca1594b025960e539ed7a9a6d434a3 \
    --watchedAddresses=0x35d1b3f3d7966a1dfe207aa4514c12a259a0492b \
    --watchedAddresses=0xa950524441892a31ebddf91d3ceefa04bf454466 \
    --watchedAddresses=0x18B4633D6E39870f398597f3c1bA8c4A41294966 \
    --watchedAddresses=0x64DE91F5A373Cd4c28de3600cB34C7C6cE410C85 \
    --watchedAddresses=0x83076a2F42dc1925537165045c9FDe9A4B71AD97 \
    --watchedAddresses=0xe0F30cb149fAADC7247E953746Be9BbBB6B5751f \
    --watchedAddresses=0x956ecD6a9A9A0d84e8eB4e6BaaC09329E202E55e \
    --watchedAddresses=0x39755357759cE0d7f32dC8dC45414CCa409AE24e \
    --watchedAddresses=0x794e6e91555438afc3ccf1c5076a74f42133d08d \
    --watchedAddresses=0xb4eb54af9cc7882df0121d26c5b97e802915abe6 \
    --watchedAddresses=0x81FE72B5A8d1A857d176C3E7d5Bd2679A9B85763 \
    --watchedAddresses=0xf36B79BD4C0904A5F350F1e4f776B81208c13069 \
    --watchedAddresses=0x77b68899b99b686F415d074278a9a16b336085A0 \
    --watchedAddresses=0xf185d0682d50819263941e5f4EacC763CC5C6C42 \
    --watchedAddresses=0x7382c066801E7Acb2299aC8562847B9883f5CD3c
