# Online Authorisation Integrator Application

- [Table of Content](#online-authorisation-integrator-application)
  - [Pre-requisites](#pre-requisites)
  - [Application Architecture](#application-architecture)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [SSL/TLS](#ssltls)
  - [Using the application](#using-the-application)

## Pre-requisites
1. [Docker](https://docs.docker.com/engine/install/ubuntu/)
2. [Docker compose](https://docs.docker.com/compose/install/linux/)
3. Run `make run_application` to start the application. This can be used as restarting the application for updates as well.
4. Head over to [dashboard](https://localhost:3000) to view the application.
5. API documentation is done through [swagger](https://localhost:1323/swagger/index.html#/) and can be viewed here.

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
7. Last 10 transactions show 10 most recent transactions based on updated time, not creation time.


