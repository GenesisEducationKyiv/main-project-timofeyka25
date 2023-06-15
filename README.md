# Welcome

### This project was written to send the latest cryptocurrency rates via email.

* The backend was fully written using **Go** and has been tested in **Postman**
* A **csv** file was used to save the email data.
* The data in the csv file is organized in order, enabling efficient duplicate detection using binary search.
* Created swagger documentation for the implemented endpoints.

## Quick Start

### Get your Google App Credentials

The mailer uses SMTP to send emails. To use this protocol, you need to authenticate with your email and app password. A
guide on how to create this password can be found **[here](https://www.getmailbird.com/gmail-app-password/)**.

### Update .env file

There is an .env.example file in the repository git. You need to rename it to .env and replace the value of
the `SMTP_USERNAME` variable with your email and `SMTP_PASSWORD` with the password you received in the paragraph above.

### Deploy in Docker

The best way to launch a project is to deploy it using docker-compose. If you haven't installed it yet, now is the time
to do so. **[Install Docker](https://docs.docker.com/engine/install/)**.

```
docker-compose up
```

This project deploys by default at the port **8000**.

## API Reference

You can find the documentation for the swagger API after starting the server and going
to http://localhost:8000/swagger/index.html.

## Third-party APIs

**[Coinbase API](https://docs.cloud.coinbase.com/)** is used to get exchange rates of a pair of cryptocurrencies from
cryptocurrency markets. 