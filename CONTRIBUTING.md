# Contribution guidelines

## Pull Requests
- `go fmt` and `go vet` are run as part of `make test` and `make integrationtest`, please make sure to check in the format changes.
- Ensure that new code is well tested, including integration testing if applicable.
- Make sure the build is passing.
- Update the README as necessary.
- Once a Pull Request has received one approval it can be merged in by a core developer.

## Formatting SQL
We use Goland's built-in sql formatting for `.sql` files. If you use
Goland, please make sure to format new or changed `.sql` files:
- Update the SQL dialect for the `.sql` file(s).
  - The default dialect will likely be `<Generic SQL>`.
  - To change to `PosgreSQL`, right click the `.sql` file and choose `Change Dialect`.
  - Choose `PostgreSQL`.
- Reformat the `.sql` file: with the file(s) selected, from the top menu choose `Code` > `Reformat
    Code`.

## Import Ordering
We follow a standard of including imports in the following order:
1. core go packages
1. external library packages
1. project packages

## Generating the Changelog
See documentation in VulcanizeDB repository: https://github.com/vulcanize/vulcanizedb/blob/staging/documentation/contributing.md#generating-the-changelog.

## Creating a new migration file
1. `make new_migration NAME=add_columnA_to_table1`
    - This will create a new timestamped migration file in `db/migrations`
1. Write the migration code in the created file, under the respective `goose` pragma
    - Goose automatically runs each migration in a transaction; don't add `BEGIN` and `COMMIT` statements.

## Code of Conduct
Vulcanize follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/1/4/code-of-conduct).
