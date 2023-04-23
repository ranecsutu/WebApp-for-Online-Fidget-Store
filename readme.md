# Fidget Spinner Shop

This is the repo for my fidget shop project.

- Built in Go version 1.18
- Uses [chi router](https://github.com/go-chi/chi)
- Uses [alex edwards SCS](https://github.com/alexedwards/scs/v2)

Prerequisites:
1. You need a stripe account. Take from there STRIPE_SECRET and STRIPE_KEY. Put these variables into the Makefile and Air file.
2. Create a DB Connection and put the necessary information into the Makefile (and Air file if want to use Air). Create the tables by running "soda migrate" at the root path of the project. You can also use "soda reset" to reset the DB to its initial state.

Steps to run:
1. Can use Air or Makefile in order to start front/back end. Use "Make start" "Make stop" to start/stop the front/back-end using make file.
2. To make a payment use the IBAN 4242 4242 4242 4242 (stripe test IBAN that is always accepted).
