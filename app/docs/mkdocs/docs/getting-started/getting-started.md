# Getting Started

This Quickstart guide will walk you through the entire process step-by-step, from spinning up your signare instance to sign a transaction. After reading this guide, you will be able to spin up the signare, create an account and sign a transaction using that account.

The target audience of this document is every user who needs to start using the signare as quickly as possible.

It's out of the scope of this guide to explain in detail how the signare works and how to build its binary.

## Concepts

Before you start, we suggest you take a look at our [glossary](../glossary/glossary.md) in order to better understand some of the signare's terminology.

## Requirements

In order to start this guide you will need:

- [x] signare binary.
- [x] PostgresSQL database on `localhost:5432`
- [x] SoftHSMv2 installed.


## Starting the signare

1. Change directory to `deployment/database/postgres` and execute the following command to create the database:

    ```console
    psql -h postgres -p 5432 -U postgres -a -f /create-databases.sql
    ```

2. Run the database migration through the signare binary: 

    ```console
    signare upgrade --config <path_to_repository>/deployment/examples
    ```
   
3. Spin up the signare, it's key to use eht `--signer-administrator` flag as in the command bellow to follow this guide successfully:

    ```console
    signare --listen-address 0.0.0.0 --http-port 32325 --rpc-port 4545 --config <path_to_repository>/deployment/examples --signer-administrator owner
    ```

!!! tip

    You can adapt the configuration and flags as needed, for more in depth information take a look at the [configuration reference](../reference/configuration.md).

## Creating an account

There are a number of necessary steps prior to creating an account.

First, you need to have an application and a user with the role of `application-administrator`. Then, you need to configure your HSM and a slot within it for your application. Lastly, **generate an account**.

1. Create an application:

    ```console
    curl --location --request POST 'http://localhost:32325/applications' \
    --header 'X-Auth-UserId: owner' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "meta": {
            "id": "my-first-application"
        },
        "spec": {
            "chainId": "44844"
        }
    }'
    ```

2. Create a user:

    ```console
    curl --location --request POST 'http://localhost:32325/applications/my-first-application/users' \
    --header 'X-Auth-UserId: owner' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "meta": {
            "id": "my-first-user"
        },
        "spec": {
            "roles": [
                "application-admin"
            ],
            "description": "a user authorized to administrate an application"
        }
    }'
    ```

3. Create a module:

    ```console
    curl --location --request POST 'http://localhost:32325/admin/modules' \
    --header 'X-Auth-UserId: owner' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "meta": {
            "id": "my-first-hsm"
        },
        "spec": {
            "configuration": {
                "hsmKind": "softHSM"
            },
            "description": "my first hsm"
        }
    }'
    ```

4. Generate a slot with SoftHSM:

    ```console
    softhsm2-util --init-token --slot 0 --label WALLET-000 --pin userpin --so-pin sopin
    ```

    !!! note

        This step is needed only if you didn't have any slot created in your soft hsm.

5. Run the following command and copy the generated slot id:

    ```console
    softhsm2-util --show-slots 
    ```

6. Create a slot  pasting the generated slot id in the `slot` attribute:

    ```console
    curl --location --request POST 'http://localhost:32325/admin/modules/my-first-hsm/slots' \
    --header 'X-Auth-UserId: owner' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "meta":{
            "id": "my-first-slot"
        },
        "spec": {
            "applicationId": "my-first-application",
            "slot": <your_generated_slot>,
            "pin": <your_slot_pin>
        }
    }'
    ``` 
   
7. Generate a new account:

    ```console
    curl --location --request POST 'http://localhost:4545' \
    --header 'X-Auth-UserId: my-first-user' \
    --header 'X-Auth-ApplicationId: my-first-application' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "jsonrpc": "2.0",
        "method": "eth_generateAccount",
        "params": [],
        "id": 1
    }'
    ```


## Signing a transaction

Let's assume that you have already followed the guide for "Creating an account", now it's the moment to create a user with a `transaction-signer` role and assigning the generated account to it.

1. Create a user:

    ```console
    curl --location --request POST 'http://localhost:32325/applications/my-first-application/users' \
    --header 'X-Auth-UserId: my-first-user' \
    --header 'X-Auth-ApplicationId: my-first-application' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "meta": {
            "id": "my-first-transaction-signer"
        },
        "spec": {
            "roles": [
                "transaction-signer"
            ],
            "description": "a user authorized to sign transactions"
        }
    }'
    ```
   
2. Enable the generated account in the new user profile:

    ```console
    curl --location --request POST 'http://localhost:32325/applications/my-first-application/users/my-first-transaction-signer/accounts' \
    --header 'X-Auth-UserId: my-first-user' \
    --header 'X-Auth-ApplicationId: my-first-application' \
    --header 'Content-Type: application/json' \
    --data-raw '{ 
        "spec": {
            "accounts": [
                <your_generated_account>
            ]
        }
    }'
    ```

3. Sign a transaction using `my-first-transaction-signer` user:

   ```console
   curl --location 'http://localhost:4545' \
   --header 'Content-Type: application/json' \
   --header 'X-Auth-UserId: my-first-transaction-signer' \
   --header 'X-Auth-ApplicationId: my-first-application' \
   --data '{
       "id": 1,
       "jsonrpc": "2.0",
       "method": "eth_signTransaction",
       "params": 
           {
            "from": <your_generated_account>,
            "to": "0xA4F666f1860D2aCbe49b342C87867754a21dE850",
            "gas": "0x3E8",
            "gasPrice": "0x0",
            "value": "",
            "nonce": "0x1",
            "data": "0x1f170873"
           }
   }' 
   ```
