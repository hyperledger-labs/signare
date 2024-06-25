# How to sign a transaction 

This document explains how to sign transactions using the signare. After reading this document, the reader will be able to sign Ethereum transactions using the application.

The target audience of these documents are users that want to learn how to effectively utilize the application.

## Introduction

In order to sign an Ethereum transaction, you must have a configured user with a valid account. In order to do so, you can
check our [How to configure users](how-to-configure-users.md) and [How to create an account](how-to-create-an-account.md) guides

## Signing an Ethereum transaction

As long as you have a user configured with a valid account, you can sign an Ethereum transaction by making the following request to the application: 

```code
curl --location --request POST 'http://localhost:4545' \
--header 'X-Auth-UserId: user-admin' \
--header 'X-Auth-ApplicationId: application' \
--header 'Content-Type: application/json' \
--data-raw '{
    "jsonrpc": "2.0",
    "method": "eth_signTransaction",
    "params": [
        {
            "from": "0x161735Fca56b7dd36768E4a6bbaEA731886fCfcA",
            "to": "0xA4F666f1860D2aCbe49b342C87867754a21dE850",
            "gas": "0x3E8",
            "gasPrice": "0x0",
            "value": "",
            "nonce": "0x1",
            "data": "0x1f170873"
        }
    ],
    "id": 1
}'
```

!!! note 
    Remember that you need valid Ethereum transaction data in order for this request to work. 

