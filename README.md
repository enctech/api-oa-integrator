# Online Authorisation Integrator Application

- [Table of Content](#online-authorisation-integrator-application)
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
  - [References](#references)

## Pre-requisites
1. [Docker](https://docs.docker.com/engine/install/ubuntu/)
2. [Docker compose](https://docs.docker.com/compose/install/linux/)

### Deployment
1. Clone this repo (https://github.com/enctech/api-oa-integrator) or download from releases (https://github.com/enctech/api-oa-integrator/releases).
2. Update SSL certificate in [cert](./cert) folder. Create one if there's none.
3. Run `make copy_cert`. This will copy the certificate to backend and frontend folder.
4. Run `make run_application` to start the application. This can be used as restarting the application for updates as well.
5. Head over to [dashboard](https://localhost:3000) to view the application.
6. API documentation is done through [swagger](https://localhost:1323/swagger/index.html#/) and can be viewed here.

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
1. As user entered. SnB will send request to this service with type `IDENTIFICATION`. This is the point where this service will call 3rd parties to verify the user.
2. If user is confirmed belong to any 3rd parties, it will send finalMessageCustomer to SnB with all relevant data.
3. If user is not confirmed belong to any 3rd parties, it will send empty message on finalMessageCustomer and the process will end here.
4. Once user entered the premise, and users is verified belong to 3rd parties, SnB will send request to this service with type `LEAVE_LOOP` indicating that user is already entered. For users that is not marked under any 3rd parties, SnB will not send this request.

### Exit flow
![Exit Flow](./screenshots/oa-exit.png)
1. As user exited. SnB will send request to this service with type `IDENTIFICATION`. OA-integrator will check if user exist in the database. 
2. If user exist, it will send finalMessageCustomer to SnB with all relevant data.
3. If user is not exist in the database, it will send empty message on finalMessageCustomer and the process will end here.
4. SnB will then send request to OA-integrator with type `PAYMENT` to request for payment. OA-integrator will then call 3rd parties to request for payment.
5. If payment is successful, OA-integrator will send finalMessageCustomer to SnB with all relevant data.
6. If payment is not successful, OA-integrator will finalMessageCustomer to SnB with all relevant data but paid amount is 0. Indicated that SnB should proceed with other payment method.
7. User is marked as "Payment success" or "Payment failed" based on payment status.
8. After user left the premise, user is marked with "User left the location".

## Using the application
> :information_source: By default, there is no need for login to view the data. However, there is a login page that can be used to login only to edit configuration.

### Home Page
![Home Page](./screenshots/homepage_1.png)

> :heavy_exclamation_mark: All data displayed here are **only for current day**.

1. Landing page of the application.
2. Total entry indicates the number of success entries (exit included).
3. Total exit indicates the number of session that has already marked requesting for exit.
4. Total payment is total successful payment made to 3rd parties.
5. 3rd parties status indicates the status of all the 3rd parties that are currently connected to the application.
6. SnB status indicates the status of SnB locations.

## References
1. [SNB Documentation](./backend/external-docs/oa-docs.PDF)
2. [Touch n' Go Documentation](./backend/external-docs/tng-docs.pdf)


