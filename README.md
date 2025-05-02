# <img src=https://media.tenor.com/Uq_-tDUQlJkAAAAi/hot-beverage-joypixels.gif height="65"/> **frappuccino**

### Welcome to the frappuccino project — where hot ideas come to life!
#### **Using Postman:** Send a request to `http://localhost:8080`.
#### Clone the repository:
```bash
   git clone git@git.platform.alem.school:ukabdoll/frappuccino.git
```

## About frappuccino

#### frappuccino is a coffee shop management system designed to help manage inventory, menu items, and customer orders through a RESTful API. This project is built with Go and aims to streamline the daily operations of a coffee shop.

#### **Inventory Management:** - Describes actions that can be performed with inventory items.
#### **Menu Management:** - Functions related to adding and editing menu items.
#### **Order Management:** - All functions related to handling orders.
#### **Support for cURL and Postman:** - Testing using popular tools.

## 🔧 Features

- **Inventory Management** – Add, edit, and remove stock
- **Menu Management** – Add or modify menu items
- **Order Management** – Place and close orders
- **Reports** – Sales summaries, popular items, leftovers
- **Tooling Support** – Easily testable via Postman or cURL

## Short tables
### API Operations for Inventory
| #   | Method | Path            | Description                                           |
| --- | ------ | --------------- | ----------------------------------------------------- |
| 1   | POST   | /inventory      | Create a new inventory item.                          |
| 2   | GET    | /inventory      | Retrieve all inventory information.                   |
| 3   | GET    | /inventory/{id} | Retrieve information for a specific item by its ID.   |
| 4   | PUT    | /inventory/{id} | Edit an existing inventory item by its ID.            |
| 5   | DELETE | /inventory/{id} | Delete an inventory item. Stock will also be removed. |

#### Menu Endpoints 
| #   | Method | Path       | Description                                         |
| --- | ------ | ---------- | --------------------------------------------------- |
| 1   | POST   | /menu      | Add a new menu item.                                |
| 2   | GET    | /menu      | Retrieve all menu information.                      |
| 3   | GET    | /menu/{id} | Retrieve information for a specific item by its ID. |
| 4   | PUT    | /menu/{id} | Edit an existing menu item by its ID.               |
| 5   | DELETE | /menu/{id} | Delete a menu item.                                 |

#### Orders Endpoints 
| Method | Path                  | Description                                          |
| ------ | --------------------- | ---------------------------------------------------- |
| POST   | /orders               | Add a new order.                                     |
| GET    | /orders               | Retrieve all order information.                      |
| GET    | /orders/{id}          | Retrieve information for a specific order by its ID. |
| PUT    | /orders/{id}          | Edit an existing order by its ID.                    |
| DELETE | /orders/{id}          | Delete an order.                                     |
| POST   | /orders/{id}/close    | Close an order.                                      |
| POST   | /orders/batch-process | Bulk Order Processing                                |

#### Reports
| Method | Path                                                                  | Description                       |
| ------ | --------------------------------------------------------------------- | --------------------------------- |
| GET    | /reports/total-sales                                                  | Get the total sales amount.       |
| GET    | /reports/popular-items                                                | Get a list of popular menu items. |
| GET    | /reports/numberOfOrderedItems?startDate={startDate}&endDate={endDate} | Number of ordered items.          |
| GET    | /reports/search                                                       | Full Text Search Report           |
| GET    | /reports/orderedItemsByPeriod?period={daymonth}&month={month}         | Ordered items by period           |
| GET    | /reports/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}  | Get leftovers                     |
| GET    | /reports/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}  | Get leftovers                     |


## Example Usage
### Inventory Endpoints

1. **POST /inventory**  - Creates a new inventory item.
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

2. **GET /inventory**  - Fetches the list of all inventory items currently stored in the system.
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

3. **GET /inventory/{id}**  - Fetches details of a specific inventory item based on the provided ID.
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
4. **PUT /inventory/{id}**  - Updates an existing inventory item with new information, such as quantity or name.
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

5. **DELETE /inventory/{id}**  - Removes an inventory item from the system based on the provided ID, and decreases stock.
   - **Example error output:**
      ```json
      {
         "error": "invent : not found - id = 1"
      }
      ```

### Menu Endpoints
1. **POST /menu**  - Add a new menu item.
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
2. **GET /menu**  - Retrieve all menu information.
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
3. **GET /menu/{id}**  - Retrieve information for a specific item by its ID.
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
4. **PUT /menu/{id}**  - Edit an existing menu item by its ID.
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

5. **DELETE /menu/{id}**  - Delete a menu item.
      - **Example error Output:**
         ```json
         {
            "error": "menu : not found"
         }
         ```

## To be con