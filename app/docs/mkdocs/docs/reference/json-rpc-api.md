# JSON RPC API Specification

This section describes the different JSON RPC API endpoints.

The target audience of this document are users that want to interact with the signare.

RBAC considerations are out of scope of this document. For more details about RBAC, please check its [documentation](rbac.md){:target="_blank"}.

## Initial considerations

The signare JSON RPC API follows the [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification){:target="_blank"}

The application always responds with a 200 OK HTTP status code, as the error details are part of the JSON RPC response.

Returned JSON RPC errors follow the JSON-RPC 2.0 specification. However, the specification reserves `-32000` to `-32099` for implementation-defined server errors. The signare defines the following ones:

| Code   | Message             | Description                                                 |
|--------|---------------------|-------------------------------------------------------------|
 | -32097 | Precondition failed | The request can not be executed in the current system state |
 | -32098 | Not found           | A specified resource was not found.                         |
 | -32099 | Unauthorized        | The request was not authorized.                             |

## Ethereum JSON RPC API supported methods

The `eth_signTransaction` method is supported following the [Ethereum's JSON RPC API](https://ethereum.org/es/developers/docs/apis/json-rpc){:target="_blank"} specification.

Please, refer to Ethereum's documentation for details abouts its input and output parameters.

!!! info
    If the ``gasPrice`` field of the request body is not informed, it is set to 0.


## Custom RPC methods

### eth_generateAccount

Generates a new key pair in the HSM slot configured for the application sent in the header and returns the Ethereum address that corresponds to the public key.

* Request:

    It does not receive any parameters.

    Example:
    ```
    curl -X POST -H "X-Auth-UserId: <user>" -H "X-Auth-ApplicationId: <application>" --data '{"jsonrpc":"2.0","method":"eth_generateAccount","params":[], "id":1}' http://localhost:4545
    ```

* Success response:
  
    Example:
    ```
    {"jsonrpc":"2.0","id":1,"result":"0xcc753268336A33e56Da47500D9C786077CC24311"}
    ```
  
* Error responses:

  | Code   | Message             |
  |--------|---------------------|
  | -32602 | Invalid params      | 
  | -32603 | Internal error      | 
  | -32097 | Precondition failed |
  | -32099 | Unauthorized        |


### eth_removeAccount

Removes a key pair from the HSM slot configured for the application sent in the header given the address of the public key.

* Request:
    
    Input parameters:

    | Name    | Type   | Required |
    |---------|--------|----------|
    | Address | String | ✔        |

    Example:
    ```
    curl -X POST -H "X-Auth-UserId: <user>" -H "X-Auth-ApplicationId: <application>" --data '{"jsonrpc":"2.0","method":"eth_removeAccount","params":[{"address": "0xcc753268336A33e56Da47500D9C786077CC24311"}], "id":1}' http://localhost:4545
    ```

* Success response:

    Example:
    ```
    {"jsonrpc":"2.0","id":1,"result":"0xcc753268336A33e56Da47500D9C786077CC24311"}
    ```

* Error responses:

  | Code   | Message             |
  |--------|---------------------|
  | -32602 | Invalid params      | 
  | -32603 | Internal error      | 
  | -32097 | Precondition failed |
  | -32098 | Not found           |
  | -32099 | Unauthorized        |

### eth_accounts

Lists all the key pairs stored in the HSM slot configured for the application sent in the header as an array of the Ethereum addresses that correspond to the stored public keys. 

* Request:

  It does not receive any parameters.

  Example:
    ```
    curl -X POST -H "X-Auth-UserId: <user>" -H "X-Auth-ApplicationId: <application>" --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[], "id":1}' http://localhost:4545
    ```

* Success response:

  Example:
    ```
    {"jsonrpc":"2.0","id":1,"result":["0xa2c16184fA76cD6D16685900292683dF905e4Bf2","0x13c21AE733fD7312b6dE09a5eb9C5710f8177239","0xcc753268336A33e56Da47500D9C786077CC24311"]}
  
* Error responses:

  | Code   | Message             |
  |--------|---------------------|
  | -32602 | Invalid params      | 
  | -32603 | Internal error      | 
  | -32097 | Precondition failed |
  | -32099 | Unauthorized        |
