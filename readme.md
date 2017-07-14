# K-Link Registry

The K-Link Registry comprises a solution for registrants managing the applications of a K-Link Network.

Each registrant (e.g. technical staff of an organization) can add applications, such as K-Boxes or websites with a K-Adapter through the graphical user interface the K-Link Registry provides. The added application is then becoming part of the K-Link Network. For each application the registrant gets an authentication token to include in the application.

These applications are then entitled to perform certain actions on the K-Link Network. Which ones? This depends on their permission level:

- `data-search` - allows an application to access the search endpoint.
- `data-add` - allows an application to add a document to the K-Link's search index.
- `data-remove-own` - allows an application to delete a document which was previously added by the same application.
- `data-remove-all` - allows an application to delete any document within the search-index (e.g. a future manage-my-search-index-super-admin-app)

The interaction of applications with K-Link are happening through the K-Search API. K-Search requests the permission level of applications using the API to the K-Link Registry. The K-Link Registry answers if the application is a) part of the network and b) has a certain set of permissions. Based on this information K-Search will continue or block the request from an application to the K-Search API.
