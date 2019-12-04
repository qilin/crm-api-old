# Database structure

## Schemas

There are two schemas:
- users: stores user's data
- store: stores storefront's data 

### Users schema

This schema stores user accounts, log of user account modifications, authentication log.

#### Database structure
- users: stores user data and profiles
- auth_providers: stores information about connected auth providers to user's profile
- auth_log: stores authentication logs
- users_params_log: stores history of user profile changes


### Store schema

This schema stores game info, storefront blocks.