# <img src=https://media.tenor.com/Uq_-tDUQlJkAAAAi/hot-beverage-joypixels.gif height="65"/> **Hot-coffee** </div>



### Welcome to the Hot Coffee project — where hot ideas come to life!
#### **Using Postman:** Send a request to `http://localhost:8080`.
#### Clone the repository:
   ```bash
   git clone git@git.platform.alem.school:ukabdoll/hot-coffee.git
```

## About Hot-Coffee

#### Hot-Coffee is a coffee shop management system designed to help manage inventory, menu items, and customer orders through a RESTful API. This project is built with Go and aims to streamline the daily operations of a coffee shop.

#### **Inventory Management:** - Describes actions that can be performed with inventory items.
#### **Menu Management:** - Functions related to adding and editing menu items.
#### **Order Management:** - All functions related to handling orders.
#### **Support for cURL and Postman:** - Testing using popular tools.

INSERT INTO inventory (name, description, quantity, reorder_level, unit, price) VALUES ('Espresso Beans', 'High-quality Arabica beans', 50, 10, 'kg', 25.99);
#### Inventory Endpoints 
| Method  | Path              | Description |
|---------|-------------------|-------------|
| POST    | /inventory        | Create a new inventory item. |
| GET     | /inventory        | Retrieve all inventory information. |
| GET     | /inventory/{id}   | Retrieve information for a specific item by its ID. |
| PUT     | /inventory/{id}   | Edit an existing inventory item by its ID. |
| PUT     | /inventory/       | Put some inventories. |
| DELETE  | /inventory/{id}   | Delete an inventory item. Stock will also be removed. |

#### Menu Endpoints 
| Method  | Path              | Description |
|---------|-------------------|-------------|
| POST    | /menu             | Add a new menu item. |
| GET     | /menu             | Retrieve all menu information. |
| GET     | /menu/{id}        | Retrieve information for a specific item by its ID. |
| PUT     | /menu/{id}        | Edit an existing menu item by its ID. |
| DELETE  | /menu/{id}        | Delete a menu item. |

#### Orders Endpoints 
| Method  | Path              | Description |
|---------|-------------------|-------------|
| POST    | /orders           | Add a new order. |
| GET     | /orders           | Retrieve all order information. |
| GET     | /orders/{id}      | Retrieve information for a specific order by its ID. |
| PUT     | /orders/{id}      | Edit an existing order by its ID. |
| DELETE  | /orders/{id}      | Delete an order. |
| POST    | /orders/{id}/close| Close an order. |

#### Aggregations 
| Method  | Path              | Description |
|---------|-------------------|-------------|
| GET     | /reports/total-sales | Get the total sales amount. |
| GET     | /reports/popular-items | Get a list of popular menu items. |
