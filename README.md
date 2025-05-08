# <img src=https://media.tenor.com/Uq_-tDUQlJkAAAAi/hot-beverage-joypixels.gif height="65"/>  frappuccino
<!-- ![Frappuccino Logo](https://media.tenor.com/Uq_-tDUQlJkAAAAi/hot-beverage-joypixels.gif) -->

Welcome to the frappuccino project — where hot ideas come to life!

## Table of contents
- [  frappuccino](#--frappuccino)
  - [Table of contents](#table-of-contents)
  - [Introduction](#introduction)
  - [🔧 Features](#-features)
  - [⚙️ Requirements](#️-requirements)
  - [📁 Directory structure](#-directory-structure)
  - [Short tables](#short-tables)
    - [API Operations for inventory](#api-operations-for-inventory)
    - [API Operations for menu](#api-operations-for-menu)
    - [API Operations for order](#api-operations-for-order)
    - [API Operations for report](#api-operations-for-report)
  - [Example Usage](#example-usage)
    - [Inventory Endpoints](#inventory-endpoints)
    - [Menu Endpoints](#menu-endpoints)
    - [Order Endpoints](#order-endpoints)




## Introduction

`frappuccino` is a coffee shop management system designed to help manage inventory, menu items, and customer orders through a RESTful API. This project is built with Go and aims to streamline the daily operations of a coffee shop.

**Inventory Management:** - Describes actions that can be performed with inventory items.
**Menu Management:** - Functions related to adding and editing menu items.
**Order Management:** - All functions related to handling orders.
**Support for cURL and Postman:** - Testing using popular tools.
## 🔧 Features

- **Inventory Management** – Add, edit, and remove stock
- **Menu Management** – Add or modify menu items
- **Order Management** – Place and close orders
- **Reports** – Sales summaries, popular items, leftovers
- **Tooling Support** – Easily testable via Postman or cURL

## ⚙️ Requirements
   To run the Frappuccino project, ensure the following dependencies are installed:

- **Go** 1.24 or higher
- **PostgreSQL** 15+
- **Docker** 
- **make** (for running development tasks)
- **Git**

## 📁 Directory structure
```
frappuccino/
├── cmd/                  # Entrypoints (main.go)
├── internal/             # Application logic (not exposed as public API)
│   ├── handler/          # HTTP handlers (Gin controllers)
│   ├── service/          # Business logic layer
│   ├── repository/       # Data access layer (PostgreSQL/sqlx)
|   └── routes/           # Route definitions and grouping
├── migrations/           # SQL migration files
└── model/                # Data models and DTOs
├── .env                  # Docker Compose services
├── docker-compose.yml    # Docker Compose services
├── Dockerfile            # Docker container setup
├── Makefile              # Task automation
├── go.mod                # Go module file
├── go.sum                # Go checksums
└── README.md             # Project documentation
```

- #### [`cmd/`](cmd/) – Entrypoint (main.go)
- #### [`internal/`](internal/) – Application logic
  - ##### [`handler/`](internal/handler/) – HTTP handlers
  - ##### [`service/`](internal/service/) – Business logic
  - ##### [`dal/`](internal/dal/) – DB access
  - ##### [`routes/`](internal/routes/) – Route definitions
- ##### [`migrations/`](migrations/) – SQL migration files

## Short tables
### API Operations for inventory
| #   | Method | Path               | Description                                           |
| --- | ------ | ------------------ | ----------------------------------------------------- |
| 1   | POST   | /inventory         | Create a new inventory item.                          |
| 2   | GET    | /inventory         | Retrieve all inventory information.                   |
| 3   | GET    | /inventory/{id}    | Retrieve information for a specific item by its ID.   |
| 4   | PUT    | /inventory/{id}    | Edit an existing inventory item by its ID.            |
| 5   | DELETE | /inventory/{id}    | Delete an inventory item. Stock will also be removed. |
| 6   | GET    | /inventory/history | Retrieve all inventory transaction history.           |
| 7   | GET    | /inventory/reorder | Receiving Reordered Inventory                         |

### API Operations for menu
| #   | Method | Path          | Description                                         |
| --- | ------ | ------------- | --------------------------------------------------- |
| 1   | POST   | /menu         | Add a new menu item.                                |
| 2   | GET    | /menu         | Retrieve all menu information.                      |
| 3   | GET    | /menu/{id}    | Retrieve information for a specific item by its ID. |
| 4   | PUT    | /menu/{id}    | Edit an existing menu item by its ID.               |
| 5   | DELETE | /menu/{id}    | Delete a menu item.                                 |
| 6   | GET    | /menu/history | Retrieve all menu price history.                    |

### API Operations for order

| №   | Method | Path                  | Description                                          |
| --- | ------ | --------------------- | ---------------------------------------------------- |
| 1   | POST   | /orders               | Add a new order.                                     |
| 2   | GET    | /orders               | Retrieve all order information.                      |
| 3   | GET    | /orders/{id}          | Retrieve information for a specific order by its ID. |
| 4   | PUT    | /orders/{id}          | Edit an existing order by its ID.                    |
| 5   | DELETE | /orders/{id}          | Delete an order.                                     |
| 6   | POST   | /orders/{id}/close    | Close an order.                                      |
| 7   | POST   | /orders/batch-process | Bulk Order Processing                                |
| 8   | GET    | /orders/history       | Retrieve all order status history.                   |


### API Operations for report
| Method | Path                                                                  | Description                       |
| ------ | --------------------------------------------------------------------- | --------------------------------- |
| GET    | /reports/total-sales                                                  | Get the total sales amount.       |
| GET    | /reports/popular-items                                                | Get a list of popular menu items. |
| GET    | /reports/numberOfOrderedItems?startDate={startDate}&endDate={endDate} | Number of ordered items.          |
| GET    | /reports/search                                                       | Full Text Search Report           |
| GET    | /reports/orderedItemsByPeriod?period={daymonth}&month={month}         | Ordered items by period           |
| GET    | /reports/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}  | Get leftovers                     |


## Example Usage
### Inventory Endpoints
1. ``POST /inventory``  - Creates a new inventory item.
   - **Example Input:**
      ```json
      {
       "name": "Baking Powder",
       "description": "Leavening agent for baking",
       "quantity": 490,
       "reorder_level": 50,
       "unit": "g",
       "price": 250
      }  
     ```

   - **Example Output:**
      ```json
      {
       "message": "inventory: created"
      }
      ```
2. ``GET /inventory``  - Fetches the list of all inventory items currently stored in the system.
    - **Example Output:**
      ```json
      [
         {
         "ingredient_id": 1,
         "name": "Double Chocolate Cake",
         "description": "Rich chocolate layer cake",
         "quantity": 5000,
         "reorder_level": 200,
         "unit": "g",
         "price": 1500
         },
         {
         "ingredient_id": 2,
         "name": "Milk",
         "description": "Fresh dairy milk Chocolate",
         "quantity": 100,
         "reorder_level": 10,
         "unit": "ml",
         "price": 50
         },
      ///
      ]
      ```
3. ``GET /inventory/{id}``  - Fetches details of a specific inventory item based on the provided ID.
    - **Example Output:**
      ```json
      {
         "ingredient_id": 1,
         "name": "Double Chocolate Cake",
         "description": "Rich chocolate layer cake",
         "quantity": 5000,
         "reorder_level": 200,
         "unit": "g",
         "price": 1500
      }
      ```
4. ``PUT /inventory/{id}``  - Updates an existing inventory item with new information, such as quantity or name.
    - **Example Input:**
      ```json
      {
         "ingredient_id": 1, // this line will be ignored
         "name": "My frappuccino",
         "description": "ust",
         "quantity": 5000,
         "reorder_level": 200,
         "unit": "g",
         "price": 1500
      }
      ```
   - **Example output:**
      ```json
      {
         "message": "updated : 1"
      }
      ```
5. ``DELETE /inventory/{id}``  - Removes an inventory item from the system based on the provided ID, and decreases stock.
   - **Example error output:**
      ```json
      {
         "error": "invent : not found - id = 1"
      }
      ```

### Menu Endpoints
1. ``POST /menu``  - Add a new menu item.
   - **Example Input:**
      ```json
      {
         "name": "My Chocolate",
         "description": "is not your Chocolate",
         "tags": [
            "dessert",
            "sweet"
         ],
         "allergens": [
            "gluten",
            "dairy"
         ],
         "price": 5,
         "ingredients": [
            {
                  "inventory_id": 4,
                  "quantity": 200
            },
            ///
         ]
      }
     ```

   - **Example Output:**
      ```json
      {
         "message": "success : menu created:"
      }
      ```
2. ``GET /menu``  - Retrieve all menu information.
   - **Example Output:**
      ```json
      [
         {
            "product_id": 1,
            "name": "Espresso",
            "description": "Rich and bold espresso shot",
            "tags": [
                  "coffee"
            ],
            "allergens": null,
            "price": 3,
            "ingredients": [
                  {
                     "inventory_id": 1,
                     "quantity": 30
                  }
            ]
         },
         ///
      ]
      ```
3. ``GET /menu/{id}``  - Retrieve information for a specific item by its ID.
   - **Example Output:**
      ```json
      {
         "product_id": 1,
         "name": "Espresso",
         "description": "Rich and bold espresso shot",
         "tags": [
            "coffee"
         ],
         "allergens": null,
         "price": 3,
         "ingredients": [
            {
                  "inventory_id": 1,
                  "quantity": 30
            },
            ///
         ]
      }
      ```
4. ``PUT /menu/{id}``  - Edit an existing menu item by its ID.
   - **Example Output:**
      ```json
      {
         "name": "my milk",
         "description": "is not your milk",
         "tags": [
            "coffee"
         ],
         "allergens": null,
         "price": 3,
         "ingredients": [
            {
                  "inventory_id": 1,
                  "quantity": 30
            },
            ///
         ]
      }
      ```
   - **Example Output:**
      ```json
      {
         "message": "Updated Menu by id:  : "
      }
      ```

5. ``DELETE /menu/{id}``  - **Delete a menu item.**
      - **Example error Output:**
         ```json
         {
            "error": "menu : not found"
         }
         ```

### Order Endpoints
1. ``POST /orders``  - **Add a new order.**
   - **Example Input:**
      ```json
      {
         "customer_name": "Mike",
         "allergens": [
               "dairy"
         ],
         "items": [
               {
                  "product_id": 2,
                  "quantity": 1
               }
         ]
      }
     ```

   - **Example Output:**
      ```json
      {
         "message": "succes : order created by : Mike"
      }
      ```
---