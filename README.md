# Twitter Lite

## Introduction
Twitter Lite is a simplified social media platform designed to allow users to share short messages, known as “tweets,” with their followers. This project demonstrates a streamlined version of the core functionalities of Twitter, focusing on user management, authentication, and basic tweet operations.

### Key Features:

- **User Management:**
  - **Create User:** Users can create an account to join the platform.
  - **Update User:** Existing users can update their profile information.

- **Authentication:**
  - **Login:** Users can authenticate themselves by logging in with their credentials.
  - **Token Refresh:** Refresh the authentication token to maintain active sessions securely.
  - **Token Revoke:** Users can revoke their authentication tokens, effectively logging out of the system.

- **Tweet Management:**
  - **Post Tweet:** Users can post new tweets to share their thoughts with the community.
  - **Get Tweet by ID:** Retrieve a specific tweet using its unique identifier.
  - **Get All Tweets:** Fetch all tweets posted on the timeline.
  - **Delete Tweet:** Users can delete their tweets by ID, removing them from the platform.

This project is built with a clean and modular architecture, following the Ports and Adapters model, which ensures that the core business logic is decoupled from external dependencies like databases. This design allows for easy adaptability and scalability as the platform grows.

### Ports and Adapters Model:

The Ports and Adapters model, also known as Hexagonal Architecture, is designed to create a loosely coupled system where the core business logic is independent of external factors like databases, user interfaces, and frameworks.

In this architecture:
- Ports represent the interfaces or entry points into the application. They define the operations that can be performed on the application, without specifying how these operations are implemented.
- Adapters are the implementation of these ports. They serve as a bridge between the application’s core logic and external systems, such as databases, web services, or user interfaces.

By using this model, I have structured the repository in such a way that it can be easily switched or replaced with another implementation without affecting the core business logic. This approach ensures that the system remains flexible, maintainable, and testable, allowing for easier integration of new technologies or changes in the external environment.

### APIs:

| API Endpoint                      | Description                                |
|-----------------------------------|--------------------------------------------|
| `POST /users`                     | Creates a new user.                        |
| `PUT /users`                      | Updates an existing user.                  |
| `POST /login`                     | Authenticates and logs in a user.          |
| `POST /refresh`                   | Refreshes the user's authentication token. |
| `POST /revoke`                    | Revokes the user's authentication token.   |
| `POST /tweets`                    | Creates a new tweet.                       |
| `GET /tweets/{tweetId}`           | Retrieves a tweet by its ID.               |
| `GET /tweets`                     | Retrieves all tweets.                      |
| `DELETE /tweets/{tweetId}`        | Deletes a tweet by its ID.                 |
