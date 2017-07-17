# K-Link Registry

The K-Link Registry creates a application management solution for registrants on a K-Link Network.

Each registrant (e.g.: technical staff of an organization) can add applications, such as K-Boxes or websites with a K-Search client through the graphical user interface the K-Link Registry provides. The added application then becomes an integral part of the K-Link Network. For each application the registrant creates, an associated authentication token is generated for its authorization on the network.

These applications are then entitled to perform certain actions on the K-Link Network. Which ones entirely depends on the permission level assigned to them:

- `data-search` - allows an application to access the search endpoint.
- `data-add` - allows an application to add a document to the K-Link's search index.
- `data-remove-own` - allows an application to delete a document which was previously added by the same application.
- `data-remove-all` - allows an application to delete any document within the search-index (e.g.: a future manage-my-search-index-super-admin-app)

Applications within a K-Link network interact through the K-Search API. A K-Search engine requests the permission level of applications by using the API of the connected K-Link Registry. This registry responds if the application is A) part of the network and B) has a certain set of required permissions. Based on this information the K-Search engine will either reply or refuse the request of an application to the K-Search API.
