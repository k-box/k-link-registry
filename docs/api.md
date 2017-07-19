# K-Link Registry: Access API

## Context

The K-Link Registry offers a user interface for `registrants` to add `applications` to a K-Link network. Each registered `application` will get a authentication token generated, which the administrator manually must place in the API call of one K-Adapter, K-Box or any K-Search.js-Library and which gives access to the K-Search-API.

The K-Link Registry also exposes the **AccessAPI** to the K-Search. The K-Search can request information on certain permissions for a certain application, based on the `app_url` and `auth_token`.

## Permissions

The system shall be able to distinguish the permissions of applications, which can be:

- `data-search` - allows an application to access the search endpoint.
- `data-edit` - allows an application to edit a document located in the K-Link's search index.
- `data-add` - allows an application to add a document to the K-Link's search index.
- `data-view` - allows an application to view a document located in the K-Link's search index.
- `data-remove-own` - allows an application to delete a document which was previously added by the same application.
- `data-remove-all` - allows an application to delete any document within the search-index (e.g.: a future manage-my-search-index-super-admin-app)

### Versioning

The exposed API is versioned according to the [Semantic Versioning](http://semver.org/). The
API version numbers use a `{MAJOR}.{MINOR}` notation. The selection of the API version is through the url.

* **Current version**: 1.0

### Endpoint

* **URL:** `https://{BASE_URL}/api/{MAJOR}.{MINOR}/{METHOD}`

## Basic API calls

### Requests

All API calls are made to `/api`. Arguments can be passed in the `POST` request as JSON.

| Property | Type    | Required   | Description |
| -------- | ------- | ---------- | ----------- |
| `id` | String | ✔ | An identifier established by the Client that MUST contain a String, Number, or NULL value if included. The value SHOULD normally not be Null and Numbers SHOULD NOT contain fractional parts |
| `params` | Object | ✔ | Arguments for the method of the request (see below under "methods") |


### Responses

The response contains a JSON response object:

| Property | Type    | Required   | Description |
| -------- | ------- | ---------- | ----------- |
| `id` | String | ✔ | The identifier established by the client in the request |
| `result` | result object | on success | REQUIRED on success; MUST NOT exist if there was an error |
| `error`  | error object | on error | REQUIRED on failure; MUST NOT exist if there was no error |

* `status`: Usually "`200` (OK)"; but might change slightly to "`201` (Created)" or "`204` (No data)". Should be ignored by the client.

## Methods

The `{METHOD}` executes to one of the offered functions provided by the K-Search API:

## application.authenticate

Get detailed information of the permission level of an application handled in the K-Kink Registry

* URL: `/api/1.0/application.authenticate`

**Request**:

| Property | Type    | Required   | Description |
| -------- | ------- | ---------- | ----------- |
| `id` | String | ✔ | An identifier established by the client that MUST contain a String, Number, or NULL value if included. The value SHOULD normally not be Null and Numbers SHOULD NOT contain fractional parts |
| `params` | Object | ✔ | A simple JSON object |
| `params[app_url]` | String | ✔ | The url a request is coming from |
| `params[app_secret]` | String | ✔ | The provided application secret used for authentication |
| `params[permissions][]` | Array of strings | | Set of permissions. Each permission has to be matched. If omitted, all application permissions will be returned |

**Successful response**

* `status`: `200` (OK)

| Property | Type    | Required   | Description |
| -------- | ------- | ---------- | ----------- |
| `id` | String | ✔ | The identifier established by the client in the request |
| `result` | Object | ✔ | Application object |
| `result[name]` | String | ✔ | The name of the application |
| `result[app_url]` | String | ✔ | The URL where the application is running and send requests |
| `result[app_id]` | String | ✔ | Unique identifier of the application within the K-Link Registry |
| `result[permissions][]` | Array of strings | ✔ | Permission names of the application |
| `result[email]` | string | ✔ | Contact email address of the application administrator |

**Unsuccessful response**

* `status`: `200` (OK)

| Property | Type    | Required   | Description |
| -------- | ------- | ---------- | ----------- |
| `id` | String | ✔ | The identifier established by the client in the request |
| `error` | Object | ✔ | Object with error information |
| `error[code]` | Integer | ✔ | JSON-RPC inspired [error codes](http://www.jsonrpc.org/specification#error_object) |
| `error[message]` | String | ✔ | Human readable error message |
