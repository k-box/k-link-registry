**### K-Link Registry: User workflows

K-Link Registry has a basic, but complete, user handling. Users are usually called `registrants` in the K-Link Registry. There are two trustfully specifications, which need to be made on users and to be modifiable:

* **Email address**
* **Passwords**

So, we need proper verification workflows:

### Email verification

Brief proposal:

* **All email added to a registrant needs to be verified.**
* A standard process is established to verify over through a verification link including the email and an randomly generated token. For exmaple `https://registry.foo.bar/?email=foo@bar.com&verification_token=12F5726B7EA91F75BED9AA993C30F87BA93EBDFBC669E68828B1A6B9CCC5FCDC`
* Token and email are getting verified against the data base, and on success a message is diplayed to the registrant and the email is taken over into the account's information.

*(Email verification works over a link with a token (random hash), which is getting compared to the one created and passed by url arguments to the application. The verification can happen either logged in or not. A positive message will be displayed anyway.)*

### Password renewal/recovery

Brief proposal:

* **Password change only through password renewal link.** (Also on an account page, registrants would send a password renewal/recovery mail to change the password)
* A standard process is established to verify over through a verification link including the email and an randomly generated token. For exmaple `https://registry.foo.bar/password?registrant=5&verification_token=12F5726B7EA91F75BED9AA993C30F87BA93EBDFBC669E68828B1A6B9CCC5FCDC`
* Token and registrant_id are getting verified against the data base, and on success - first - the entry is deleted from the database and - second - a form is displayed for new password definition.

### Registrant workflows

In order to show these on concrete examples I'm drafting the different scenarios an email or a password could possibly be changed.
*I'm happy to provide also activity diagrams for them. But writing them just down was easier for now.*

This fits accordingly to the design proposals made in issue #295- K-Link-Registry fronted.

##### 1. Change Email

- **Registrant** logs into the system.
- **Registrant** visits the account page.
- **Registrant** applies changes to the email address.
- **Email** is sent to the email the registrant specified in the recent changes, providing them with an email verification link.
- **Registrant** clicks on the verification link
- **Message** about the successful email verification is presented to the registrant.

##### 2. Password change

###### 2.1. Recover password

- **Registrant** visit the Recover Password page on the system.
- **Registrant** introduces the email address associated to their account.
- **Email** is sent to the email the registrant specified in the account, providing them with an email an password renewal link.
- **Registrant** clicks on the password renewal link and gets a form to introduce a new password.
- **Registrant** introduces two times the same valid password.
- **Message** about the successful password verification is presented to the registrant.
- If user has a timestamp in `last_login` and **Email** is sent out to the address of the account confirming a change of the password.


###### 2.2. Change password

- **Registrant** signs up on the system.
- **Registrant** visits the account page.
- **Registrant** clicks on the `Send link to change password` option.
- **Email** is sent to the email the registrant specified in the account, providing them with an email an password renewal link.
- **Registrant** clicks on the password renewal link and gets a form to introduce a new password.
- **Registrant** introduces two times the same valid password.
- **Message** about the successful password verification is presented to the registrant.
- If user has a timestamp in `last_login` and **Email** is sent out to the address of the account confirming a change of the password.

##### 3. Account creation

###### 3.1. Create account with admin approval

- **Registrant** signs up on the system.
- **Email** is sent to the email the registrant specified in the sign-up, providing them with an email verification link.
- **Registrant** clicks on the verification link
- **Message** about the successful email verification is presented to the registrant.
- **Message** to please wait until an administrator approves the sign-up.
- **Email** is sent out to all admin users on the system, asking for approval of a new registrant.
- **Admin** approves new registrant through the user interface.
- **Email** notification is sent to registrant that the account is ready to use.
- **Registrant** logs in for the first time.

###### 3.2. Admin creates registrant

- **Admin** creates a registrant through the user interface.
- **Email** is sent to the registrant notifying of the new account and asked to verify the email with a verification link.
- **Registrant** clicks on the verification link
- **Message** about the successful email verification is presented to the registrant.
- **Registrant** logs in for the first time.

##### 4. Invite users

###### 4.1 Send an invite

* **Registrant** requests an invitation by entering affiliated e-mail address & desired access (cannot exceed his own).
* **Administrator** receives request & approves/disapproves.
* **E-mail** is sent to new potential registrant & account confirmation ensues.
* **Message** is sent to both registrants once ready.

##### 5. Access audit

###### 5.1 Request account/app access review

* **Registrant** requests a review with desired access.
* **Administrator** receives request & approves/disapproves.
* **E-mail** is sent to registrant indicating decision result.
* **Message** is sent to registrants once ready.
