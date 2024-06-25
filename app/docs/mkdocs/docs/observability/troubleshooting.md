# Troubleshooting guide

This document is a guide for investigating errors, with the goal of resolving the issue or generate a report.  

The target audience of this document are all the tyo of users in the signare. 

**How to use this guide:**

1. Scan through [specific scenarios](#specific-scenarios) to see if any of these applies to you. If any do follow the instructions in the subsection.
2. Scan through [general scenarios](#general-scenarios) to find which scenario(s) applies to you. If any do, follow the instructions to update your instance or collect information for the issue report.
3. If you cannot resolve the problem on your own, create an issue on the [project's repository](https://github.com/hyperledger-labs/signare/issues)


!!! NOTE

    This guide assumes you have site admin privileges and the RBAC default configuration of the signare instance. If you are non-admin user, report your error to the site admin and have them walk through this guide.

## Specific scenarios

### **Scenario: `signare` fails to start because administrator already exist**

If this is the case, you may have already initialized your signare with another administrator name. To confirm, follow these steps:

1. Go to your PostgreSQL instance.
2. Find the `cfg_admin` table.
3. Check if there is any row in the table.

**Solution:** replace the value of the `--signer-administrator` flag with one of the ids in the `cfg_admin` table.

### **Scenario: cannot create HSM slot**

The signare attempts to connect to the HSM slot before creating a slot. If the response is a `412` HTTP status code
the connection has failed.

To confirm, follow these steps:

1. List the slot partitions in your HSM. 
2. Confirm that your `application`'s HSM slot exists in the HSM device.

**Solution 1:** if the slot exists, confirm that the pin of the slot, defined in the request, is correct.

**Solution 2:** if the slot does not exist, manually create the slot in your HSM and repeat the operation.

### **Scenario: create an account for a user fails**

The signare connects to the `application`'s slot in the HSM to validate if the account exists. A `412` HTTP status code
means that the account was not found in the HSM slot. 
 
#### Case A: the HSM slot is not accessible

To confirm, follow these steps:

1. List the slot partitions in your HSM.
2. Confirm that your `application`'s HSM slot exists in the HSM device.

**Solution 1:** if the slot exists, the slot configured for your `application` check if your configured slot's pin is correct.

**Solution 2:** if the slot does not exist, create the slot in your HSM and try again.

#### Case B: the address does not exist in the HSM

To confirm, follow these steps:

1. Call the `eth_accounts` RPC method 
    ```console
    curl --location --request POST 'http://localhost:<rpc_port>' \
    --header 'X-Auth-UserId: <signare_admin>' \
    --header 'Content-Type: application/json' \
    --data-raw '{
       "id": 1,
       "jsonrpc": "2.0",
       "method": "eth_accounts",
       "params": 
           { 
           }
    }'
    ```
2. Confirm that the `address` is not in the result's list 

**Solution 1:** if you want to use this specific `address`, import it to your HSM manually.

**Solution 2:** if you can use any `address`, generate a new account using the signare and replace it in your original request.
```console
    curl --location 'http://localhost:<rpc_port>' \
    --header 'X-Auth-UserId: <applciation_admin>' \
    --header 'X-Auth-ApplicationId: <application>' \
    --header 'Content-Type: application/json' \
    --data '{
        "jsonrpc": "2.0",
        "method": "eth_generateAccount",
        "params": [],
        "id": 1
    }'
```  

### **Scenario: cannot sign a transaction**
#### Case A: forbidden access

If this is the case, the HTTP status code of the response is `200` but the RPC response has the error code `-32099`. To confirm, follow these steps:

1. Read your `X-Auth-UserId` and `X-Auth-ApplicationId` in the request's headers.
2. Get the user configuration by calling the signare endpoint:
   ```console
       curl --location --request GET 'http://localhost:32325/applications/<applicationId>/users/<userId>' \
      --header 'X-Auth-UserId: <user_application_admin>' \
      --header 'Content-Type: application/json' 
   ```  
3. The user's `roles` list must have the `transaction-signer` role
4. The user's `accounts` list must have the account in the `from` field of your transaction.

**Solution 1:** if the point 3 is not satisfied, edit your user's roles list to have the `transaction-signer` role.

**Solution 2:** if the point 4 is not satisfied, add the account in the `from` field of the transaction to your user's `accounts` list.

#### Case B: slot is not available

If this is the case, the HTTP status code of the response is a `200` but RPC response has the error code `-32097`. To confirm, follow these steps:

1. Open your HSM.
2. Check the slots list in your HSM.
3. Confirm the slot is not created.

**Solution 1:** if your application's slot does not exist, create it or change your application slot.

**Solution 2:** if your application slot exists, check if the pin configured for your slot is correct.

## General scenarios

This section contains a list of scenarios, each of which contains instructions that include actions that are appropriate to the scenario. Note that more than one scenario may apply to a given issue.

### **Scenario: cannot remove a resource**

The signare data model define dependencies between the created resources that cannot be broken.

The signare follows a non-intervention approach instead of a fallback removal approach. This means that, if you try to remove a resource, but there are resources 
in the database that depends on it, the request fails instead of automatically removing all the resources that depend on it.

For instance, attempting to remove an `application` resource with associated `user` resources underneath it will result in a failure, typically indicated by a `412` status code.

Steps to follow:

1. Check if the resource you intend to remove has any dependencies in the [resources relations diagram](../reference/database.md#resources-relations).
2. If there are dependencies, list them using the corresponding signare endpoints (e.g. list the `users`).
3. Remove any existing dependencies of the resource you want to remove.
4. Finally, remove the resource itself.
 

### **Scenario: all RPC requests fail**

The signare RPC requests are a proxy to your HSM, all the RPC methods connect with a slot in your HSM.

The signare returns a `200` HTTP status code with a `-32097` RPC error code when your HSM slot is misconfigured. If this is your case, follow these steps:

1. List the HSM devices in the signare, this must be done by a user with the role of `signer-administrator`.
   ```console
       curl --location --request GET 'http://localhost:<port>/admin/modules' \
      --header 'X-Auth-UserId: <signer_admin_user>' \
      --header 'Content-Type: application/json' 
   ```  
2. Find the HSM of your `application` in the response. 
3. List the HSM slots of your `application`'s HSM.
   ```console
       curl --location --request GET 'http://localhost:<port>/admin/modules/<hsmId>/slots' \
      --header 'X-Auth-UserId: <signer_admin_user>' \
      --header 'Content-Type: application/json' 
   ```  
4. Find the HSM slot of your `application` in the response.
5. List the slot partitions in your HSM.
6. Confirm if your `application`'s HSM slot exists in the HSM device.

**Solution 1:** if the slot exists, the slot configured for your `application` may have the wrong pin configured.

**Solution 2:** if the slot does not exist, create the slot in your HSM and repeat the operation.


