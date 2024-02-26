# sadeem-user-api

**API Documentation**

**Authentication**

Note that certain endpoints are restricted to users with admin privileges. These endpoints are denoted with "(admin only)" after the request method.

**Endpoints**

**GET Endpoints**

- **GET /api/users/:id (admin only)**
  - **Description:** Retrieves detailed information about a specific user identified by their ID.
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "username": "string",
        "email": "string",
        "image_path": "string",
        "category": [
          {
            "id": "string",
            "name": "string",
            "activated": boolean
          }
        ]
      }
    ]
    ```
  - **Example request:** `GET /api/users/123` (where 123 is the user's ID)
- **GET /api/users/list (admin only)**
  - **Description:** Lists all users with basic information (ID, username, email, image path).
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "username": "string",
        "email": "string",
        "image_path": "string"
      }
    ]
    ```
- **GET /api/category**
  - **Description:** Lists all categories with their names and activated status.
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "name": "string",
        "activated": boolean
      }
    ]
    ```
- **GET /api/users/me**
  - **Description:** Retrieves information about the currently logged-in user.
  - **Response:** (Same structure as GET /api/users/:id)
    ```json
    {
      "id": "string",
      "username": "string",
      "email": "string",
      "image_path": "string",
      "category": [
        {
          "id": "string",
          "name": "string",
          "activated": boolean
        }
      ]
    }
    ```

**POST Endpoints**

- **POST /api/category**

  - **Description:** Creates a new category.
  - **Request Body:**
    ```json
    {
      "name": "string",
      "activated": boolean
    }
    ```
  - **Example request:**

    ```json
    POST /api/category HTTP/1.1
    Content-Type: application/json

    {
      "name": "New Category",
      "activated": true
    }
    ```

- **POST /api/users/:id/categorize (admin only)**

  - **Description:** Adds a specific user (identified by their ID) to a category (specified in the request body).
  - **Request Body:**
    ```json
    {
      "category_id": "string"
    }
    ```
  - **Example request:**

    ````json
    POST /api/users/123/categorize HTTP/1.1
    Content-Type: application/json

    {
      "category_id": "456"
    }
    ``` (where 123 is the user's ID and 456 is the category ID)
    ````

- **POST /api/users/:id/uncategorize (admin only)**
  - **Description:** Removes a specific user (identified by their ID) from a category (specified in the request body).
  - **Request Body:** (Same structure as POST /api/users/:id/categorize)
  - **Example request:** (Same structure as POST /api/users/:id/categorize)
- **POST /api/users/me/change_password**
  - **Description:** Changes the password of the currently logged-in user.
  - **Request Body:**
    ```json
    {
      "old_password": "string",
      "new_password": "string"
    }
    ```

* **POST /api/users/me/change_image**
  - **Description:** Changes the profile image of the currently logged-in user.
  - **Request Body:**
    - **Multipart form data** containing an image file named "image".
  - **Note:** File size limited to 4 MB
* **POST /api/login**
  - **Description:** Logs in a user using their email or username and password.
  - **Request Body:**
    ```json
    {
      "email": "string" (or "username": "string"),
      "password": "string"
    }
    ```
  - **Response:** (Successful login response will typically include an authentication token or other necessary information)
* **POST /api/register**
  - **Description:** Registers a new user.
  - **Request Body:**
    ```json
    {
      "email": "string",
      "username": "string",
      "password": "string"
    }
    ```
  - **Response:** (Successful registration response may include details about the newly created user)

**PUT Endpoints**

- **PUT /api/users/:id (admin only)**

  - **Description:** Edits information about a specific user (identified by their ID).
  - **Request Body:**
    ````json
    {
      "email": "string",
      "username": "string"
    }
    ``` (Only the properties you want to update need to be included)
    ````
  - **Example request:**

    ````json
    PUT /api/users/123 HTTP/1.1
    Content-Type: application/json

    {
      "email": "updated_email@example.com"
    }
    ``` (where 123 is the user's ID)
    ````

- **PUT /api/category/:id (admin only)**
  - **Description:** Edits information about a specific category (identified by its ID).
  - **Request Body:** (Same structure as PUT /api/users/:id)
- **PUT /api/users/me**
  - **Description:** Edits information about the currently logged-in user.
  - **Request Body:** (Same structure as PUT /api/users/:id)

**DELETE Endpoints**

- **DELETE /category/:id (admin only)**
  - **Description:** Deletes a specific category (identified by its ID).

## Running the Project

**Prerequisites:**

- **PostgreSQL:** Ensure you have PostgreSQL installed and running on your system.
- **Database Creation:** Create a new database for your project.

**Steps:**

1. **Execute Schema:**

   - Navigate to the project directory in your terminal.
   - Run the following command to execute the `schema.sql` file and create the necessary tables in your database:

   ```bash
   psql -d <your_database_name> < schema.sql
   ```

   Replace `<your_database_name>` with the actual name of the database you created.

2. **Set Environment Variable:**

   - Set the `DB_SOURCE_NAME` environment variable to your PostgreSQL connection string. This string contains information about your database server, username, password, and database name.

   **Example (Linux/macOS):**

   ```bash
   export DB_SOURCE_NAME="postgresql://<username>:<password>@<host>:<port>/<database_name>"
   ```
