# rbac-validator

## Purpose

Validates that a set of files used for RBAC (defining roles, permissions and actions) are correctly defined (there are not broken references).

Specifically, the tool runs the following checks:
* Actions map 1 to 1 with operationIDs.
* Permissions point to existing actions.
* Roles point to existing permissions.
* Every action is pointed by at least one permission that is also pointed by a role (in other words: check that every action is assigned to at least one role).

## How to use

Run `make tools.run_default`. It uses RBAC files located in signare/app/include/rbac.

Use `make tools.help` for more info about the command and its flags.
