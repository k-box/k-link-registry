# Api Contract

The frontend depends on a REST Api in the background to supply it with information.

The following endpoints are expected to be present:

## Authentication `/auth`

* `POST /auth/session` Create a new session.
  > Expects object with:
  >
  > * string `username`
  > * string `password`
  >
  > Returns object:
  >
  > * string `token` token for authenticating 
  >
  > Returns object:
  > * string`error` (string)
