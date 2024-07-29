# Online Authorisation Integrator Application

## Table of Contents
- [Terminology](#terminology)
- [Pre-requisites](#pre-requisites)
- [Deployment](#deployment)
- [Application Architecture](#application-architecture)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [SSL/TLS](#ssltls)
    - [Database backup](#database-backup)
- [Understanding the system](#understanding-the-system)
    - [Entry flow](#entry-flow)
    - [Exit flow](#exit-flow)
- [Using the application](#using-the-application)
    - [Home Page](#home-page)
    - [Transactions](#transactions)
        - [Online Authorisation Transactions](#online-authorisation-transactions)
            - [Default](#default)
            - [Filter](#filter)
    - [Logs](#logs)
    - [Login Page](#login-page)
    - [Configuration](#configuration)
        - [SnB Configuration](#snb-configuration)
        - [3rd Party Configuration](#3rd-party-configuration)
            - [TnG Configuration](#tng-configuration)
- [References](#references)

## Terminology

1. Here is some terminology that will be used in this system (which also mostly based on SnB)

| Term                 | Description                                                                                                 |
|----------------------|-------------------------------------------------------------------------------------------------------------|
| 3rd parties          | Integrators that wants to integrate with premise parking system                                             |
| OA-integrator        | The name of this application                                                                                |
| Job                  | Every entry / exit will trigger a job. Each job will have a specific id defined which is used by SnB        |
| Facility             | Parking Location / Premise                                                                                  |
| Lane                 | Entry / Exit lane that customer used to enter or exit the premise                                           |
| Job - IDENTIFICATION | On every entry / exit, SnB will create job for identification to check if customer belongs to any 3rd party |
| Job - LEAVE_LOOP     | Vehicle entered / left the premise                                                                          |
| Job - PAYMENT        | SnB will create this job to inform OA-integrator to perform payment from the 3rd party                      |

## Pre-requisites

1. [Docker](https://docs.docker.com/engine/install/ubuntu/)
2. [Docker compose](https://docs.docker.com/compose/install/linux/)

### Deployment

# Step-by-Step Instructions

## 1. Create a Directory

Open your Ubuntu terminal and execute the following commands:

1. **Create a directory named `api-oa-integrator` under `/home/Downloads` and navigate to it:**
    ```sh
    mkdir -p /home/Downloads/api-oa-integrator
    cd /home/Downloads/api-oa-integrator
    ```

## 2. Clone the Repository

1. **Clone the repository:**
    ```sh
    git clone https://github.com/enctech/api-oa-integrator.git
    ```
    When prompted, enter the following credentials:
    ```
    Username for 'https://github.com': enctech
    Password for 'https://enctech@github.com': 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
    ```
    **Note:** If you encounter permission issues, please get a personal access token from enctech Admin.

## 3. Create `cert` Folder

1. **Navigate to the `cert` folder within the cloned repository:**
    ```sh
    cd /home/Downloads/api-oa-integrator/api-oa-integrator
    mkdir -p cert
    cd cert
    ```

## 4. Update SSL Certificate

1. **Update or create the SSL certificate:**
    - To create a new certificate, you can use a tool like `openssl`:
        ```sh
        openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout cert.key -out cert.crt
        ```
    - During the creation process, use the following details:
        ```
        Country Name (2 letter code) [AU]:MY
        State or Province Name (full name) [Some-State]:Selangor
        Locality Name (eg, city) []:Puchong
        Organization Name (eg, company) [Internet Widgits Pty Ltd]:Enctech Services Sdn Bhd
        Organizational Unit Name (eg, section) []:Parking
        Common Name (e.g. server FQDN or YOUR name) []:enctech.com
        Email Address []:enctech@gmail.com
        ```

## 5. Copy the Certificate

1. **Navigate back to the root of the cloned repository:**
    ```sh
    cd /home/Downloads/api-oa-integrator/api-oa-integrator
    ```

2. **Run the `make copy_cert` command to copy the certificate to the backend and frontend folders:**
    ```sh
    make copy_cert
    ```

## 6. Run the Application

1. **Start the application using the `make run_application` command:**
    ```sh
    make run_application
    ```

    This command will perform the following steps:
    - **Update the repository:**
        ```sh
        make update
        ```
    - **Build the new version:**
        ```sh
        make build_new
        ```
    - **Clear old images:**
        ```sh
        make clear_images
        ```
    - **Pull the latest changes from the repository:**
        ```sh
        git pull
        ```
        When prompted, enter the following credentials:
        ```
        Username for 'https://github.com': enctech
        Password for 'https://enctech@github.com': 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
        ```
    **Note:** This command will start the application and can also be used to restart it for updates.

## 7. Test the Application

1. **Open your web browser and navigate to:**
    - [https://localhost:3000](https://localhost:3000) or
    - [https://[ipaddress]:3000](https://[ipaddress]:3000)


## Application Architecture

### Backend

1. Language [Go](https://golang.org/).
2. Framework [Echo](https://echo.labstack.com/).
3. [Postgres](https://www.postgresql.org/) database.

### Frontend

1. Language [Typescript](https://www.typescriptlang.org/).
2. Framework [React](https://reactjs.org/).
3. UI Framework
    1. [Material UI](https://material-ui.com/).
    2. [React Hook Form](https://react-hook-form.com/).
    3. [Tailwind CSS](https://tailwindcss.com/).

### SSL/TLS

1. The application is served over HTTPS (currently using self-signed certificate).
2. Backend and Frontend are served with same SSL certificate.

### Database backup

1. Database backup scripts is defined in [scripts](./scripts/db_backup.sh).
2. Suggestion to use cronjob to run the script periodically.
3. Run `chmod +x ./scripts/db_backup.sh` to make the script executable.
4. Run `crontab -e` to edit the cronjob.
5. Add the following line to the end of the file (This is for running cronjob every 12 hours).
   ```bash
   0 */12 * * * /path/to/script/db_backup.sh
   ```

## Understanding the system

1. This service act as middleware between the 3rd parties and SnB.
2. All final data should be finalized in SnB or 3rd parties.

### Entry flow

![Entry Flow](./screenshots/oa-entry.png)

1. As user entered. SnB will send request to this service with type `IDENTIFICATION`. This is the point where this
   service will call 3rd parties to verify the user.
2. If user is confirmed belong to any 3rd parties, it will send finalMessageCustomer to SnB with all relevant data.
3. If user is not confirmed belong to any 3rd parties, it will send empty message on finalMessageCustomer and the
   process will end here.
4. Once user entered the premise, and users is verified belong to 3rd parties, SnB will send request to this service
   with type `LEAVE_LOOP` indicating that user is already entered. For users that is not marked under any 3rd parties,
   SnB will not send this request.

### Exit flow

![Exit Flow](./screenshots/oa-exit.png)

1. As user exited. SnB will send request to this service with type `IDENTIFICATION`. OA-integrator will check if user
   exist in the database.
2. If user exist, it will send finalMessageCustomer to SnB with all relevant data.
3. If user is not exist in the database, it will send empty message on finalMessageCustomer and the process will end
   here.
4. SnB will then send request to OA-integrator with type `PAYMENT` to request for payment. OA-integrator will then call
   3rd parties to request for payment.
5. If payment is successful, OA-integrator will send finalMessageCustomer to SnB with all relevant data.
6. If payment is not successful, OA-integrator will finalMessageCustomer to SnB with all relevant data but paid amount
   is 0. Indicated that SnB should proceed with other payment method.
7. User is marked as "Payment success" or "Payment failed" based on payment status.
8. After user left the premise, user is marked with "User left the location".

## Using the application

> :information_source: By default, there is no need for login to view the data. However, there is a login page that can
> be used to login only to edit configuration.

### Home Page

![Home Page](./screenshots/homepage_1.png)

> :heavy_exclamation_mark: All data displayed here are **only for current day**.
> :heavy_exclamation_mark: This page does not require login

1. Landing page of the application.
2. Total entry indicates the number of success entries (exit included).
3. Total exit indicates the number of session that has already marked requesting for exit.
4. Total payment is total successful payment made to 3rd parties.
5. 3rd parties status indicates the status of all the 3rd parties that are currently connected to the application.
6. SnB status indicates the status of SnB locations. The locations are currently in location code that is used by SnB to
   call this service.

### Transactions

> :heavy_exclamation_mark: This page does not require login

#### Online Authorisation Transactions

##### Default

![OA Transactions](./screenshots/transactions-oa.png)

1. This page shows all transactions that are made through Online Authorisation.
2. It will show every status of the transaction that comes in (
   Refer [Understanding the system](#understanding-the-system))

##### Filter

![OA Transactions](./screenshots/transactions-oa-filter.png)

1. This page allows filter to scope down the transactions to find specific condition.
2. Current filters available are:
    1. Date range
    2. Entry Lane
    3. Exit Lane
    4. License Plate Number
    5. Facility
    6. Job

### Logs

> :heavy_exclamation_mark: This page does not require login

![logs](./screenshots/logs.png)

1. This page shows all the application logs.
2. The logs are limited to application logs, not system logs.

#### Filter
1. As there are lots of logs generated, the logs can be filtered by:
    1. Date range
    2. Message
    3. Fields
2. Every logs carry message and fields. Fields are data that carried with logs.
3. More info in fields can be found in `Fields` column.
4. Typically LPN is search in the fields.

### Login Page

![login](./screenshots/login.png)

1. Login button can be accessed from the top right corner, and it's always visible in every page.
2. Enter username and password to login.
3. Currently, there is no recover password page, in case of forgotten password, get admin to recreate the account.
4. Upon login, logout button will appear at the bottom left of the page, and it's always visible in every page.

### Configuration

#### SnB Configuration
![oa-config-list](./screenshots/oa-config-list.png)
![oa-config-details](./screenshots/oa-config-details.png)

1. The default page shows listing of the SnB configurations set. By default, this will be empty.
2. Click on the `New` button to create new configuration or click on the item to edit the configuration.
3. This system is meant to be handle multiple SnB configurations defined by Facilities and Lanes.
4. Each of the value is defined by SnB system.

#### 3rd Party Configuration
![3rd-party-list](./screenshots/3rd-party-list.png)
![3rd-party-details](./screenshots/3rd-party-details.png)

1. The default page shows listing of the 3rd party configurations set. By default, this will be empty.
2. Click on the `New` button to create new configuration or click on the item to edit the configuration.
3. This system is meant to be handle multiple 3rd party configurations, and there are custom configuration defined by
   the 3rd party.
4. Name field is the name of the 3rd party.
5. URL field is the URL that will be called by this system to the 3rd party.
6. ProviderID field is the ID that will be used by SnB to identify the 3rd party.
7. ClientID and Service Provider field is identifier that is used by 3rd party. This value should be defined by 3rd party.
8. Tax Rate and Surcharge is defined by the business.
9. Surcharge comes with 2 types, which is percentage and fixed amount. This is defined by the business.
10. 3rd Party selection is used to select which 3rd party to be used for the configuration. This is defined by the
    business. The values are pre-defined by the system upon integrating to the 3rd party.
11. Integration to new 3rd party is done manually by the developer following the spec of the 3rd party system.
12. Plaza ID mapper is for mapping the location id of the 3rd party to the location id of SnB. This is defined by the
    3rd party as well as SnB system.

##### TnG Configuration
1. TnG specific configuration is Private SSH Key.
2. As per stated by TnG representative
> Please create using online tool here : https://cryptotools.net/rsagen. Use 2048 key length. Please share only public key to TNG and use private key to generate signature.

## References

1. [SNB Documentation](./backend/external-docs/oa-docs.PDF)
2. [Touch n' Go Documentation](./backend/external-docs/tng-docs.pdf)


