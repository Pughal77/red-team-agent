import { TabItem, Tabs } from '@astrojs/starlight/components'

Secrets are a secure way to store and use sensitive values such as API keys, tokens, and passwords across your Daytona sandboxes. They let code running in a sandbox authenticate to external services without hardcoding credentials into your source, snapshots, or sandbox environment variables.

Secrets are organization-scoped, encrypted credentials that Daytona injects into a sandbox's outbound HTTPS traffic without ever placing the plaintext inside the sandbox. Instead of an environment variable holding the actual API key, the sandbox holds an opaque placeholder token. An outbound proxy replaces the placeholder with the real value before the request reaches its destination, and only when the request goes to a host you have allowed.

This lets you give a sandbox access to a credential such as an LLM API key or a database password while keeping the plaintext out of the sandbox environment, file system, process arguments, and logs. Code running in the sandbox can use the credential to reach an approved host, but cannot read the value or send it anywhere else.

## How secrets work

A secret never enters the sandbox in plaintext. Daytona uses an opaque placeholder token and an outbound proxy to swap the placeholder for the real value at request time:

1. You store a secret in your organization. Daytona encrypts the value at rest and assigns it an opaque placeholder token, for example **`dtn_secret_<random_string>`**.
2. When you create a sandbox, you map an environment variable to a secret by name. Daytona sets that environment variable to the placeholder, not to the real value.
3. When the sandbox makes an outbound HTTPS request, the proxy inspects it. If a request header carries a placeholder and the destination host matches the secret's allowlist, the proxy replaces the placeholder with the decrypted value before the request reaches its destination.
4. For any other destination, the placeholder is forwarded unchanged. The real value is never sent to a host you did not approve.

Because the substitution happens in the proxy, the plaintext value is never present inside the sandbox. The sandbox sees the placeholder in its environment, and any request to a non-allowed host carries the harmless placeholder rather than the secret.

A secret has the following fields:

| **Field**         | **Description**                                                                                         |
| ----------------- | ------------------------------------------------------------------------------------------------------- |
| **`name`**        | Unique name within the organization. Used to reference the secret when creating a sandbox.              |
| **`value`**       | The plaintext value. Encrypted at rest and never returned by the API after creation.                    |
| **`description`** | Optional human-readable description.                                                                    |
| **`hosts`**       | Allowlist of hosts the secret may be sent to. Supports exact hosts and `*.` wildcards. When omitted, the secret is unrestricted. |
| **`placeholder`** | Opaque token Daytona generates for the secret. This is the value injected into the sandbox environment. |

### Substitution scope

The proxy substitutes placeholders in HTTPS request headers only. Send the secret in a request header such as `Authorization` or `X-Api-Key` over HTTPS. Placeholders anywhere else pass through unchanged:

- **Plain HTTP requests** are never substituted. Substituting them would send the real value across the network in cleartext.
- **Request bodies** (including JSON) are forwarded as-is. Body substitution is not supported.
- **URL query parameters** are forwarded as-is.

If a service only accepts credentials in the body or query string, secrets cannot authenticate to it. Use header-based authentication.

### Response scrubbing

The proxy also scrubs responses. If an upstream response contains the real secret value, the proxy rewrites it back to the placeholder before the response reaches the sandbox. Code in the sandbox can never read the plaintext, even when the destination echoes it back.

This means echo services cannot confirm substitution: the response always shows the placeholder, even when the request that reached the service carried the real value. To confirm substitution works, make an authenticated call to the real service and check the status code. See [Verify substitution](#verify-substitution).

## Host allowlist

The host allowlist is the set of hosts a secret may be sent to. Set your allowed hosts using the `hosts` array when [creating a secret](#create-a-secret) or [updating a secret](#update-a-secret). The proxy replaces the placeholder only for requests whose destination host matches an entry in the allowlist; a request to any other host carries the unmodified placeholder.

Omitting `hosts` leaves the secret unrestricted: the proxy replaces the placeholder for requests to any host. Set an allowlist for every secret unless you have a specific reason not to.

- **Hosts only**: use hostnames; omit protocols, paths, ports, or query strings
- **Wildcards supported**: prefix a host with `*.` to allow the base domain and its subdomains
- **Case-insensitive**: host matching ignores letter case
- **Scope tightly**: list only the hosts the secret requires

The following examples are valid:

- **Single host**: `["api.example.com"]`
- **Wildcard host**: `["*.example.com"]`
- **Multiple hosts**: `["api.example.com", "*.example.org", "service.example.net"]`

## Create a secret

Daytona provides methods to create secrets. 

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
from daytona import CreateSecretParams, Daytona

daytona = Daytona()

secret = daytona.secret.create(CreateSecretParams(
    name="my-secret",
    value="secret-value",
    description="Optional description",
    hosts=["api.example.com"],
))
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
import { Daytona } from '@daytona/sdk'

const daytona = new Daytona()

const secret = await daytona.secret.create({
  name: 'my-secret',
  value: 'secret-value',
  description: 'Optional description',
  hosts: ['api.example.com'],
})
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
require 'daytona'

daytona = Daytona::Daytona.new

secret = daytona.secret.create(
  'my-secret',
  'secret-value',
  description: 'Optional description',
  hosts: ['api.example.com']
)
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
package main

import (
	"context"
	"log"

	"github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
)

func main() {
	client, err := daytona.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	description := "Optional description"
	secret, err := client.Secret.Create(context.Background(), &types.CreateSecretParams{
		Name:        "my-secret",
		Value:       "secret-value",
		Description: &description,
		Hosts:       []string{"api.example.com"},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Secret ID: %s", secret.ID)
}
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import io.daytona.sdk.Daytona;
import io.daytona.sdk.model.CreateSecretParams;
import io.daytona.sdk.model.Secret;

import java.util.List;

public class App {
    public static void main(String[] args) {
        try (Daytona daytona = new Daytona()) {
            CreateSecretParams params = new CreateSecretParams();
            params.setName("my-secret");
            params.setValue("secret-value");
            params.setDescription("Optional description");
            params.setHosts(List.of("api.example.com"));

            Secret secret = daytona.secret().create(params);
        }
    }
}
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/secret' \
  --request POST \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
  "name": "my-secret",
  "value": "secret-value",
  "description": "Optional description",
  "hosts": ["api.example.com"]
}'
```

</TabItem>
</Tabs>

## Use a secret in a sandbox

Daytona provides methods to inject a secret into a sandbox.

1. When creating the sandbox, pass **`secrets`** as a map of environment variable name to secret name.
2. Daytona sets that environment variable to the secret's placeholder, not to the real value.
3. Your code reads the environment variable as usual and sends it in a request header of an outbound HTTPS request, for example `Authorization` or `X-Api-Key`.
4. The outbound proxy substitutes the real value into request headers sent to the secret's allowed hosts, and leaves the placeholder unchanged for any other host. See [Substitution scope](#substitution-scope) for what is and is not substituted.

Each entry must have exactly one key. The environment variable name (the key) can differ from the secret name (the value), so the same stored secret can be exposed under different variable names in different sandboxes.


<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
from daytona import CreateSandboxFromSnapshotParams, Daytona

daytona = Daytona()

sandbox = daytona.create(CreateSandboxFromSnapshotParams(
    language="python",
    secrets={
        "MY_API_KEY": "my-secret",
    },
))
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
import { Daytona } from '@daytona/sdk'

const daytona = new Daytona()

const sandbox = await daytona.create({
  language: 'typescript',
  secrets: {
    MY_API_KEY: 'my-secret',
  },
})
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
require 'daytona'

daytona = Daytona::Daytona.new

sandbox = daytona.create(Daytona::CreateSandboxFromSnapshotParams.new(
  language: Daytona::CodeLanguage::PYTHON,
  secrets: {
    'MY_API_KEY' => 'my-secret'
  }
))
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
package main

import (
	"context"
	"log"

	"github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
)

func main() {
	client, err := daytona.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	sandbox, err := client.Create(context.Background(), types.SnapshotParams{
		SandboxBaseParams: types.SandboxBaseParams{
			Language: types.CodeLanguagePython,
			Secrets: map[string]string{
				"MY_API_KEY": "my-secret",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sandbox ID: %s", sandbox.ID)
}
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import io.daytona.sdk.Daytona;
import io.daytona.sdk.Sandbox;
import io.daytona.sdk.model.CreateSandboxFromSnapshotParams;

import java.util.Map;

public class App {
    public static void main(String[] args) {
        try (Daytona daytona = new Daytona()) {
            CreateSandboxFromSnapshotParams params = new CreateSandboxFromSnapshotParams();
            params.setLanguage("python");
            params.setSecrets(Map.of("MY_API_KEY", "my-secret"));

            Sandbox sandbox = daytona.create(params);
        }
    }
}
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/sandbox' \
  --request POST \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
  "secrets": [
    { "MY_API_KEY": "my-secret" }
  ]
}'
```

</TabItem>
</Tabs>

## Verify substitution

Do not use echo services to test substitution. [Response scrubbing](#response-scrubbing) rewrites the real value back to the placeholder in every response, so an echo service always shows the placeholder even when substitution worked. This is a false negative, not a failure.

To verify substitution, call the real service and check the status code. For example, store a GitHub token as a secret allowed for `api.github.com`, then call `https://api.github.com/user` from the sandbox. A `200` response proves the proxy replaced the placeholder with the real token, because GitHub rejects the placeholder with a `401`.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
from daytona import CreateSandboxFromSnapshotParams, CreateSecretParams, Daytona

daytona = Daytona()

daytona.secret.create(CreateSecretParams(
    name="github-token",
    value="ghp_your_real_token",
    hosts=["api.github.com"],
))

sandbox = daytona.create(CreateSandboxFromSnapshotParams(
    secrets={"GITHUB_TOKEN": "github-token"},
))

response = sandbox.process.exec(
    'curl -s -o /dev/null -w "%{http_code}" '
    '-H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/user'
)
print(response.result)  # 200 proves substitution works
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
import { Daytona } from '@daytona/sdk'

const daytona = new Daytona()

await daytona.secret.create({
  name: 'github-token',
  value: 'ghp_your_real_token',
  hosts: ['api.github.com'],
})

const sandbox = await daytona.create({
  secrets: {
    GITHUB_TOKEN: 'github-token',
  },
})

const response = await sandbox.process.executeCommand(
  'curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/user'
)
console.log(response.result) // 200 proves substitution works
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
require 'daytona'

daytona = Daytona::Daytona.new

daytona.secret.create(
  'github-token',
  'ghp_your_real_token',
  hosts: ['api.github.com']
)

sandbox = daytona.create(Daytona::CreateSandboxFromSnapshotParams.new(
  secrets: { 'GITHUB_TOKEN' => 'github-token' }
))

response = sandbox.process.exec(
  command: 'curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/user'
)
puts response.result # 200 proves substitution works
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
package main

import (
	"context"
	"log"

	"github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
)

func main() {
	client, err := daytona.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	_, err = client.Secret.Create(ctx, &types.CreateSecretParams{
		Name:  "github-token",
		Value: "ghp_your_real_token",
		Hosts: []string{"api.github.com"},
	})
	if err != nil {
		log.Fatal(err)
	}

	sandbox, err := client.Create(ctx, types.SnapshotParams{
		SandboxBaseParams: types.SandboxBaseParams{
			Secrets: map[string]string{
				"GITHUB_TOKEN": "github-token",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	response, err := sandbox.Process.ExecuteCommand(ctx,
		`curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/user`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(response.Result) // 200 proves substitution works
}
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import io.daytona.sdk.Daytona;
import io.daytona.sdk.Sandbox;
import io.daytona.sdk.model.CreateSandboxFromSnapshotParams;
import io.daytona.sdk.model.CreateSecretParams;
import io.daytona.sdk.model.ExecuteResponse;

import java.util.List;
import java.util.Map;

public class App {
    public static void main(String[] args) {
        try (Daytona daytona = new Daytona()) {
            CreateSecretParams secretParams = new CreateSecretParams();
            secretParams.setName("github-token");
            secretParams.setValue("ghp_your_real_token");
            secretParams.setHosts(List.of("api.github.com"));
            daytona.secret().create(secretParams);

            CreateSandboxFromSnapshotParams params = new CreateSandboxFromSnapshotParams();
            params.setSecrets(Map.of("GITHUB_TOKEN", "github-token"));
            Sandbox sandbox = daytona.create(params);

            ExecuteResponse response = sandbox.getProcess().executeCommand(
                "curl -s -o /dev/null -w \"%{http_code}\" -H \"Authorization: Bearer $GITHUB_TOKEN\" https://api.github.com/user");
            System.out.println(response.getResult()); // 200 proves substitution works
        }
    }
}
```

</TabItem>
</Tabs>

## Update secrets in a sandbox

Daytona provides methods to update secrets mounted in a sandbox without recreating it.

Pass `secrets` as a map of environment variable name to secret name, the same format used when [creating the sandbox](#use-a-secret-in-a-sandbox). The map replaces the previously mounted set: new entries are attached, entries no longer present are detached, and an empty map detaches all secrets.

Attached, detached, and rotated secrets take effect for outbound requests within seconds. New environment variables are visible only to processes spawned after the update; already-running processes keep their environment. A sandbox created without any secrets must be restarted before newly attached secrets work.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
sandbox.update_secrets({
    "MY_API_KEY": "my-secret",
})

# Detach all secrets
sandbox.update_secrets({})
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
await sandbox.updateSecrets({
  MY_API_KEY: 'my-secret',
})

// Detach all secrets
await sandbox.updateSecrets({})
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
sandbox.update_secrets({ 'MY_API_KEY' => 'my-secret' })

# Detach all secrets
sandbox.update_secrets({})
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
err := sandbox.UpdateSecrets(context.Background(), map[string]string{
	"MY_API_KEY": "my-secret",
})

// Detach all secrets
err = sandbox.UpdateSecrets(context.Background(), map[string]string{})
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import java.util.Map;

sandbox.updateSecrets(Map.of("MY_API_KEY", "my-secret"));

// Detach all secrets
sandbox.updateSecrets(Map.of());
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/sandbox/SANDBOX_ID_OR_NAME/secrets' \
  --request PUT \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
  "secrets": [
    { "MY_API_KEY": "my-secret" }
  ]
}'
```

</TabItem>
</Tabs>

## List secrets

Daytona provides methods to list secrets. 

List operations use cursor-based pagination and return one page at a time: the page's secrets in `items`, the `total` count, and a cursor for the next page (unset when there are no more pages). Pass the cursor from a previous response to fetch the next page. List operations return metadata only. Secret values are never returned after creation.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
page = daytona.secret.list()
for secret in page.items:
    print(secret.name)
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
const page = await daytona.secret.list()
page.items.forEach((secret) => console.log(secret.name))
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
page = daytona.secret.list
page.items.each { |secret| puts secret.name }
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
page, err := client.Secret.List(context.Background(), nil)
if err != nil {
	log.Fatal(err)
}
for _, secret := range page.Items {
	log.Println(secret.Name)
}
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import io.daytona.sdk.model.ListSecretsResponse;
import io.daytona.sdk.model.Secret;

ListSecretsResponse page = daytona.secret().list();
for (Secret secret : page.getItems()) {
    System.out.println(secret.getName());
}
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/secret/paginated' \
  --header 'Authorization: Bearer YOUR_API_KEY'
```

</TabItem>
</Tabs>

## Get a secret

Daytona provides methods to get a single secret by its ID. 

Get operations return metadata only. Secret values are never returned after creation.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
secret = daytona.secret.get(secret_id)
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
const secret = await daytona.secret.get(secretId)
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
secret = daytona.secret.get(secret_id)
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
secret, err := client.Secret.Get(context.Background(), secretId)
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
Secret secret = daytona.secret().get(secretId);
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/secret/SECRET_ID' \
  --header 'Authorization: Bearer YOUR_API_KEY'
```

</TabItem>
</Tabs>

## Update a secret

Daytona provides methods to update a secret's value, description, or host allowlist. 

Updating the value rotates the credential. The placeholder stays the same, so existing sandboxes that reference the secret pick up the new value on subsequent requests without recreation. The change takes effect in all sandboxes that use the secret within 15 seconds. Send only the fields you want to change. Omitting a field leaves it unchanged.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
from daytona import UpdateSecretParams

secret = daytona.secret.update(secret_id, UpdateSecretParams(
    value="new-secret-value",
    hosts=["api.example.com", "*.example.com"],
))
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
const secret = await daytona.secret.update(secretId, {
  value: 'new-secret-value',
  hosts: ['api.example.com', '*.example.com'],
})
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
secret = daytona.secret.update(
  secret_id,
  value: 'new-secret-value',
  hosts: ['api.example.com', '*.example.com']
)
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
newValue := "new-secret-value"
secret, err := client.Secret.Update(context.Background(), secretId, &types.UpdateSecretParams{
	Value: &newValue,
	Hosts: []string{"api.example.com", "*.example.com"},
})
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
import io.daytona.sdk.model.UpdateSecretParams;

import java.util.List;

UpdateSecretParams params = new UpdateSecretParams();
params.setValue("new-secret-value");
params.setHosts(List.of("api.example.com", "*.example.com"));

Secret secret = daytona.secret().update(secretId, params);
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/secret/SECRET_ID' \
  --request PATCH \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
  "value": "new-secret-value",
  "hosts": ["api.example.com", "*.example.com"]
}'
```

</TabItem>
</Tabs>

## Delete a secret

Daytona provides methods to delete a secret. 

Deleted secrets cannot be recovered, and sandboxes that reference a deleted secret can no longer resolve its value.

<Tabs syncKey="language">
<TabItem label="Python" icon="seti:python">

```python
daytona.secret.delete(secret_id)
```

</TabItem>
<TabItem label="TypeScript" icon="seti:typescript">

```typescript
await daytona.secret.delete(secretId)
```

</TabItem>
<TabItem label="Ruby" icon="seti:ruby">

```ruby
daytona.secret.delete(secret_id)
```

</TabItem>
<TabItem label="Go" icon="seti:go">

```go
err := client.Secret.Delete(context.Background(), secretId)
```

</TabItem>
<TabItem label="Java" icon="seti:java">

```java
daytona.secret().delete(secretId);
```

</TabItem>
<TabItem label="API" icon="seti:json">

```bash
curl 'https://app.daytona.io/api/secret/SECRET_ID' \
  --request DELETE \
  --header 'Authorization: Bearer YOUR_API_KEY'
```

</TabItem>
</Tabs>

## Permissions

Secrets are scoped to an organization. Managing them requires the `manage:secrets` permission, which you can grant to an [organization role](/docs/en/organizations) or an [API key](/docs/en/api-keys). [Updating the secrets mounted in a sandbox](#update-secrets-in-a-sandbox) is a sandbox operation and requires the `write:sandboxes` permission instead. Create, update, and delete operations are recorded in the [audit logs](/docs/en/audit-logs). Secret values are masked in audit entries.