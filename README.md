# About
This project emulates a basic backend for a food ordering service. In addition to the basic requirements for a food ordering service, it also implements basic monitoring with Prometheus.

# Prebuild Installations
1. Docker
2. Go
3. Protoc

# Build Steps
1. Clone the project into $HOME/go/src/github.com/avvarikrish
2. cd into project directory and run ./setup.sh
3. Run docker-compose up -d

# Client Binary Usage
1. ccgobbles

    ./ccgobbles function-name
    
        Functions:
            register_user:
                Register a user in the database

            login_user:
                Login a user with the proper credentials (email and password)

            update_user:
                Update user info

            delete_user:
                Delete a user

            add_restaurant:
                Add restaurant to database

            create_order
                Create an order with the proper user email, restaurant id, and order items

2. metrics

    ./metrics metric-name
    
        Metric Names
            average:
                Calculate the average number of items per order in the last [<time-interval>]

            percentile
                Calculate the 95th percentile of number of items per order in the last [<time-interval>]

        Time Intervals
            s:
                Last 5 seconds

            m:
                Last 1 minute

            h:
                Last 1 hour

        
